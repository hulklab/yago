package etcd

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/hulklab/yago"
	"github.com/hulklab/yago/coms/etcd"
	"github.com/hulklab/yago/coms/locker/lock"
	"go.etcd.io/etcd/clientv3/concurrency"
)

func init() {
	lock.RegisterLocker("etcd", func(name string) lock.ILocker {
		driverInsId := yago.Config.GetString(name + ".driver_instance_id")
		retry := yago.Config.GetInt(name + ".retry")
		if retry == 0 {
			retry = 3
		}
		// eIns := etcd.Ins(driverInsId)
		val := &etcdLock{
			eInsId: driverInsId,
			retry:  retry,
		}

		val.errc = make(chan error)
		return val
	})
}

type etcdLock struct {
	// eIns  *etcd.Etcd
	eInsId string
	retry  int
	key    string
	ctx    context.Context
	mu     sync.Mutex
	mutex  *concurrency.Mutex
	errc   chan error
}

func (e *etcdLock) eIns() *etcd.Etcd {
	return etcd.Ins(e.eInsId)
}

func (e *etcdLock) Lock(key string, opts ...lock.SessionOption) error {
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

	e.ctx = ctx

	var err error

	for i := 0; i < e.retry; i++ {
		err = e.lock(key, ops.TTL, ops.ErrorNotify)
		if err == nil {
			break
		}

		if errors.Is(err, context.DeadlineExceeded) {
			break
		}

		log.Printf("etcd lock err:%s,retry:%d", err.Error(), i)
	}

	if err != nil {
		return err
	}

	if ops.DisableKeepAlive {
		go func() {
			<-time.After(time.Duration(ops.TTL) * time.Second)
			e.Unlock()
		}()
	}
	return nil
}

func (e *etcdLock) lock(key string, ttl int64, errNotify bool) error {
	response, err := e.eIns().Client.Grant(e.ctx, ttl)
	if err != nil {
		return err
	}

	session, err := concurrency.NewSession(e.eIns().Client, concurrency.WithLease(response.ID))
	if err != nil {
		return err
	}

	mutex := concurrency.NewMutex(session, "/lock_"+key)

	log.Printf("[EtcdLock] lease id: %d\n", response.ID)

	err = mutex.Lock(e.ctx)
	if err != nil {
		return err
	}

	log.Printf("[EtcdLock] lock key: %s\n", mutex.Key())

	// get, err := e.eIns().Get(e.ctx, mutex.Key())
	// log.Printf("get: %v,%s,%s\n", get.Kvs, mutex.Key(), err)

	res, err := e.eIns().Txn(e.ctx).If(mutex.IsOwner()).Commit()
	// log.Printf("if: %v,%s,%s,%+v\n", res.Succeeded, mutex.Key(), err, mutex)

	if err != nil {
		return fmt.Errorf("etcd judge mutex is owner err: %w", err)
	}

	if res.Succeeded == false {
		return errors.New("mutex is not belongs to you")
	}

	e.mutex = mutex

	go func() {
		<-session.Done()
		log.Println("[EtcdLock] session is close")

		if errNotify {
			e.errc <- errors.New("[EtcdLocker] session is close")
		}
	}()

	return nil
}

func (e *etcdLock) Unlock() {
	e.mu.Lock()
	defer e.mu.Unlock()

	if e.mutex != nil {
		log.Println("[EtcdLock] unlock start")
		err := e.mutex.Unlock(e.ctx)
		log.Println("[EtcdLock] unlock err:", err)
		e.mutex = nil
	}
}

func (r *etcdLock) ErrC() <-chan error {
	return r.errc
}
