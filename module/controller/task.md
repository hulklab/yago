# Task 控制器

task 控制器内部调度使用了 [cron](https://github.com/robfig/cron)，这使得 task 控制器可以实现类似 crontab 的定时任务功能。同时在此基础上我们增加了 @loop 关键字，用来支持常驻进程模式，应用场景就是一些异步队列消费永不退出的这种情况。

## 路由注册

task init 函数中通过 AddTaskRouter 完成路由注册。

```go
func init() {
	homeTask := new(HomeTask)
	yago.AddTaskRouter("@loop", homeTask.HelloLoopAction)
	yago.AddTaskRouter("0 */1 * * * *", homeTask.HelloSchduleAction)
}
```

AddTaskRouter 参数说明

| 参数位置 | 参数类型 | 说明 |
| ------- | ------- | ------- |
| 1 | String | 执行任务计划，参考Spec表|
| 2 | Func | 接口对应的 Action Func |

Spec
```bash
# ┌─────────────── second (0 - 59)
# | ┌───────────── minute (0 - 59)
# │ | ┌─────────── hour (0 - 23)
# │ │ | ┌───────── day of the month (1 - 31)
# │ │ │ | ┌─────── month (1 - 12)
# │ │ │ │ | ┌───── day of the week (0 - 6) (Sunday to Saturday;7 is also Sunday on some systems)
# │ │ │ │ │ |                                  
# │ │ │ │ │ |
# │ │ │ │ │ |
# * * * * * * command to execute
```

| Entry | Description | Equivalent to |
| ------ | -------------| --------------|
| @yearly (or @annually) | Run once a year at midnight of 1 January |	0 0 0 1 1 * |
| @monthly | Run once a month at midnight of the first day of the month | 0 0 0 1 * * |
| @weekly | Run once a week at midnight on Sunday morning | 0 0 0 * * 0 |
| @daily (or @midnight) | Run once a day at midnight | 0 0 0 * * * |
| @hourly | Run once an hour at the beginning of the hour | 0 0 * * * * |

还有一个特殊的 @loop，需要注意的是 @loop 必须在 Action 内搭配 RunLoop 函数运行，否则 @loop 只会执行一次便退出。

我们在 RunLoop 内监听了全局关闭信号，用来平滑地完成单次循环，同时还有 Wait 函数用来帮助 task 收到关闭信号时，做一些清理工作。

在 RunLoop 函数内我们传递一个回调函数和一个可选的执行间隔参数，如果执行间隔不传，默认没有等待，直接进入下个loop。

## TaskAction

```go
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

```

## 开启 Task

控制是否需要在此机器上开启 task 任务，有两种方式

* 修改配置文件中的 app.task_enable，默认为开启
* 修改环境变量 export {{配置文件中的app_name}}_APP_TASK_ENABLE=1, 1 表示开启，0 表示关闭，配置文件与环境变量同时存在时环境变量生效