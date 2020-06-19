package semalib

import (
	"context"
	"sync"
	"sync/atomic"
	"unsafe"
)

type semaphore struct {
	bufSize int
	channel chan struct{}
	wg      *sync.WaitGroup

	error   unsafe.Pointer
	lock sync.Mutex
	ctx    context.Context
	cancel context.CancelFunc
}

func New(concurrencyNum int) *semaphore {
	s := new(semaphore)

	s.channel = make(chan struct{}, concurrencyNum)
	s.bufSize = concurrencyNum
	s.wg = &sync.WaitGroup{}
	s.ctx, s.cancel = context.WithCancel(context.Background())

	return s
}

func (s *semaphore) TryAcquire() bool {
	select {
	case s.channel <- struct{}{}:
		s.wg.Add(1)
		return true
	default:
		return false
	}
}

func (s *semaphore) Acquire() {
	s.channel <- struct{}{}
	s.wg.Add(1)
}

func (s *semaphore) Release() {
	<-s.channel
	s.wg.Done()
}

// add goruntine
// trigger error returns early for concurrency
func (s *semaphore) Add(f func() error) {
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()

		select {
		case <-s.ctx.Done():
			return
		default:
			err := f()

			if err != nil {
				if ! atomic.CompareAndSwapPointer(&s.error, nil, unsafe.Pointer(&err)) {
					return
				}

				s.cancel()
			}

			return
		}
	}()
}

func (s *semaphore) Wait() error {
	s.wg.Wait()

	err := (*error)(atomic.LoadPointer(&s.error))
	if err != nil {
		return *err
	}

	return nil
}

func (s *semaphore) AvailablePermits() int {
	return s.bufSize - len(s.channel)
}
