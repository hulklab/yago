package demotask

import (
	"log"
	"time"

	"github.com/hulklab/yago"
	"github.com/hulklab/yago/base/basetask"
)

type DemoTask struct {
	basetask.BaseTask
}

func init() {
	t := new(DemoTask)
	yago.AddTaskRouter("@loop", t.HelloLoopAction)
	yago.AddTaskRouter("0 */1 * * * *", t.HelloSchduleAction)
}

func (t *DemoTask) HelloLoopAction() {
	t.RunLoop(func() {
		log.Println("Start Task homeTask.HelloLoopAction")
		log.Println("Doing Task homeTask.HelloLoopAction")
		time.Sleep(time.Second * time.Duration(5))
		log.Println("End Task homeTask.HelloLoopAction")
	})
}

func (t *DemoTask) HelloSchduleAction() {
	log.Println("Start Task homeTask.HelloSchduleAction")
	log.Println("Doing Task homeTask.HelloSchduleAction")
	time.Sleep(time.Second * time.Duration(1))
	log.Println("End Task homeTask.HelloSchduleAction")
}
