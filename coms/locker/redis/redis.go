package redis

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/hulklab/yago/libs/str"

	"github.com/hulklab/yago/coms/locker/lock"

	"github.com/garyburd/redigo/redis"
	"github.com/hulklab/yago"
	"github.com/hulklab/yago/coms/rds"
)

//const timeDelta = 0

type redisLock struct {
	rInsId  string
	retry   int
	key     string
	expired int64
	token   string
	ctx     context.Context
	done    chan struct{}
	errc    chan error
	lockc   chan struct{}
	sub *rds.Subscriber
}

func init() {
	lock.RegisterLocker("redis", func(name string) lock.ILocker {
		driverInsId := yago.Config.GetString(name + ".driver_instance_id")
		retry := yago.Config.GetInt(name + ".retry")
		if retry == 0 {
			retry = 3
		}
		//rIns := rds.Ins(driverInsId)
		val := &redisLock{
			rInsId: driverInsId,
			//rIns:   rIns,
			retry: retry,
		}

		val.errc = make(chan error)
		val.lockc = make(chan struct{})
		val.done = make(chan struct{})
		return val
	})
}

func (r *redisLock) rIns() *rds.Rds {
	return rds.Ins(r.rInsId)

}

func (r *redisLock) autoRenewal(ttl int64, errNotify bool) {

	go func() {
		ticker := time.NewTicker(time.Duration(ttl)*time.Second - time.Millisecond*800)
		//log.Println("now:", time.Now())
		defer ticker.Stop()

		for {
			select {
			case <-r.done:
				log.Println("[RedisLock] renewal done")
				return
			case <-ticker.C:

				//log.Println("ticker:", time.Now())
				reply, err := redis.Int64(r.rIns().Expire(r.key, ttl))
				//log.Printf("[Redislock] %d %v", reply, err)
				if err != nil {
					log.Printf("[RedisLock] renewal err: %s\n", err.Error())
					if errNotify {
						go func() {
							r.errc <- fmt.Errorf("%s:%w", "lock renewal err", err)
						}()
					}
					break
				} else if reply == 0 {
					//  续约失败
					log.Printf("[RedisLock] renewal fail: key %s is not exists", r.key)
					if errNotify {
						go func() {
							r.errc <- fmt.Errorf("%s:%w", "lock renewal err", err)
						}()
					}
					break
				}

			}
		}
	}()

}

func (r *redisLock) Lock(key string, opts ...lock.SessionOption) error {
	var ctx context.Context
	ctx = context.Background()

	ops := &lock.SessionOptions{TTL: lock.DefaultSessionTTL}
	for _, opt := range opts {
		opt(ops)
	}

	if ops.WaitTime > 0 {
		var cancelFunc context.CancelFunc

		ctx, cancelFunc = context.WithTimeout(context.Background(), ops.WaitTime)
		defer cancelFunc()
	}

	r.ctx = ctx
	r.key = key

	var err error

	err = r.listen()
	if err != nil {
		return err
	}

	for i := 0; i < r.retry; i++ {
		err = r.tryLock(ops.TTL)
		if err == nil {
			break
		}

		if errors.Is(err, context.DeadlineExceeded) {
			return err
		}

		log.Printf("[RedisLock] lock err:%s,retry:%d", err.Error(), i)
	}

	if err == nil && !ops.DisableKeepAlive {
		r.autoRenewal(ops.TTL, ops.ErrorNotify)
	}

	return err
}

func (r *redisLock) topicKey() string {
	return fmt.Sprintf("__yago_lock_topic_%s", r.key)
}

func (r *redisLock) listen() error {
	sub, err := r.rIns().NewSubscriber(r.topicKey())
	if err != nil {
		return fmt.Errorf("[RedisLock] %s new listener err:%w", r.key, err)
	}

	r.sub = sub

	go func() {
		err := r.sub.Subscribe(func(topic string, bytes []byte) {
			r.lockc <- struct{}{}
		})

		if err != nil {
			log.Println("[RedisLock] redis listen err:", err.Error(), "key:", r.key)
			r.errc <- fmt.Errorf("listen err %w", err)
		}
	}()

	return nil
}

func (r *redisLock) tryLock(timeout int64) error {
	ok, err := r.lock(timeout)
	if err != nil {
		return err
	}

	if ok {
		return nil
	}

	for {
		rn := rand.Int63n(timeout) + 1
		//log.Println(r.key, rn)

		select {
		case <-r.ctx.Done():
			return r.ctx.Err()
		case <-r.lockc:
			//log.Println("[RedisLock] redis get unlock topic:", r.topicKey())
			ok, err := r.lock(timeout)
			if err != nil {
				return err
			}

			if ok {
				return nil
			}
		case <-time.After(time.Second * time.Duration(rn)):
			//log.Printf("[RedisLock] try to get lock %s after time %d", r.key, rn)
			ok, err := r.lock(timeout)
			if err != nil {
				return err
			}

			if ok {
				return nil
			}

		}
	}
}

func (r *redisLock) lock(timeout int64) (ok bool, err error) {
	token := str.UniqueId()

	status, err := redis.String(r.rIns().Do("SET", r.key, token, "EX", timeout, "NX"))
	if err == redis.ErrNil {
		// The lock was not successful, it already exists.
		return false, nil
	}

	if err != nil {
		return false, err
	}

	if status == "OK" {
		r.token = token
		//log.Println("ok:", time.Now())
		return true, nil
	}

	return false, fmt.Errorf("set lock status err:%s", status)
}

var unlockScript = redis.NewScript(1, `
	if redis.call("get", KEYS[1]) == ARGV[1]
	then
		return redis.call("del", KEYS[1])
	else
		return 0
	end
`)

func (r *redisLock) Unlock() {
	if r.done != nil {
		close(r.done)
		r.done = nil
	}

	if r.sub != nil {
		r.sub.Close()
	}

	rc := r.rIns().GetConn()
	defer func(rc redis.Conn) {
		err := rc.Close()
		if err != nil {
			log.Println("[RedisLock] close redis conn err: ", err.Error())
		}
	}(rc)

	reply, err := redis.Int(unlockScript.Do(rc, r.key, r.token))
	log.Printf("[RedisLock] lock del %s,%v,%v", r.key, err, reply)

	if err != nil {
		return
	}

	if reply > 0 {
		_, err = r.rIns().Publish(r.topicKey(), r.token)
		if err != nil {
			log.Printf("[RedisLock] lock %s release and broadcast err %v", r.key, err)
		}
	}

}

func (r *redisLock) ErrC() <-chan error {
	return r.errc
}
