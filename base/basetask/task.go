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
			ch <- true
		}()

		//handlerFunc()
		select {
		case <-yago.TaskCloseChan:
			return
		case <-ch:
		}

		if intervalOne > 0 {
			time.Sleep(intervalOne)
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
			ch <- handlerFunc()
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

		if intervalOne > 0 {
			time.Sleep(intervalOne)
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
