# Redis 组件
Redis 组件我们依赖的开源包是 `github.com/garyburd/redigo/redis`。

按照组件的设计，我们定义了自己的 Redis 结构对其进行了组合，在保留其原生的功能之外，以便扩展。

```go
// yago/coms/rds/redis.go
type Rds struct {
	*redis.Pool
}
```

所以你可以查看 [redigo 官方文档](https://github.com/gomodule/redigo) 来获取所有支持的 api。

本文中仅介绍部分常用的 api 以及扩展的 api。


## 配置 Redis 组件
```toml
[redis]
addr = "127.0.0.1:6379"
auth = ""
db = 0
max_idle = 5
idle_timeout = 30
```
我们在模版 app.toml 中默认配置开启了日志组件，可根据实际情况进行调整。

## 使用 Redis 组件
* 使用 Do 命令

yago 对 Do 命令进行了封装，执行完 Do 命令之后会回收连接，若想使用原生的 Do，请调用 `rds.Ins().GetConn().Do()` 

```go
rc := rds.Ins()
reply, err := rc.Do("SET","test_key","senyuan","NX")

v, err := redis.String(rc.Do("GET","test_key"))
```
> redigo 对 Do 的返回值做了一些封装，除了 redis.String 外，其他的参考 [参考文件](https://github.com/gomodule/redigo/blob/master/redis/reply.go)


* 使用 yago 封装的 cmd

```go
rc := rds.Ins()
s, err := redis.String(rc.Set("test_key2","shanglin","ex",10))

s2, err := redis.String(rc.Get("test_key2"))

```
>  由于 redis 命令太多，此处不一一举例了，常用的命令都封装在 `yago/coms/rds/redis_cmd.go` 中

* 发布与订阅

```go
subscriber, err := rds.Ins().NewSubscriber(topic)
if err != nil {
    return
}
defer subscriber.Close()

go func() {
    r := rds.Ins()
    defer r.Close()
    time.Sleep(time.Second)
    // 发布
    r.Do("publish", topic, "hello")
}()

// 订阅
err = subscriber.Subscribe(func(bytes []byte) {
    fmt.Println("msg:", string(bytes))
})
```

