package orm

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"github.com/hulklab/yago"
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

		dsn := dbUser + ":" + dbPassword + "@tcp(" + dbHost + ":" + dbPort + ")/" + dbName + "?charset=utf8"

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

		return orm
	})

	return v.(*Orm)
}
