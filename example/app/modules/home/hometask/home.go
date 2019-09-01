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
	yago.AddTaskRouter("@loop", homeTask.HelloLoopAction)
	yago.AddTaskRouter("0 */1 * * * *", homeTask.HelloSchduleAction)
}

func (t *HomeTask) HelloLoopAction() {
	t.RunLoop(func() {
		log.Println("Start Task homeTask.HelloLoopAction")
		log.Println("Doing Task homeTask.HelloLoopAction")
		time.Sleep(time.Second * time.Duration(5))
		log.Println("End Task homeTask.HelloLoopAction")
	})
}

func (t *HomeTask) HelloSchduleAction() {
	log.Println("Start Task homeTask.HelloSchduleAction")
	log.Println("Doing Task homeTask.HelloSchduleAction")
	time.Sleep(time.Second * time.Duration(1))
	log.Println("End Task homeTask.HelloSchduleAction")
}
