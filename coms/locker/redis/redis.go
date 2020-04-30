package redis

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/hulklab/yago/libs/str"

	"github.com/hulklab/yago/coms/locker/lock"

	"github.com/garyburd/redigo/redis"
	"github.com/hulklab/yago"
	"github.com/hulklab/yago/coms/rds"
)

const timeDelta = 0

type redisLock struct {
	//rIns    *rds.Rds
	rInsId  string
	retry   int
	key     string
	expired int64
	token   string
	ctx     context.Context
	done    chan struct{}
	errc    chan error
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
		return val
	})
}

func (r *redisLock) rIns() *rds.Rds {
	return rds.Ins(r.rInsId)

}

func (r *redisLock) autoRenewal(ttl int64, errNotify bool) {
	r.done = make(chan struct{})

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

	for i := 0; i < r.retry; i++ {
		err = r.lock(ops.TTL)
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

func (r *redisLock) lock(timeout int64) error {
	i := 1
	token := str.UniqueId()

	for {
		select {
		case <-r.ctx.Done():
			return r.ctx.Err()
		default:

			status, err := redis.String(r.rIns().Do("SET", r.key, token, "EX", timeout, "NX"))
			if err == redis.ErrNil {
				// The lock was not successful, it already exists.
				break
			}

			if err != nil {
				return err
			}

			if status == "OK" {
				r.token = token
				//log.Println("ok:", time.Now())
				return nil
			}

			// 超过 1 分钟归零
			if i >= 60*1000*10 {
				i = 1
			}

			time.Sleep(time.Duration(i*100) * time.Microsecond)

			i++
		}
	}
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

	//if r.errors != nil {
	//	fmt.Println("close err chan")
	//	close(r.errors)
	//	r.errors = nil
	//}

	_, err := unlockScript.Do(r.rIns().GetConn(), r.key, r.token)
	log.Printf("[RedisLock] lock del %s,%v", r.key, err)
}

func (r *redisLock) ErrC() <-chan error {
	return r.errc
}
