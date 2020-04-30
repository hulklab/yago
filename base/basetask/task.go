package basetask

import (
	"log"
	"runtime"
	"strings"
	"time"

	"github.com/hulklab/yago/coms/locker/lock"

	"github.com/hulklab/yago/coms/locker"

	"github.com/hulklab/yago"
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
		case <-yago.TaskCloseChan:
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
		case <-yago.TaskCloseChan:
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

	handlerFunc()

	for {
		select {
		case <-yago.TaskCloseChan:
			return
		case err := <-mu.ErrC():
			log.Println("[RunLoopWithLock] some err occur in lock:", err)
			mu.Unlock()
			goto HEAVEN
		case <-time.After(arg.Interval):
			handlerFunc()
		}
	}
}

//exit while return false, and continue while return true
func (b BaseTask) RunLoopWhile(handlerFunc func() bool, interval ...time.Duration) {
	var intervalOne time.Duration
	if len(interval) > 0 {
		intervalOne = interval[0]
	}
	for {
		if !handlerFunc() {
			return
		}
		if intervalOne > 0 {
			time.Sleep(intervalOne)
		}
		select {
		case <-yago.TaskCloseChan:
			return
		default:
		}
	}
}

func (b BaseTask) Wait(cb func()) {
	select {
	case <-yago.TaskCloseChan:
		if cb != nil {
			cb()
		}
	}
}
