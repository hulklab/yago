package logger

import (
	"fmt"
	"github.com/hulklab/yago"
	"github.com/natefinch/lumberjack"
	"github.com/sirupsen/logrus"
	"runtime"
	"strings"
)

type Logger struct {
	*logrus.Logger
}

func Ins(id ...string) *Logger {
	var name string

	if len(id) == 0 {
		name = "logger"
	} else if len(id) > 0 {
		name = id[0]
	}

	v := yago.Component.Ins(name, func() interface{} {
		conf := yago.Config.GetStringMap(name)

		formatter := conf["formatter"].(string)
		filePath := conf["file_path"].(string)
		maxSize := int(conf["max_size"].(int64))
		maxBackups := int(conf["max_backups"].(int64))
		maxAge := int(conf["max_age"].(int64))
		level := logrus.Level(conf["level"].(int64))
		compress := conf["compress"].(bool)
		var stdoutEnable bool
		if v, b := conf["stdout_enable"]; b {
			stdoutEnable = v.(bool)
		}

		val := &Logger{logrus.New()}
		// 设置最低log level
		val.SetLevel(level)

		// 日志中显示记录的文件和函数名, 注意：textField 中需要避开 file 和 func 字段
		val.SetReportCaller(true)

		if formatter == "json" {
			val.Formatter = &logrus.JSONFormatter{CallerPrettyfier: CallerPretty}
		} else {
			val.Formatter = &logrus.TextFormatter{CallerPrettyfier: CallerPretty}
		}
		val.Out = &lumberjack.Logger{
			Filename:   filePath,
			MaxSize:    maxSize,
			MaxBackups: maxBackups,
			MaxAge:     maxAge,
			Compress:   compress,
		}

		if stdoutEnable {
			val.AddHook(NewStdoutHook())
		}
		return val
	})

	return v.(*Logger)
}

func (l *Logger) SetHookFields(kv logrus.Fields) {
	hook := new(Hook)

	hook.SetFields(kv)

	l.AddHook(hook)
}

func (l *Logger) Category(c string) *logrus.Entry {
	return l.WithFields(logrus.Fields{
		"category": c,
	})
}

// 将日志中记录的文件名 file 和方法名 func 转成短名字
func CallerPretty(caller *runtime.Frame) (function string, file string) {
	if caller == nil {
		return "", ""
	}

	short := caller.File
	i := strings.LastIndex(caller.File, "/")
	if i != -1 && i != len(caller.File)-1 {
		short = caller.File[i+1:]
	}

	fun := caller.Function
	j := strings.LastIndex(caller.Function, "/")
	if j != -1 && j != len(caller.Function)-1 {
		fun = caller.Function[j+1:]
	}

	return fun, fmt.Sprintf("%s:%d", short, caller.Line)
}
