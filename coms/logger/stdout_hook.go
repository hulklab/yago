package logger

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
)

type stdoutHook struct {
}

func NewStdoutHook() *stdoutHook {
	return &stdoutHook{}
}

func (h *stdoutHook) Fire(entry *logrus.Entry) error {
	var code string
	line, err := entry.String()
	if err != nil {
		code = "\033[31m"
		_, _ = fmt.Fprintf(os.Stderr, "\033[0m%sUnable to read entry, %v\033[0m", code, err)
		return err
	}

	fmt.Print(line)
	return nil
}

func (h *stdoutHook) Levels() []logrus.Level {
	return logrus.AllLevels
}
