package orm

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/url"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/hulklab/yago"
	"github.com/hulklab/yago/coms/logger"
	"github.com/sirupsen/logrus"
	"xorm.io/xorm"
	xormLog "xorm.io/xorm/log"
)

type Orm struct {
	*xorm.Engine
}

// 扩展了一个事务功能
func (o *Orm) Transactional(f func(session *xorm.Session) error, opts ...OrmOption) (err error) {
	session := o.NewSession()
	ormArg := ExtractOption(opts...)
	if ormArg.ctx != nil {
		session = session.Context(ormArg.ctx)
	}

	defer session.Close()
	err = session.Begin()
	if err != nil {
		return err
	}

	defer func() {
		if p := recover(); p != nil {
			err1 := session.Rollback()
			if err1 != nil {
				log.Println("err occur in db transaction:", err1.Error())
			}
			panic(p)
		} else if err != nil {
			if err2 := session.Rollback(); err2 != nil {
				err = fmt.Errorf("execute %s, rollback err: %s", err.Error(), err2.Error())
			}
		} else {
			err = session.Commit()
		}
	}()

	err = f(session)
	return err
}

type OrmArg struct {
	Session *xorm.Session
	ctx     context.Context
}

type OrmOption func(arg *OrmArg)

func WithSession(session *xorm.Session) OrmOption {
	return func(arg *OrmArg) {
		arg.Session = session
	}
}

func WithContext(ctx context.Context) OrmOption {
	return func(arg *OrmArg) {
		arg.ctx = ctx
	}
}

func ExtractOption(opts ...OrmOption) OrmArg {
	arg := OrmArg{}
	for _, opt := range opts {
		opt(&arg)
	}
	return arg
}

// 添加或者修改的原子操作，但是要求 columns 里面必须包含唯一键，否则会一直执行添加操作
// Upsert("table_name",g.Hash{"name":"zhangsan","uuid":"abcdef"})
// Upsert("table_name",g.Hash{"name":"zhangsan","uuid":"abcdef"},orm.WithSession(session)) 事务中使用
func (o *Orm) Upsert(table interface{}, columns map[string]interface{}, opts ...OrmOption) (sql.Result, error) {
	if table == nil {
		return nil, errors.New("table is required in orm upsert")
	}

	if len(columns) == 0 {
		return nil, errors.New("columns is required in orm upsert")
	}

	cols := make([]string, 0)
	args := make([]interface{}, 0)
	values := make([]interface{}, 0)
	placeholders := make([]interface{}, 0)

	for field, value := range columns {
		cols = append(cols, Ins().Quote(field))
		values = append(values, value)
		placeholders = append(placeholders, value)
	}

	tableName := o.TableName(table)
	colStr := strings.Join(cols, ", ")
	valuePlaceStr := strings.TrimLeft(strings.Repeat(", ?", len(cols)), ", ")
	updateStr := strings.Join(cols, " = ?, ") + "= ?"

	statement := "INSERT INTO " + tableName + "(" + colStr + ") values(" + valuePlaceStr + ") on DUPLICATE key update " + updateStr + ";"

	args = append(args, statement)

	for _, place := range placeholders {
		args = append(args, place)
	}

	for _, val := range values {
		args = append(args, val)
	}

	var session *xorm.Session
	ormArg := ExtractOption(opts...)

	if ormArg.Session != nil {
		session = ormArg.Session

		if ormArg.ctx != nil {
			session = session.Context(ormArg.ctx)
		}
	}

	if session == nil {
		if ormArg.ctx != nil {
			return o.Context(ormArg.ctx).Exec(args...)
		}

		return o.Exec(args...)
	} else {
		return session.Exec(args...)
	}
}

// 返回 orm 组件单例
func Ins(id ...string) *Orm {
	var name string

	if len(id) == 0 {
		name = "db"
	} else if len(id) > 0 {
		name = id[0]
	}

	v := yago.Component.Ins(name, func() interface{} {

		// 注：从 Config 里面取出的整型是 int64
		conf := yago.Config.GetStringMap(name)

		dbHost := conf["host"].(string)
		dbPort := conf["port"].(string)
		dbUser := conf["user"].(string)
		dbPassword := conf["password"].(string)
		dbName := conf["database"].(string)
		charset := conf["charset"].(string)

		if dbHost == "" {
			log.Fatalf("Fatal error: Sql host is empty")
		}
		if dbPort == "" {
			log.Fatalf("Fatal error: Sql port is empty")
		}

		dsn := dbUser + ":" + dbPassword + "@tcp(" + dbHost + ":" + dbPort + ")/" + dbName + "?charset=" + charset

		timezone, ok := conf["timezone"]
		if ok {
			dsn = dsn + "&loc=" + url.QueryEscape(timezone.(string))
		}

		driver, ok := conf["driver"]
		if !ok {
			driver = "mysql"
		}

		val, err := xorm.NewEngine(driver.(string), dsn)
		if err != nil {
			log.Fatalf("[ORM] Fatal error: new orm engine err, %s", err.Error())
		}

		orm := &Orm{
			val,
		}

		// 连接生存时间
		maxLife, ok := conf["max_life_time"]
		if ok {
			orm.DB().SetConnMaxLifetime(time.Duration(maxLife.(int64)) * time.Second)
		}

		// 最大空闲连接
		maxIdle, ok := conf["max_idle_conn"]
		if ok {
			orm.DB().SetMaxIdleConns(int(maxIdle.(int64)))
		}

		// 最大打开连接数
		maxOpen, ok := conf["max_open_conn"]
		if ok {
			orm.DB().SetMaxOpenConns(int(maxOpen.(int64)))
		}

		// 设置日志
		showLog, ok := conf["show_log"]
		if ok {
			if ctxLogger != nil {
				orm.SetLogger(ctxLogger)
				orm.ShowSQL(showLog.(bool))
			} else {
				orm.SetLogger(getLogger(showLog.(bool)))
			}
		}

		return orm
	})

	return v.(*Orm)
}

func getLogger(show bool) *Logger {
	entry := logger.Ins().WithFields(logrus.Fields{"category": "orm.sql"})

	lg := &Logger{
		Entry: entry,
		show:  show,
	}

	return lg
}

var ctxLogger xormLog.ContextLogger

func RegisterCtxLogger(ctxLog xormLog.ContextLogger) {
	ctxLogger = ctxLog
}

type Logger struct {
	*logrus.Entry
	show bool
}

func (l *Logger) Level() xormLog.LogLevel {
	return 0
}

func (l *Logger) SetLevel(c xormLog.LogLevel) {
}

func (l *Logger) ShowSQL(show ...bool) {
	if len(show) > 0 {
		l.show = show[0]
	}
}

func (l *Logger) IsShowSQL() bool {
	return l.show
}
