# RpcThird

`yago/base/basethird/rpc.go` 是对 `grpc.ClientConn` 的再次封装，
主要的作用是统一规范第三方 rpc 接口的调用方式，简化业务层的调用，统一记录调用的日志。

## 如何使用
### 定义第三方调用的 ThirdApi
每个第三方单独定义一个目录（package），统一放置于 `app/third` 目录下。
本文以 yago 中 example 里面的 homeapi 为例，homeapi 即为一个第三方调用，
它具体的 homeApi 结构体定义于 homeapi/home.go 中。

```go
type homeApi struct {
	basethird.RpcThird
}
```

由于 homeApi 组合了 basethird.RpcThird，即拥有了所有 RpcThird 的调用方法，下文中给出调用样例。

### 配置第三方 rpc 接口
```toml
[home_api]
address = "127.0.0.1:50051"
# Host 配置，如果域名已解析， hostname 可以设置为空串
hostname = "localhost"
# 读写超时时间，单位 s
timeout = 10
# 如果 rpc 服务端开启 ssl，客户端需要打开 ssl_on，并指定公钥
# ssl_on = true
# cert_file = "./conf/server.pem"

```
我们在模版 app.toml 中给出了配置 homeapi 的样例。

### 定义 protobuf 文件
客户端的 protobuf 一般由服务端提供，这里我们把 home.proto 文件放在 homeapi/protobuf/homepb 目录下。
go 中间文件 `home.pb.go` 的生成参考 [Protobuf 规范](/module/controller/protobuf.md)

### 实现实例化 api 的 Ins 方法
```go
func Ins() *homeApi {
	name := "home_api"
	v := yago.Component.Ins(name, func() interface{} {
		api := new(homeApi)
		api.Address = yago.Config.GetString(name + ".address")
		api.Hostname= yago.Config.GetString(name + ".hostname")
		api.Timeout = yago.Config.GetInt(name + ".timeout")
        // api.SslOn = yago.Config.GetBool(name + ".ssl_on")
        // api.CertFile = yago.Config.GetString(name + ".cert_file")
		return api
	})
	return v.(*homeApi)
}
```
可以使用组件的方式来实例化 ThirdApi，`grpc` 的连接采用的是连接池。

### 实现接口调用的方法
通常第三方的每个接口参数和返回值都不一样，我们需要在 ThirdApi 中为每个接口定义一个方法，下面给出示例。

```go
func (a *homeApi) Hello(name string) () {
    rep, err := a.Call(func(conn *grpc.ClientConn, ctx context.Context) (rep proto.Message, e error) {

        c := pb.NewHomeClient(conn)

        return c.Hello(ctx, &pb.HelloRequest{Name: name})

    }, name)
    
	if err != nil {
		fmt.Println(err)
		return
	}

	v, ok := rep.(*pb.HelloReply)
	if ok {
		fmt.Println("ok:", v.Data)
	} else {
		fmt.Println("not match", v)
	}
}
```


### 调用 api
```go
homeapi.Ins().Hello()

```


