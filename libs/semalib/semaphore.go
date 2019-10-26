package semalib

import "sync"

type semaphore struct {
	bufSize int
	channel chan struct{}
	wg      *sync.WaitGroup
}

func New(concurrencyNum int) *semaphore {
	return &semaphore{
		channel: make(chan struct{}, concurrencyNum),
		bufSize: concurrencyNum,
		wg:      &sync.WaitGroup{},
	}
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

func (s *semaphore) Wait() {
	s.wg.Wait()
}

func (s *semaphore) AvailablePermits() int {
	return s.bufSize - len(s.channel)
}
