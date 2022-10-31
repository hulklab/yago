package gmodel

import (
	"context"

	"github.com/hulklab/yago/coms/orm"
	"xorm.io/xorm"
)

type Arg struct {
	session *xorm.Session
	ctx     context.Context
	db      string
}

type Option func(arg *Arg)

func WithDb(db string) Option {
	return func(arg *Arg) {
		arg.db = db
	}
}

func WithSession(session *xorm.Session) Option {
	return func(arg *Arg) {
		arg.session = session
	}
}

func WithCtx(ctx context.Context) Option {
	return func(o *Arg) {
		o.ctx = ctx
	}
}

// WithCtxAndSess Ctx and Session
func WithCtxAndSess(ctx context.Context, session *xorm.Session) Option {
	opt := func(o *Arg) {
		o.ctx = ctx
		o.session = session
	}

	return opt
}

// 处理 model 里面的事务问题
type BaseModel struct {
	arg Arg
}

func (m *BaseModel) Init(opts ...Option) {
	arg := Arg{}

	for _, opt := range opts {
		opt(&arg)
	}

	m.arg = arg
}

func (m *BaseModel) GetSession() *xorm.Session {
	if m.arg.session == nil {
		ctx := context.Background()
		if m.arg.ctx != nil {
			ctx = m.arg.ctx
		}

		id := make([]string, 0)
		if len(m.arg.db) > 0 {
			id = append(id, m.arg.db)
		}

		return orm.Ins(id...).Context(ctx)
	}

	if m.arg.ctx != nil {
		return m.arg.session.Context(m.arg.ctx)
	}

	return m.arg.session
}

func (m *BaseModel) UpdateById(id int64, bean interface{}, cols ...string) error {
	session := m.GetSession()

	if len(cols) > 0 {
		session.Cols(cols...)
	}

	_, err := session.Where("id=?", id).Update(bean)
	return err
}
