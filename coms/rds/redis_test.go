package rds

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	"log"
	"testing"
	"time"
)

// go test -v ./coms/rds -test.run TestString -args "-c=${PWD}/app.toml"

func TestApi(t *testing.T) {
	rc := Ins()
	defer rc.Close()
	// 测试 keys
	b, err := redis.Bool(rc.Exists("zjl_test_001"))
	if err != nil {
		t.Error("test exists err:", err.Error())
	} else {
		if b == false {
			t.Log("test exists ok")
		} else {
			t.Error("test exists fail")

		}
	}

	// 测试 strings
	s, err := redis.String(rc.Set("zjl_test_string", "hello", "ex", 10))
	if err != nil {
		t.Error("test strings set err:", err.Error())
	} else {
		t.Log("test set ok ", s)
	}

	s2, err := redis.String(rc.Get("zjl_test_string"))
	if err != nil {
		t.Error("test strings get err:", err.Error())
	} else {
		t.Log("test get ok ", s2)
	}

	// 测试 list
	i, err := redis.Int(rc.RPush("zjl_test_list", "zjl"))
	if err != nil {
		t.Error("test rpush err:", err.Error())
	} else {
		t.Log("test rpush ok ", i)
	}

	ls, err := redis.String(rc.LPop("zjl_test_list"))
	if err != nil {
		t.Error("test lpop err:", err.Error())
	} else {
		t.Log("test lpop ok ", ls)
	}

	// 测试 hash
	hi, err := redis.Int(rc.HSet("zjl_test_hash", "name", []byte("zjl")))
	if err != nil {
		t.Error("test hset err:", err.Error())
	} else {
		t.Log("test hset ok ", hi)
	}

	hbs, err := redis.Bytes(rc.HGet("zjl_test_hash", "name"))
	if err != nil {
		t.Error("test hget err:", err.Error())
	} else {
		t.Log("test hget ok ", hbs)
	}

	hsm, err := redis.StringMap(rc.HGetAll("zjl_test_hash"))
	if err != nil {
		t.Error("test hgetall err:", err.Error())
	} else {
		t.Logf("test hgetall ok %T, %v", hsm, hsm)
	}

	// 测试 set
	si, err := redis.Int(rc.SAdd("zjl_test_set", "zjl"))
	if err != nil {
		t.Error("test sadd err:", err.Error())
	} else {
		t.Log("test sadd ok", si)
	}

	ssl, err := redis.Strings(rc.SMembers("zjl_test_set"))
	if err != nil {
		t.Error("test smembers err:", err.Error())
	} else {
		t.Logf("test smembers ok %T, %v", ssl, ssl)
	}

	// 测试 order set
	zi, err := redis.Int(rc.ZAdd("zjl_test_zset", 1, "php", 2, "go", 3, "python"))
	if err != nil {
		t.Error("test zadd err:", err.Error())
	} else {
		t.Log("test zadd ok", zi)
	}

	zsm, err := redis.StringMap(rc.ZRange("zjl_test_zset", 0, -1, "WITHSCORES"))
	if err != nil {
		t.Error("test zrang err:", err.Error())
	} else {
		t.Logf("test zrang ok, %T, %v", zsm, zsm)
	}

	// 删除 key
	li, err := redis.Int(rc.Del("zjl_test_list", "zjl_test_string", "zjl_test_001", "zjl_test_hash"))
	if err != nil {
		t.Error("test del err:", err.Error())
	} else {
		t.Log("test del ok ", li)
	}
}

func TestString(t *testing.T) {

	rc := Ins()

	// 用完后返回连接池
	defer rc.Close()

	reply, err := rc.Do("set", "test_key", "zjl", "NX")
	fmt.Printf("%T,%v,%T,%s\n", reply, reply, err, err)

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

	subscriber, err := Ins().NewSubscriber(topic)
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
		subscriber.Close()
	})

	if err != nil {
		fmt.Println(err.Error())
	}
}
