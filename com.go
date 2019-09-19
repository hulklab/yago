package yago

import (
	"sync"
)

type components struct {
	m sync.Map
}

func (c *components) Ins(key string, f func() interface{}) interface{} {
	// @todo 监听配置重载信号
	v, ok := c.m.Load(key)
	if !ok {
		val := f()
		v, _ = c.m.LoadOrStore(key, val)
	}
	return v

}

func (c *components) Del(key string, cb ...func()) {
	c.m.Delete(key)

	if len(cb) > 0 {
		// 执行回调关闭链接
		go cb[0]()
	}
}

var Component = new(components)

// example @see coms/kafka
