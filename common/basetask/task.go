package basetask

import (
	"github.com/hulklab/yago"
	"time"
)

type BaseTask struct {
}

func (b BaseTask) RunLoop(handlerFunc yago.TaskHandlerFunc, interval ...time.Duration) {
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

func (b BaseTask) Wait(cb func()) {
	select {
	case <-yago.TaskCloseChan:
		if cb != nil {
			cb()
		}
	}
}
