package trace

import (
	"context"

	"github.com/hulklab/yago/coms/logger"
	"github.com/hulklab/yago/example/app/g"
	"github.com/sirupsen/logrus"
)

const (
	traceIdKey = "__trace_id__"
)

type Context struct {
	context.Context
}

func New() *Context {
	c := &Context{
		Context: context.Background(),
	}
	return c
}

func NewWithCtx(ctx context.Context) *Context {
	c := &Context{
		Context: ctx,
	}
	return c
}

func (t *Context) Set(key string, val interface{}) {
	t.Context = context.WithValue(t.Context, key, val)
}

func (t *Context) Get(key string) interface{} {
	return t.Context.Value(key)
}

func (t *Context) GetString(key string) string {
	val := t.Get(key)
	if stringVal, ok := val.(string); ok {
		return stringVal
	}
	return ""
}

func (t *Context) GetTraceId() string {
	return t.GetString(traceIdKey)
}

func (t *Context) SetTraceId(traceId string) {
	t.Set(traceIdKey, traceId)
}

func (t *Context) Logger() *logrus.Entry {
	field := g.Hash{
		"trace_id": t.GetTraceId(),
	}

	return logger.Ins().WithFields(field)
}
