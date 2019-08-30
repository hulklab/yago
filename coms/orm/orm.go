package orm

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"github.com/hulklab/yago"
	"github.com/hulklab/yago/libs/logger"
	"github.com/sirupsen/logrus"
	"log"
	"net/url"
	"time"
	"xorm.io/core"
)

type Orm struct {
	*xorm.Engine
}

// 扩展了一个事务功能
func (orm *Orm) Transactional(f func(session *xorm.Session) error) (err error) {
	session := orm.NewSession()
	defer session.Close()
	err = session.Begin()
	if err != nil {
		return err
	}

	defer func() {
		if p := recover(); p != nil {
			err1 := session.Rollback()
			log.Println("err occur in db transaction:", err1.Error())
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

		val, _ := xorm.NewEngine(driver.(string), dsn)

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
			orm.SetLogger(getLogger(showLog.(bool)))
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

type Logger struct {
	*logrus.Entry
	show bool
}

func (l *Logger) Level() core.LogLevel {
	return 0
}

func (l *Logger) SetLevel(c core.LogLevel) {
}

func (l *Logger) ShowSQL(show ...bool) {
	if len(show) > 0 {
		l.show = show[0]
	}
}

func (l *Logger) IsShowSQL() bool {
	return l.show
}
