package basemodel

import (
	"context"
	"errors"

	"github.com/hulklab/yago/coms/orm"
	"xorm.io/builder"
	"xorm.io/xorm"
)

type Options struct {
	session       *xorm.Session
	ctx           context.Context
	defaultPage   int
	defaultSize   int
	defaultOrders []Order
}

type BaseModel struct {
	options *Options
}

func (m *BaseModel) GetSession() *xorm.Session {
	if m.options.session == nil {
		ctx := context.Background()
		if m.options.ctx != nil {
			ctx = m.options.ctx
		}

		return orm.Ins().Context(ctx)
	}

	if m.options.ctx != nil {
		return m.options.session.Context(m.options.ctx)
	}

	return m.options.session
}

type Option func(o *Options)

func (m *BaseModel) Init(opts ...Option) {
	var opt Options
	for _, option := range opts {
		option(&opt)
	}
	m.options = &opt
}

func WithSession(session *xorm.Session) Option {
	return func(o *Options) {
		o.session = session
	}
}

func WithCtx(ctx context.Context) Option {
	return func(o *Options) {
		o.ctx = ctx
	}
}

func WithDefaultPageSize(page, size int) Option {
	return func(o *Options) {
		o.defaultPage = page
		o.defaultSize = size
	}
}

func WithDefaultOrder(orders ...Order) Option {
	return func(o *Options) {
		o.defaultOrders = orders
	}
}

// Funnel screening condition
// For example, "filters":{"user_state":[0,1],"phone_state":[0]} means "where user_state in(0,1) and phone_state = 0"
type List []interface{}

type Filters map[string]List

// Order condition
// 1 = ASC, -1 = DESC
// For example, "orders":[{"created_at":1},{"id":-1}] means "order by created_at asc, id desc"
type Order map[string]int

type Orders []Order

// Fuzzy search
// For example, "sheldon": []string{"username", "phone"} means "username like %sheldon% or phone like %sheldon%"
type Q map[string][]string

type QueryHandler func(session *xorm.Session)

// List query without page
type Query struct {
	Q          Q
	Filters    Filters
	Orders     Orders
	ExtraQuery QueryHandler
}

// List query with page
// For example,
// {
//  "page": 1,
//  "size": 20,
//  "q": "x",
//  "orders": [
//    {
//      "created_at": 1
//    },
//    {
//      "id": -1
//    }
//  ],
//  "filters": {
//    "user_state": [
//      0,
//      1
//    ],
//    "phone_state": [
//      0
//    ]
//  }
// }
type PageQuery struct {
	Page       int
	Size       int
	Q          Q
	Filters    Filters
	Orders     Orders
	ExtraQuery QueryHandler
}

// Build filters condition
func buildFilters(session *xorm.Session, filters Filters) {
	for k, v := range filters {
		if len(v) > 0 {
			session.In(k, v)
		}
	}
}

// Build q condition
func buildQ(session *xorm.Session, q Q) {
	for k, v := range q {
		if len(k) > 0 && len(v) > 0 {
			cond := builder.NewCond()
			for _, f := range v {
				cond = cond.Or(builder.Like{f, k})
			}
			session.And(cond)
		}
	}
}

// Build order condition
func buildOrders(session *xorm.Session, orders Orders) {
	for _, v := range orders {
		for f, s := range v {
			if s == 1 {
				session.Asc(f)
			} else if s == -1 {
				session.Desc(f)
			}
		}
	}
}

// List with page
func (m *BaseModel) PageList(query *PageQuery, items interface{}) (int64, error) {
	if m.options == nil {
		return 0, errors.New("model.init function needs to be executed in new model")
	}
	if len(m.options.defaultOrders) == 0 {
		return 0, errors.New("default orders can not be nil")
	}
	if m.options.defaultPage == 0 {
		return 0, errors.New("default page can not be 0")
	}
	if m.options.defaultSize == 0 {
		return 0, errors.New("default size can not be 0")
	}

	session := m.GetSession()

	buildFilters(session, query.Filters)

	buildQ(session, query.Q)

	if len(query.Orders) > 0 {
		buildOrders(session, query.Orders)
	} else {
		buildOrders(session, m.options.defaultOrders)
	}

	if query.ExtraQuery != nil {
		query.ExtraQuery(session)
	}

	if query.Page <= 0 {
		query.Page = m.options.defaultPage
	}

	if query.Size <= 0 {
		query.Size = m.options.defaultSize
	}

	session.Limit(query.Size, (query.Page-1)*query.Size)

	total, err := session.FindAndCount(items)

	if err != nil {
		return 0, err
	}

	return total, nil
}

// todo group list
