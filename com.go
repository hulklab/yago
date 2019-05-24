package yago

import (
	"sync"
)

type Components struct {
	m sync.Map
}

func (c *Components) Ins(key string, f func() interface{}) interface{} {
	// @todo 监听配置重载信号
	v, ok := c.m.Load(key)
	if !ok {
		val := f()
		v, _ = c.m.LoadOrStore(key, val)
	}
	return v

}

var Component = new(Components)

// example @see libs/kafka
