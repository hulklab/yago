package kafka

import (
	"github.com/Shopify/sarama"
	"github.com/bsm/sarama-cluster"
	"github.com/hulklab/yago"
	"log"
)

type Kafka struct {
	connect string
	config  *cluster.Config
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
		connect, ok := conf["cluster"]
		if !ok {
			// @todo
		}
		val := NewKafka(connect.(string), config)

		return val
	})

	return v.(*Kafka)
}

// 实例化一个全新的 Kafka
func NewKafka(connect string, config *cluster.Config) *Kafka {
	return &Kafka{
		connect: connect,
		config:  config,
	}
}

type Consumer struct {
	conn      *cluster.Consumer
	closeChan chan bool
}

func (q *Kafka) NewConsumer(topic string, group string) (*Consumer, error) {
	consumer, err := cluster.NewConsumer([]string{q.connect}, group, []string{topic}, q.config)
	if err != nil {
		log.Println("Kafka", "init consumer failed", err.Error())
		return nil, err
	}

	go func() {
		for ntf := range consumer.Notifications() {
			log.Println("Kafka", "rebalanced:", ntf)
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
			log.Println("Kafka", err.Error())
			return err
		case <-c.closeChan:
			return nil
		}
	}
}

func (c *Consumer) Close() {
	c.closeChan <- true
}

func (q *Kafka) Produce(topic string, value string) (partition int32, offset int64, err error) {

	producer, err := sarama.NewSyncProducer([]string{q.connect}, nil)
	if err != nil {
		log.Println("Kafka", "init producer failed", err.Error())
		return
	}

	defer func() {
		if err := producer.Close(); err != nil {
			log.Println("Kafka", "close producer failed", err.Error())
		}
	}()

	msg := &sarama.ProducerMessage{Topic: topic, Value: sarama.StringEncoder(value)}

	partition, offset, err = producer.SendMessage(msg)
	if err != nil {
		log.Println("Kafka.produce", "msg", "Failed to send message", "err", err.Error())
	} else {
		log.Println("Kafka.produce", "msg", "Message send success", "partition", partition, "offset", offset)
	}

	return partition, offset, err
}
