package logger

import (
	"testing"
	"time"

	"github.com/hulklab/yago"
	"github.com/sirupsen/logrus"
)

// go test -v ./coms/logger -bench="."  -args "-c=${PWD}/example/conf/app.toml"

func TestPlainText(t *testing.T) {
	Ins().Info("this is a info level msg")
	Ins().Warn("this is a warn level msg")
	Ins().Error("this is a error level msg")
	// Ins().Panic("this is a panic level msg")
	// Ins().Fatal("this is a fatal level msg")
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
	Ins().WithFields(logrus.Fields{
		"hello": "logger",
	}).Info()
}

func BenchmarkFile(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Ins().WithFields(logrus.Fields{
			"hello": "logger",
		}).Info()
	}
}
