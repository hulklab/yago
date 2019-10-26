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

func (this *semaphore) TryAcquire() bool {
	select {
	case this.channel <- struct{}{}:
		this.wg.Add(1)
		return true
	default:
		return false
	}
}

func (this *semaphore) Acquire() {
	this.channel <- struct{}{}
	this.wg.Add(1)
}

func (this *semaphore) Release() {
	<-this.channel
	this.wg.Done()
}

func (this *semaphore) Wait() {
	this.wg.Wait()
}

func (this *semaphore) AvailablePermits() int {
	return this.bufSize - len(this.channel)
}
