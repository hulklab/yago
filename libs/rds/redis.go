package rds

import (
	"github.com/garyburd/redigo/redis"
	"github.com/hulklab/yago"
	"github.com/pkg/errors"
	"log"
	"time"
)

type Rds struct {
	redis.Conn
}

// 返回 redis 的一个连接
func Ins(id ...string) *Rds {

	var name string

	if len(id) == 0 {
		name = "redis"
	} else if len(id) > 0 {
		name = id[0]
	}

	v := yago.Component.Ins(name, func() interface{} {

		val := initRedisConnPool(name)
		return val
	})

	redisPool := v.(*redis.Pool)

	rds := redisPool.Get()
	return &Rds{Conn: rds}
}

func initRedisConnPool(name string) *redis.Pool {

	config := yago.Config.GetStringMap(name)

	addr := config["addr"].(string)

	var maxIdle = 5
	mIdle, ok := config["max_idle"]
	if ok {
		maxIdle = int(mIdle.(int64))
	}

	var idleTimeout = time.Duration(240) * time.Second
	iTimeout, ok := config["idle_timeout"]
	if ok {
		idleTimeout = time.Duration(iTimeout.(int64)) * time.Second
	}

	var dialOptions = make([]redis.DialOption, 0)

	connTimeout, ok := config["conn_timeout"]
	if ok {
		ct := time.Duration(connTimeout.(int64)) * time.Millisecond
		dialOptions = append(dialOptions, redis.DialConnectTimeout(ct))
	}

	readTimeout, ok := config["read_timeout"]
	if ok {
		rt := time.Duration(readTimeout.(int64)) * time.Millisecond
		dialOptions = append(dialOptions, redis.DialReadTimeout(rt))
	}

	writeTimeout, ok := config["write_timeout"]
	if ok {
		wt := time.Duration(writeTimeout.(int64)) * time.Millisecond
		dialOptions = append(dialOptions, redis.DialWriteTimeout(wt))
	}

	passwd, ok := config["auth"]
	if ok {
		dialOptions = append(dialOptions, redis.DialPassword(passwd.(string)))
	}

	db, ok := config["db"]
	if ok {
		dialOptions = append(dialOptions, redis.DialDatabase(int(db.(int64))))
	}

	return &redis.Pool{
		MaxIdle:     maxIdle,
		IdleTimeout: idleTimeout,
		//MaxActive:   maxActive,
		//Wait:        true,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", addr, dialOptions...)
			if err != nil {
				return nil, err
			}
			return c, err
		},
		TestOnBorrow: pingRedis,
	}
}

// @todo log
func pingRedis(c redis.Conn, t time.Time) error {
	_, err := c.Do("ping")
	if err != nil {
		log.Println("[ERROR] ping redis fail", err)
	}
	return err
}

type Subscriber struct {
	topic     string
	conn      *redis.PubSubConn
	closeChan chan bool
}

func NewSubscriber(conn redis.Conn, topic string) (*Subscriber, error) {
	s := new(Subscriber)
	s.closeChan = make(chan bool)
	s.topic = topic
	prc := redis.PubSubConn{Conn: conn}
	err := prc.Subscribe(s.topic)
	if err != nil {
		log.Println("redis: ", err.Error())
		return nil, err
	}
	s.conn = &prc
	return s, nil
}

func (s *Subscriber) Subscribe(cb func([]byte)) error {
	for {
		select {
		case <-s.closeChan:
			return nil
		default:
			switch v := s.conn.Receive().(type) {
			case redis.Message:
				cb(v.Data)
			case redis.Subscription:
				if v.Count == 0 {
					s.closeChan <- true
				}
			case error:
				return errors.New(v.Error())
			}
		}
	}
}

func (s *Subscriber) Close() {
	s.closeChan <- true
}
