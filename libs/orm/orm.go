package orm

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/core"
	"github.com/go-xorm/xorm"
	"github.com/hulklab/yago"
	"github.com/hulklab/yago/libs/logger"
	"github.com/sirupsen/logrus"
	"net/url"
	"time"
)

type Orm struct {
	*xorm.Engine
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
			orm.ShowSQL(showLog.(bool))
			orm.SetLogger(getLogger())
		}

		return orm
	})

	return v.(*Orm)
}

func getLogger() *Logger {

	entry := logger.Ins().WithFields(logrus.Fields{"category": "orm.sql"})

	lg := &Logger{
		Entry: entry,
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
