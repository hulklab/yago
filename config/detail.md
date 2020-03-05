# Yago 程序的配置文件

## 应用配置

```toml
[app]

app_name = "app"
env = "dev"
debug = true
# 如果不设置则不会创建 pidfile
# pidfile = "/var/run/app.pid"

# 是否开启http服务
http_enable = true
# http服务地址
http_addr = ":8080"
# http服务关闭最大等待时长, 秒
http_stop_time_wait = 10

# http ssl config
# http_ssl_on = true
# http_cert_file = "./yourdomain.crt"
# http_key_file = "./yourdomain.key"

# cors 跨域
# http_cors_allow_all_origins = true
# http_cors_allow_origins = []
# http_cors_allow_methods = ["GET", "POST", "PUT", "PATCH", "DELETE", "HEAD"]
# http_cors_allow_headers = ["Origin", "Content-Length", "Content-Type"]
# http_cors_expose_headers = []
# http_cors_allow_credentials = true
# http_cors_max_age = "12h"

# gzip 模式 1:Default, 2:Best Speed, 3:Best Compression
# http_gzip_on = true
# http_gzip_level = 1

# pprof route: /debug/pprof
# http_pprof_on = false

# http html 模版配置
# http_view_render = true
# http_view_path = "views/*"
# http_static_path = "static/js"

# 是否开启rpc服务
rpc_enable = true
# rpc服务地址
rpc_addr = ":50051"
# rpc服务关闭最大等待时长, 秒
rpc_stop_time_wait = 10
# rpc reflection
# rpc_reflect_on = true

# rpc ssl config
# rpc_ssl_on = true
# rpc_cert_file = "./conf/server.pem"
# rpc_key_file = "./conf/server.key"

# 是否开启task任务
task_enable = true
# http服务关闭最大等待时长, 秒
task_stop_time_wait = 10

# 组件资源关闭最大等待时长, 秒
com_stop_time_wait = 10
```

| 配置项 | 类型 | 说明 |
| ------- | ------- |------- |
| app.app_name | String | 程序名称， 默认 “APP”， 作为配置文件环境变量的前缀 <br/> 比如 export APP_APP_TASK_ENABLE=false，<br/>就相当于 app.task_enable=false，<br/>而且环境变量配置项会比配置文件的中的值有更高的优先级 |
| app.env | String | 运行环境: dev, prod |
| app.debug | Bool | 是否开启debug<br/>debug=true会输出更详细的控制台调试信息|
| app.pidfile | String| 设置保存进程号的文件，默认不设置不创建 |
| app.http_enable | Bool | 是否开启http服务 |
| app.http_addr | String | http服务监听的ip地址和端口 |
| app.http_stop_time_wait | Duration | http服务收到关闭信号时的最大等待时长，默认10s |
| app.http_ssl_on | Bool | 是否开启https |
| app.http_cert_file | String | https证书文件地址 |
| app.http_cert_key | String | https证书私钥 |
| app.http_cors_allow_all_origins | Bool | 是否允许所有origin访问 |
| app.http_cors_allow_origins | StringSlice | 允许哪些指定的origin访问<br>使用时需要将app.http_cors_allow_all_origins 设为 false |
| app.http_cors_allow_methods | StringSlice | 允许客户端使用哪些方法发起请求 |
| app.http_cors_allow_headers | StringSlice | 服务器支持的所有头信息字段 |
| app.http_cors_allow_credentials | Bool | 是否允许发送Cookie |
| app.http_cors_max_age| Duration | 用来指定预检请求的有效期，单位为秒 |
| app.http_gzip_on | Bool | 是否开启http压缩，默认为true |
| app.http_gzip_level | Int | http压缩模式<br/>1：默认，2：速度优先 3：大小优先 |
| app.http_pprof_on | Bool | 是否开启pprof，默认为false |
| app.http_view_render | Bool | 是否开启模版渲染 |
| app.http_view_path | String | 模版路径 |
| app.http_static_path | String | 静态资源路径 |
| app.rpc_enable | Bool | 是否开启grpc服务 |
| app.rpc_addr | String | grpc服务监听的ip地址和端口 |
| app.rpc_stop_time_wait | Duration | grpc服务收到关闭信号时的最大等待时长，默认10s | 
| app.rpc_reflect_on | Bool | 开启grpc_cli调用 |
| app.rpc_ssl_on | Bool | 是否grpc SSL服务 |
| app.rpc_cert_file | String | grpc SSL证书文件地址 |
| app.rpc_key_file | String | grpc SSL私钥文件地址 |
| app.task_enable | Bool | 是否开启task服务 ｜
| app.task_stop_time_wait | Duration | task服务收到关闭信号时的最大等待时长，默认10s |
| app.com_stop_time_wait | Duration | 组件资源关闭最大等待时长，默认10s |


## Logger 组件

```toml
[logger]
# json | text, default text
formatter = "json"
# 日志最低等级 Panic = 0, Fatal = 1, Error = 2, Warn = 3, Info = 4, Debug = 5, Trace = 6
level = 5
# 文件路径
file_path = "./logs/app.log"
# 最大保留的备份数
max_backups = 20
# 日志最大保留天数
max_age = 30
# 文件最大大小(mb)
max_size = 500
# 是否开启压缩
compress = true
# write log to stdout
# stdout_enable = true
```

| 配置项 | 类型 | 说明 |
| ------- | ------- |------- |
| logger.formatter | String | logger组件的日志格式，默认text，还有json可选 |
| logger.level | Int | 日志最低等级 <br> Panic = 0, Fatal = 1, Error = 2, Warn = 3, Info = 4, Debug = 5, Trace = 6 <br> 数字小于等于该值的日志才会显示 |
| logger.file_path | String | 日志文件的输出地址 |
| logger.max_backups | Int | 最大保留的备份数 |
| logger.max_age | Duration | 日志最大保留天数 |
| logger.max_size | Int | 文件最大大小(mb) |
| logger.compress | Bool | 是否开启日志文件Gzip压缩 |
| logger.stdout_enable | Bool | 是否打印到标准输出 |

## 数据库组件(SQL)

```toml
[db]
host = "127.0.0.1"
user = "user"
password = "password"
port = "3306"
database = "db"
prefix =""
timezone = "Asia/Shanghai"
charset = "utf8"
max_life_time = 8
max_idle_conn = 20
max_open_conn = 500
show_log = true
```

| 配置项 | 类型 | 说明 |
| ------- | ------- |------- |
| db.host | String | 数据库ip地址，默认127.0.0.1 |
| db.port | String | 数据库端口地址，默认3306 |
| db.user | String | 数据库用户 |
| db.password | String | 数据库用户密码 |
| db.database | String | 默认数据库 |
| db.prefix | String | 数据库表名前缀 |
| db.timezone | String | 数据库时区 Asia/Shanghai |
| db.charset | String | 默认字符集 |
| db.max_life_time | Duration | 最大连接生命周期 |
| db.max_idle_conn | Int | 连接池最大空闲连接数，默认20 |
| db.max_open_conn | Int | 连接池最大打开连接数，默认500 |
| db.show_log | Bool | 是否打印 SQL 日志，默认true |

## Redis 组件

```toml
[redis]
addr = "127.0.0.1:6379"
auth = ""
db = 0
max_idle = 5
idle_timeout = 30
```

| 配置项 | 类型 | 说明 |
| ------- | ------- |------- |
| redis.addr | String | Redis Server ip加端口号，默认 "127.0.0.1:6379" |
| redis.auth | String | Redis Server 认证信息 |
| redis.db | Int | Redis Server 默认db，默认 0 |
| redis.max_idle | Int | Redis 连接池最大空闲连接数 |
| redis.idle_timeout | Int | Redis 连接池空闲超时时间 |

## Mongodb 组件

```toml
[mongodb]
# https://docs.mongodb.com/manual/reference/connection-string/
mongodb_uri = "mongodb://user:password@127.0.0.1:27017/?connectTimeoutMS=5000&socketTimeoutMS=5000&maxPoolSize=100"
database = "test"
```

| 配置项 | 类型 | 说明 |
| ------- | ------- |------- |
| mongodb.mongodb_uri | String | Mongodb 连接 uri<br>参考https://docs.mongodb.com/manual/reference/connection-string/ |
| mongodb.database | String | 默认数据库 |

## Kafka 组件

```toml
[kafka]
cluster = "127.0.0.1:9092"
```

| 配置项 | 类型 | 说明 |
| ------- | ------- |------- |
| kafka.cluster | String | Kafka 集群地址 |

## 第三方 API 调用

```toml
[home_api]
domain = "http://127.0.0.1:8080"
hostname = "localhost"
rpc_address = "127.0.0.1:50051"
timeout = 10
max_recv_msgsize_mb =  10
max_send_msgsize_mb =  10
# ssl_on = true
# cert_file = "./conf/server.pem"
```

| 配置项 | 类型 | 说明 |
| ------- | ------- |------- |
| *_api.domain | String | http请求ip+端口或者是已经域名解析的域名+端口 |
| *_api.host | String | http请求指定的host |
| *_api.rpc_address | String | grpc服务地址 |
| *_api.timeout | Duration | 请求超时时间 |
| *_api.max_recv_msgsize_mb | Int | grpc可以接收的最大消息大小，默认 4m |
| *_api.max_send_msgsize_mb | Int | grpc可以发送的最大消息大小，默认 4m |
| *_api.ssl_on | Bool | grpc请求是否开启证书校验 |
| *_api.cert_file | String | grpc证书地址 |

## 配置嵌套
Yago 采用 import 字段自己实现了一套配置嵌套的方法，下面演示几种常用格式的使用方法：

### toml
```toml
import = "./conf/app.toml"

[app]
http_addr = ":8088"

```

### json
```json
{
  "import": "./conf/app.toml",
  "app": {"http_addr": ":8088"}
}
```

### yaml
```yaml
import: "./conf/app.toml"
app:
  http_addr: ":8080"

```

其中 `import` 设置的值表示继承配置的路径，可以是相对路径或者绝对路径，Yago 支持多级继承，但是使用时要避开循环引用。
