package rds

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	"github.com/hulklab/yago"
	"log"
	"testing"
	"time"
)

// go test -v ./app/libs/rds -test.run TestString -args "-c=${PWD}/app.toml"

func TestString(t *testing.T) {
	yago.Config = yago.NewAppConfig("/Users/zhangjiulong/projects/go/src/github.com/hulklab/yago/example/conf/app.toml")

	rc := Ins()

	// 用完后返回连接池
	defer rc.Close()

	rc.Do("set", "test_key", "zjl")

	v, err := redis.String(rc.Do("get", "test_key"))

	if v != "zjl" {
		t.Error("test write fail", v)
	} else {
		t.Log("test write ok", v)
	}

	// 删除
	rc.Do("del", "test_key")

	v, err = redis.String(rc.Do("get", "test_key"))
	if err == redis.ErrNil {
		t.Log("test delete ok", v)
	} else {
		t.Error("test delete fail", v)
	}
}

func TestExpire(t *testing.T) {
	rc := Ins()
	defer rc.Close()

	key := "test_expire_key"

	rc.Do("set", key, "zjl")

	rc.Do("expire", key, 5)

	for i := 0; i < 5; i++ {

		v, _ := redis.String(rc.Do("get", key))
		fmt.Println("v:", v)
		time.Sleep(1 * time.Second)

	}

	v, err := redis.String(rc.Do("get", key))
	if err == redis.ErrNil {
		t.Log("test expire ok,v:", v)
	} else {
		t.Error("test expire fail,v:", v)
	}
}

// 测试列表
func TestList(t *testing.T) {
	rc := Ins()
	defer rc.Close()

	key := "test_list_key"

	rc.Do("lpush", key, "redis")
	rc.Do("lpush", key, "mongo")
	rc.Do("lpush", key, "es")

	v1, _ := redis.String(rc.Do("rpop", key))
	v2, _ := redis.String(rc.Do("rpop", key))
	v3, _ := redis.String(rc.Do("rpop", key))

	if v1 == "redis" && v2 == "mongo" && v3 == "es" {
		t.Log("test list ok,v:", v1, v2, v3)
	} else {
		t.Error("test expire fail,v:", v1, v2, v3)
	}

}

// 测试 hash
func TestMap(t *testing.T) {
	rc := Ins()
	defer rc.Close()

	key := "test_map_key"

	rc.Do("hset", key, "username", "zhangjiulong")
	rc.Do("hmset", key, "name", "zjl", "age", 3)
	rc.Do("hset", key, "phone", "13900000000")

	username, _ := redis.String(rc.Do("hget", key, "username"))
	fmt.Println(username)

	v, _ := redis.Values(rc.Do("hmget", key, "name", "age", "phone"))
	for _, val := range v {
		fmt.Println(string(val.([]byte)))
	}
}

// 测试 set
func TestSet(t *testing.T) {
	rc := Ins()
	defer rc.Close()

	key := "test_set_key"

	rc.Do("sadd", key, "redis")
	rc.Do("sadd", key, "redis")
	rc.Do("sadd", key, "mongo")

	v, err := redis.Int(rc.Do("scard", key))
	if err != nil {
		t.Error("test set err", err)
	} else if v == 2 {
		t.Log("test set succ", v)
	} else {
		t.Error("test set fail", v)
	}

	b, err := redis.Bool(rc.Do("sismember", key, "redis"))
	if b {
		t.Log("test set succ", b)
	} else {
		t.Error("test set fail", b)

	}
}

// 原生 sub 调用
//func TestPubSub(t *testing.T) {
//
//	rc := Ins()
//	defer rc.Close()
//
//	pubrc := Ins()
//	defer pubrc.Close()
//
//	topic := "test_subscribe"
//
//	var wg sync.WaitGroup
//	wg.Add(2)
//
//	prc := redis.PubSubConn{Conn: rc}
//
//	prc.Subscribe(topic)
//
//	go func() {
//		defer wg.Done()
//
//		for {
//			switch v := prc.Receive().(type) {
//			case redis.Message:
//				fmt.Printf("%s: message: %s\n", v.Channel, v.Data)
//			case redis.Subscription:
//				fmt.Printf("%s: %s %d\n", v.Channel, v.Kind, v.Count)
//				if v.Count == 0 {
//					return
//				}
//			case error:
//				fmt.Println("err:", v)
//				return
//
//			}
//		}
//
//	}()
//
//	go func() {
//		defer wg.Done()
//
//		pubrc.Do("publish", topic, "hello")
//		pubrc.Do("publish", topic, "world")
//
//		prc.Unsubscribe(topic)
//	}()
//
//	wg.Wait()
//}

func TestPubSub(t *testing.T) {
	topic := "test_subscribe"

	subscriber, err := NewSubscriber(Ins(), topic)
	if err != nil {
		log.Println(err.Error())
		return
	}

	go func() {
		r := Ins()
		defer r.Close()
		time.Sleep(time.Second)
		r.Do("publish", topic, "hello")
	}()

	err = subscriber.Subscribe(func(bytes []byte) {
		fmt.Println("msg:", string(bytes))
	})

	if err != nil {
		fmt.Println(err.Error())
	}
}
