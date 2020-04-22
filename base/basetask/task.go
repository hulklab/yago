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
		handlerFunc()
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
