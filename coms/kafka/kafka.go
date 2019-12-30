package kafka

import (
	"log"
	"sync"

	"github.com/Shopify/sarama"
	cluster "github.com/bsm/sarama-cluster"
	"github.com/hulklab/yago"
	"github.com/hulklab/yago/libs/str"
)

type Kafka struct {
	connect []string
	config  *cluster.Config
	topic   string

	syncProducers  sync.Map
	asyncProducers sync.Map
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

		topic, ok := conf["topic"]
		if !ok || len(topic.(string)) == 0 {
			log.Fatal("kafka: default topic is empty")
		}
		val := NewKafka(conn, config, topic.(string))

		return val
	})

	return v.(*Kafka)
}

// 实例化一个全新的 Kafka
func NewKafka(connect []string, config *cluster.Config, topic string) *Kafka {
	return &Kafka{
		connect: connect,
		config:  config,
		topic:   topic,
	}
}

type Consumer struct {
	conn      *cluster.Consumer
	closeChan chan bool
}

func (q *Kafka) NewConsumer(group string, topic ...string) (*Consumer, error) {
	var topics []string
	if len(topic) == 0 {
		topics = []string{q.topic}
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

	//defer consumer.Close()

	c := new(Consumer)
	c.closeChan = make(chan bool)
	c.conn = consumer

	return c, nil
}

func (c *Consumer) Consume(cb func([]byte)) error {
	for {
		select {
		case msg, ok := <-c.conn.Messages():
			if ok {
				// 回调函数处理消息
				cb(msg.Value)
				// mark message as processed
				c.conn.MarkOffset(msg, "")
			}
		case err := <-c.conn.Errors():
			log.Println(err.Error())
			return err
		case <-c.closeChan:
			return nil
		}
	}
}

func (c *Consumer) Close() {
	c.closeChan <- true
}

type SyncProducer struct {
	topic string
	conn  sarama.SyncProducer
}

func (q *Kafka) SyncProducer(topic ...string) (*SyncProducer, error) {
	var p SyncProducer
	if len(topic) == 0 {
		p.topic = q.topic
	} else {
		p.topic = topic[0]
	}

	v, ok := q.syncProducers.Load(p.topic)
	if ok {
		return v.(*SyncProducer), nil
	}

	producer, err := sarama.NewSyncProducer(q.connect, nil)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	p.conn = producer

	q.syncProducers.LoadOrStore(p.topic, &p)

	return &p, nil
}

func (p *SyncProducer) Produce(value string) (partition int32, offset int64, err error) {
	msg := &sarama.ProducerMessage{Topic: p.topic, Value: sarama.StringEncoder(value)}
	partition, offset, err = p.conn.SendMessage(msg)
	if err != nil {
		log.Println(err.Error())
	}
	return partition, offset, err
}

func (p *SyncProducer) Close() {
	_ = p.conn.Close()
}

type AsyncProducer struct {
	topic string
	conn  sarama.AsyncProducer
	wg    sync.WaitGroup

	successes int64
	errors    int64
	enqueued  int64
}

func (q *Kafka) AsyncProducer(topic ...string) (*AsyncProducer, error) {
	var p AsyncProducer
	if len(topic) == 0 {
		p.topic = q.topic
	} else {
		p.topic = topic[0]
	}

	v, ok := q.asyncProducers.Load(p.topic)
	if ok {
		return v.(*AsyncProducer), nil
	}

	producer, err := sarama.NewAsyncProducer(q.connect, nil)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	p.wg.Add(2)
	go func() {
		defer p.wg.Done()
		for range p.conn.Successes() {
			p.successes++
		}
	}()

	go func() {
		defer p.wg.Done()
		for err := range p.conn.Errors() {
			log.Printf("Kafka error: %+v, err：%s\n", err.Msg, err.Err)
			p.errors++
		}
	}()
	p.conn = producer

	q.asyncProducers.LoadOrStore(p.topic, &p)

	return &p, nil
}

func (p *AsyncProducer) Produce(value string) {
	msg := &sarama.ProducerMessage{Topic: p.topic, Value: sarama.StringEncoder(value)}
	p.conn.Input() <- msg
	p.enqueued++
}

func (p *AsyncProducer) Close() {
	p.conn.AsyncClose()
	p.wg.Wait()
}
