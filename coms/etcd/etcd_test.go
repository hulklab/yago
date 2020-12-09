package etcd

import (
	"context"
	"fmt"
	"testing"
	"time"

	"go.etcd.io/etcd/clientv3"
)

// go test -v ./coms/etcd -run TestPut -args "-c=${PWD}/example/conf/app.toml"
func TestPut(t *testing.T) {
	// put without requestTimeout
	key1, value1 := "key1", "value1"
	_, err := Ins().Put(context.TODO(), key1, value1)
	if err != nil {
		t.Fatal("put err", err)
	}

	// put with requestTimeout
	key2, value2 := "key2", "value2"
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	_, err = Ins().Put(ctx, key2, value2)
	cancel()
	if err != nil {
		t.Fatal("put err with requestTime", err)
	}
}

// go test -v ./coms/etcd -run TestGet -args "-c=${PWD}/example/conf/app.toml"
func TestGet(t *testing.T) {
	// get key
	key := "key1"
	res, err := Ins().Get(context.TODO(), key, clientv3.WithLimit(3))
	if err != nil {
		t.Fatal("get key err", err)
	}

	if len(res.Kvs) == 0 {
		fmt.Println("key is not exist")
	}

	for _, item := range res.Kvs {
		fmt.Println(string(item.Key), "=>", string(item.Value))
	}

	// get with prefix by desc order
	prefix := "k"
	res, err = Ins().Get(context.TODO(), prefix, clientv3.WithPrefix(), clientv3.WithSort(clientv3.SortByKey, clientv3.SortDescend))
	if err != nil {
		t.Fatal("get key err", err)
	}

	if len(res.Kvs) == 0 {
		fmt.Println("the prefix key is not exist")
	}

	for _, item := range res.Kvs {
		fmt.Println(string(item.Key), "=>", string(item.Value))
	}
}

// go test -v ./coms/etcd -run TestDel -args "-c=${PWD}/example/conf/app.toml"
func TestDel(t *testing.T) {
	key := "key1"
	_, err := Ins().Delete(context.TODO(), key)
	if err != nil {
		t.Fatal("del err", err)
	}
}

// go test -v ./coms/etcd -run TestEtcdWatch -args "-c=${PWD}/example/conf/app.toml"
func TestEtcdWatch(t *testing.T) {
	etcd := Ins()

	putTimes := 0
	delTimes := 0

	watchPutTs := 0
	WatchDelTs := 0

	watchChan := etcd.Watch(context.TODO(), "", clientv3.WithPrefix()) // watch all key
	go func() {
		for {
			msg := <-watchChan
			for _, event := range msg.Events {
				if event.Type == clientv3.EventTypePut {
					watchPutTs += 1
					fmt.Println("watch:", string(event.Kv.Key), "=>put=> ", string(event.Kv.Value))
				} else if event.Type == clientv3.EventTypeDelete {
					WatchDelTs += 1
					fmt.Println("watch:", string(event.Kv.Key), "=>delete=> ", string(event.Kv.Value))
				}

			}
		}
	}()

	go func() {
		_, err := etcd.Put(context.TODO(), "username", "dq")
		if err == nil {
			putTimes += 1
		}

		_, err = etcd.Put(context.TODO(), "age", "18")
		if err == nil {
			putTimes += 1
		}

		_, err = etcd.Delete(context.TODO(), "username")
		if err == nil {
			delTimes += 1
		}

		_, err = etcd.Delete(context.TODO(), "age")
		if err == nil {
			delTimes += 1
		}
	}()

	time.Sleep(time.Second * 10)

	if putTimes != watchPutTs {
		t.Fatal("put times not equal watch put times")
	}

	if delTimes != WatchDelTs {
		t.Fatal("del times not equal watch del times")
	}
}
