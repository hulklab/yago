package lock

import (
	"sync"
	"time"
)

const DefaultSessionTTL = 60
const DefaultErrorBufferSize = 4086

type SessionOptions struct {
	TTL              int64
	DisableKeepAlive bool
	WaitTime         time.Duration
	ErrorNotify      bool
}

// SessionOption configures Session.
type SessionOption func(*SessionOptions)

// If TTL is <= 0, the default 60 seconds TTL will be used.
func WithTTL(ttl int64) SessionOption {
	return func(so *SessionOptions) {
		if ttl > 0 {
			so.TTL = ttl
		}
	}
}

func WithDisableKeepAlive() SessionOption {
	return func(so *SessionOptions) {
		so.DisableKeepAlive = true
	}
}

func WithWaitTime(waitTime time.Duration) SessionOption {
	return func(so *SessionOptions) {
		if waitTime > 0 {
			so.WaitTime = waitTime
		}
	}
}

func WithErrorNotify() SessionOption {
	return func(so *SessionOptions) {
		so.ErrorNotify = true
	}
}

var locks sync.Map

type ILocker interface {
	Lock(key string, opts ...SessionOption) error
	Unlock()
	Errors() <-chan error
}

type NewFunc func(configId string) ILocker

func RegisterLocker(name string, f NewFunc) {
	locks.Store(name, f)
}

func LoadLocker(name string) (NewFunc, bool) {
	val, b := locks.Load(name)
	if !b {
		return nil, false
	}

	newFunc, _ := val.(NewFunc)
	return newFunc, true

}
