package mqtt

import (
	"github.com/hulklab/yago/libs/logger"
	"os"
	"os/signal"
	"syscall"
	"testing"
)

// go test -v ./app/libs/mqtt -test.run Test* -args "-c=${PWD}/app.toml"

func TestPub(t *testing.T) {
	// 注意， pub 和 sub 不同使用同一个 mqtt 实例
	logger.Ins().Info("PUB HERE")
	err := Ins().Pub("event_topic", "hello,world", 0)
	t.Log(err)
}

func TestSub(t *testing.T) {

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGTERM, syscall.SIGINT)

	err := Ins().Sub("event_topic", 0, func(value string) {
		logger.Ins().Info(value, "here")
	})

	if err != nil {
		return
	}

	<-signals
}
