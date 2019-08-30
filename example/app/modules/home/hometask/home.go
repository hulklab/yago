package hometask

import (
	"github.com/hulklab/yago"
	"github.com/hulklab/yago/base/basetask"
	"log"
	"time"
)

type HomeTask struct {
	basetask.BaseTask
}

func init() {
	homeTask := new(HomeTask)
	yago.AddTaskRouter("@loop", homeTask.HelloAction)
	// yago.AddTaskRouter("0 * * * * *", homeTask.HelloAction)
	// yago.AddTaskRouter("0 * * * * *", homeTask.HelloAction)
}

func (t *HomeTask) HelloAction() {
	t.RunLoop(func() {
		log.Println("start task")
		log.Println("doing")
		time.Sleep(time.Second * time.Duration(5))
		log.Println("end task")
	})
}
