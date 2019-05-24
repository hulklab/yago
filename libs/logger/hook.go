package logger

import (
	"github.com/sirupsen/logrus"
)

type Hook struct {
	fields logrus.Fields
}

func (h *Hook) SetFields(kv logrus.Fields) {
	h.fields = kv
}

func (h *Hook) Fire(entry *logrus.Entry) error {

	for k, v := range h.fields {
		entry.Data[k] = v
	}

	return nil
}

func (h *Hook) Levels() []logrus.Level {
	return logrus.AllLevels
}
