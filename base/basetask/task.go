package basetask

import (
	"log"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/hulklab/yago"
	"github.com/hulklab/yago/coms/locker"
	"github.com/hulklab/yago/coms/locker/lock"
)

type BaseTask struct{}

func (b BaseTask) RunLoop(handlerFunc func(), interval ...time.Duration) {
	var intervalOne time.Duration
	if len(interval) > 0 {
		intervalOne = interval[0]
	}

	handlerFunc()

	for {
		select {
		case <-yago.StopChan:
			return
		case <-time.After(intervalOne):
			handlerFunc()
		}
	}
}

type LoopArg struct {
	Interval     time.Duration
	LockKey      string
	LockConf     string
	LockInterval time.Duration
	LockTTL      int64
}

type RunLoopOption func(arg *LoopArg)

func WithInterval(interval time.Duration) RunLoopOption {
	return func(arg *LoopArg) {
		arg.Interval = interval
	}
}

func WithLockKey(key string) RunLoopOption {
	return func(arg *LoopArg) {
		arg.LockKey = key
	}
}

func WithLockConf(confName string) RunLoopOption {
	return func(arg *LoopArg) {
		arg.LockConf = confName
	}
}

func WithLockInterval(interval time.Duration) RunLoopOption {
	return func(arg *LoopArg) {
		arg.LockInterval = interval
	}
}

func WithLockTTL(ttl int64) RunLoopOption {
	return func(arg *LoopArg) {
		arg.LockTTL = ttl
	}
}

func (b BaseTask) RunLoopWithLock(handlerFunc func(), opts ...RunLoopOption) {
	var localLocker sync.Mutex
	wrapHandler := func() {
		localLocker.Lock()
		defer localLocker.Unlock()

		handlerFunc()
	}

	arg := LoopArg{}
	for _, opt := range opts {
		opt(&arg)
	}

	if len(arg.LockKey) == 0 {
		pcs := make([]uintptr, 2)
		_ = runtime.Callers(1, pcs)

		key := runtime.FuncForPC(pcs[1]).Name()
		arg.LockKey = strings.NewReplacer("(", "", ")", "", "*", "").Replace(key)
	}

	if len(arg.LockConf) == 0 {
		arg.LockConf = "locker"
	}

	if arg.LockInterval <= 0 {
		arg.LockInterval = time.Second
	}

HEAVEN:
	log.Printf("[RunLoopWithLock] args: %+v", arg)

	mu := locker.New(arg.LockConf)

	for {
		ch := make(chan bool)

		go func(mu lock.ILocker, ch chan bool) {

			err := mu.Lock(arg.LockKey, lock.WithTTL(arg.LockTTL), lock.WithErrorNotify())

			if err == nil {
				ch <- true
			} else {
				time.Sleep(arg.LockInterval)

				ch <- false
				log.Printf("[RunLoopWithLock] err: %s", err.Error())
			}
		}(mu, ch)

		select {
		case <-yago.StopChan:
			return
		case b := <-ch:
			if b {
				goto HELL
			}
		}
	}

HELL:

	log.Println("[RunLoopWithLock] get lock success")

	defer mu.Unlock()

	wrapHandler()

	for {
		select {
		case <-yago.StopChan:
			return
		case err := <-mu.ErrC():
			log.Println("[RunLoopWithLock] some err occur in lock:", err)
			mu.Unlock()
			goto HEAVEN
		case <-time.After(arg.Interval):
			wrapHandler()
		}
	}
}

func (b BaseTask) Wait(cb func()) {
	select {
	case <-yago.StopChan:
		if cb != nil {
			cb()
		}
	}
}
