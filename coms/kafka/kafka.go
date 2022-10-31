package kafka

import (
	"errors"
	"log"
	"sync"

	"github.com/Shopify/sarama"
	cluster "github.com/bsm/sarama-cluster"
	"github.com/hulklab/yago"
	"github.com/hulklab/yago/libs/str"
)

type Kafka struct {
	connect       []string
	config        *cluster.Config
	asyncProducer *AsyncProducer
	syncProducer  *SyncProducer
	mu            sync.Mutex
}

// 返回 kafka 组件单例
func Ins(id ...string) *Kafka {

	var name string

	if len(id) == 0 {
		name = "kafka"
	} else if len(id) > 0 {
		name = id[0]
	}

	v := yago.Component.Ins(name, func() interface{} {

		config := cluster.NewConfig()
		config.Consumer.Return.Errors = true
		config.Group.Return.Notifications = true

		conf := yago.Config.GetStringMap(name)
		brokers, ok := conf["cluster"]
		if !ok || len(brokers.(string)) == 0 {
			log.Fatal("kafka: cluster is empty")
		}
		conn := str.Split(brokers.(string))

		val := NewKafka(conn, config)

		return val
	})

	return v.(*Kafka)
}

// 实例化一个全新的 Kafka
func NewKafka(connect []string, config *cluster.Config) *Kafka {
	return &Kafka{
		connect: connect,
		config:  config,
	}
}

func (q *Kafka) Close() error {
	if q.asyncProducer != nil {
		q.asyncProducer.close()
	}
	if q.syncProducer != nil {
		q.syncProducer.close()
	}
	return nil
}

type Consumer struct {
	conn      *cluster.Consumer
	closeChan chan bool
}

func (q *Kafka) NewConsumer(group string, topics ...string) (*Consumer, error) {
	if len(topics) == 0 {
		return nil, errors.New("kafka: consumer topics can not be empty")
	}

	consumer, err := cluster.NewConsumer(q.connect, group, topics, q.config)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	go func() {
		for ntf := range consumer.Notifications() {
			log.Println("Kafka:", "rebalanced:", ntf)
		}
	}()

	// defer consumer.Close()

	c := new(Consumer)
	c.closeChan = make(chan bool)
	c.conn = consumer

	return c, nil
}

func (c *Consumer) Consume(cb func(topic string, data []byte) bool) error {
	for {
		select {
		case msg, ok := <-c.conn.Messages():
			if ok {
				processed := cb(msg.Topic, msg.Value)
				// mark message as processed
				if processed {
					c.conn.MarkOffset(msg, "")
				}
			}
		case err := <-c.conn.Errors():
			log.Println(err.Error())
		case <-c.closeChan:
			return nil
		}
	}
}

func (c *Consumer) Close() {
	c.closeChan <- true
}

type SyncProducer struct {
	conn sarama.SyncProducer
}

func (q *Kafka) SyncProducer() (*SyncProducer, error) {
	q.mu.Lock()
	defer q.mu.Unlock()
	if q.syncProducer != nil {
		return q.syncProducer, nil
	}

	producer, err := sarama.NewSyncProducer(q.connect, nil)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	var p SyncProducer
	p.conn = producer
	q.syncProducer = &p

	return &p, nil
}

func (p *SyncProducer) Produce(topic, value string) (partition int32, offset int64, err error) {
	msg := &sarama.ProducerMessage{Topic: topic, Value: sarama.StringEncoder(value)}
	partition, offset, err = p.conn.SendMessage(msg)
	if err != nil {
		log.Println(err.Error())
	}
	return partition, offset, err
}

func (p *SyncProducer) close() {
	_ = p.conn.Close()
}

type AsyncProducer struct {
	conn sarama.AsyncProducer
	wg   sync.WaitGroup
}

func (q *Kafka) AsyncProducer() (*AsyncProducer, error) {
	q.mu.Lock()
	defer q.mu.Unlock()

	if q.asyncProducer != nil {
		return q.asyncProducer, nil
	}

	producer, err := sarama.NewAsyncProducer(q.connect, nil)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	var p AsyncProducer
	p.conn = producer

	p.wg.Add(1)
	go func() {
		defer p.wg.Done()
		for range p.conn.Successes() {
		}
	}()

	p.wg.Add(1)
	go func() {
		defer p.wg.Done()
		for err := range p.conn.Errors() {
			log.Println(err)
		}
	}()

	q.asyncProducer = &p

	return &p, nil
}

func (p *AsyncProducer) Produce(topic, value string) {
	msg := &sarama.ProducerMessage{Topic: topic, Value: sarama.StringEncoder(value)}
	p.conn.Input() <- msg
}

func (p *AsyncProducer) close() {
	p.conn.AsyncClose()
	p.wg.Wait()
}
