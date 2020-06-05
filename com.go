package yago

import (
	"io"
	"log"
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

func (c *components) Del(key interface{}, cb ...func()) {
	c.m.Delete(key)

	if len(cb) > 0 {
		// 执行回调关闭链接
		go cb[0]()
	}
}

func (c *components) clear() {
	comKeys := make([]interface{}, 0)
	c.m.Range(func(key, value interface{}) bool {
		comKeys = append(comKeys, key)
		return true
	})

	for _, key := range comKeys {
		c.m.Delete(key)
	}
}

func (c *components) Close() {
	c.m.Range(func(key, value interface{}) bool {
		if v, ok := value.(io.Closer); ok {
			if err := v.Close(); err != nil {
				log.Printf("Com %s close error: %s\n", key, err)
			}
		}
		return true
	})
}

var Component = new(components)
