package logger

import (
	"github.com/hulklab/yago"
	"github.com/sirupsen/logrus"
	"testing"
	"time"
)

func TestPlainText(t *testing.T) {
	Ins().Info("this is a info level msg")
	Ins().Warn("this is a warn level msg")
	Ins().Error("this is a error level msg")
	//Ins().Panic("this is a panic level msg")
	//Ins().Fatal("this is a fatal level msg")
}

func TestLogFieldText(t *testing.T) {
	Ins().WithFields(logrus.Fields{
		"hello": "logger",
	}).Info("this is a info level msg with fields")
}

func TestLogJson(t *testing.T) {
	Ins().SetFormatter(&logrus.JSONFormatter{})
	Ins().WithFields(logrus.Fields{
		"hello": "logger",
	}).Info()
}

func TestLogHook(t *testing.T) {
	Ins().SetHookFields(logrus.Fields{
		"timestamp": time.Now().Unix(),
		"endpoint":  yago.Hostname(),
	})
	for i := 0; i < 1000000000; i++ {
		Ins().WithFields(logrus.Fields{
			"hello": "logger",
		}).Info()
	}
}
