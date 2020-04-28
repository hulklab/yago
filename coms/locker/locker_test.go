package locker

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/hulklab/yago/coms/locker/lock"

	"github.com/hulklab/yago/example/app/g"

	"github.com/hulklab/yago"
	_ "github.com/hulklab/yago/coms/locker/etcd"
)

// go test -v . -args "-c=${APP_PATH}/app.toml"

func initRedis() {
	yago.Config.Set("locker", g.Hash{
		"driver":             "redis",
		"driver_instance_id": "redis",
	})
	yago.Config.Set("redis", g.Hash{
		"addr": "127.0.0.1:6379",
	})

}

func initEtcd() {
	yago.Config.Set("locker", g.Hash{
		"driver":             "etcd",
		"driver_instance_id": "etcd",
	})
	yago.Config.Set("etcd", g.Hash{
		"endpoints": []string{"127.0.0.1:2379"},
	})

}
func TestRedisExpire(t *testing.T) {
	//t.Skip()
	initRedis()
	doTest()
}

func TestRedisForever(t *testing.T) {
	initRedis()
	doTestForever()
}

func TestRedisWaitTime(t *testing.T) {
	//t.Skip()
	initRedis()
	doTestWaitTime()
}

func TestEtcd(t *testing.T) {
	initEtcd()

	doTest()
}

func TestEtcdForever(t *testing.T) {
	initEtcd()

	doTestForever()
}

func TestEtcdWaitTime(t *testing.T) {
	initEtcd()

	doTestWaitTime()
}

func doTest() {
	key := "lock_test"

	go func() {
		r1 := New()
		err := r1.Lock(key, lock.WithTTL(5), lock.WithDisableKeepAlive())
		if err != nil {
			fmt.Println("get lock in fun1 err", err.Error())
			return
		}
		defer r1.Unlock()
		fmt.Println("get lock in fun1")
		for i := 0; i < 10; i++ {
			fmt.Println("fun1:", i)
			time.Sleep(1 * time.Second)
		}
	}()

	go func() {
		r2 := New()
		err := r2.Lock(key, lock.WithTTL(5), lock.WithDisableKeepAlive())
		if err != nil {
			fmt.Println("get lock in fun2 err", err.Error())
			return
		}
		defer r2.Unlock()
		fmt.Println("get lock in fun2")
		for i := 0; i < 10; i++ {
			fmt.Println("fun2:", i)
			time.Sleep(1 * time.Second)
		}
	}()

	time.Sleep(12 * time.Second)

	r3 := New()
	err := r3.Lock(key, lock.WithTTL(10))
	if err != nil {
		fmt.Println("get lock in fun3 err:", err.Error())
		return
	}

	fmt.Println("get lock in fun3")
	r3.Unlock()

}

func doTestForever() {
	key := "lock_test_forever"

	go func() {
		r1 := New()
		err := r1.Lock(key, lock.WithTTL(1))
		if err != nil {
			fmt.Println("get lock in forever fun1 err", err.Error())
			return
		}
		defer r1.Unlock()
		fmt.Println("get forever lock in fun1")
		for i := 0; i < 10; i++ {
			fmt.Println("forever fun1:", i)
			time.Sleep(1 * time.Second)
		}
	}()

	time.Sleep(2 * time.Second)
	r2 := New()
	err := r2.Lock(key, lock.WithTTL(12))
	if err != nil {
		fmt.Println("get forever lock in fun2 err:", err.Error())
		return
	}

	fmt.Println("get forever lock in fun2")
	r2.Unlock()
}

func doTestWaitTime() {
	key := "lock_test_wait_time"

	go func() {
		r1 := New()
		err := r1.Lock(key, lock.WithWaitTime(time.Second*10))
		if err != nil {
			fmt.Println("get wait time lock in fun1 err", err.Error())
			return
		}
		defer r1.Unlock()
		fmt.Println("get wait time lock in fun1")
		for i := 0; i < 10; i++ {
			fmt.Println("fun1:", i)
			time.Sleep(1 * time.Second)
		}
	}()

	time.Sleep(2 * time.Second)
	r2 := New()
	err := r2.Lock(key, lock.WithWaitTime(time.Second*2))
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			fmt.Println("wait time out and then return")
			return
		}
		fmt.Println("get wait time lock in fun2 err:", err.Error())
		return
	}

	fmt.Println("get wait time lock in fun2")
	r2.Unlock()
}
