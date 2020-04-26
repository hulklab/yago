package basetask

import (
	"time"

	"github.com/hulklab/yago"
)

type BaseTask struct{}

func (b BaseTask) RunLoop(handlerFunc func(), interval ...time.Duration) {
	var intervalOne time.Duration
	if len(interval) > 0 {
		intervalOne = interval[0]
	}
	for {
		ch := make(chan bool)

		go func() {
			handlerFunc()
			if intervalOne > 0 {
				time.Sleep(intervalOne)
			}
			ch <- true
		}()

		//handlerFunc()
		select {
		case <-yago.TaskCloseChan:
			return
		case <-ch:
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
		ch := make(chan bool)

		go func() {
			b := handlerFunc()

			if intervalOne > 0 {
				time.Sleep(intervalOne)
			}

			ch <- b
		}()

		//if !handlerFunc() {
		//	return
		//}
		select {
		case <-yago.TaskCloseChan:
			return
		case b := <-ch:
			if !b {
				return
			}
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
