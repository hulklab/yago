module github.com/hulklab/yago

go 1.12

require (
	github.com/astaxie/beego v1.12.0 // indirect
	github.com/bsm/sarama-cluster v2.1.15+incompatible // indirect
	github.com/bwmarrin/snowflake v0.3.0
	github.com/garyburd/redigo v1.6.0 // indirect
	github.com/gin-contrib/cors v1.3.0
	github.com/gin-contrib/gzip v0.0.1
	github.com/gin-contrib/pprof v1.2.1
	github.com/gin-gonic/gin v1.4.0
	github.com/go-sql-driver/mysql v1.4.1
	github.com/go-xorm/xorm v0.7.6
	github.com/golang/protobuf v1.3.1
	github.com/natefinch/lumberjack v2.0.0+incompatible
	github.com/robfig/cron v1.2.0
	github.com/sirupsen/logrus v1.4.0
	github.com/spf13/cobra v0.0.5
	github.com/spf13/viper v1.4.0
	golang.org/x/net v0.0.0-20190603091049-60506f45cf65
	google.golang.org/grpc v1.23.0
	xorm.io/core v0.7.0
)

replace github.com/ugorji/go v1.1.4 => github.com/ugorji/go/codec v0.0.0-20190204201341-e444a5086c43
