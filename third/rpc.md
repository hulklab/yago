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
    
    req := &pb.HelloRequest{Name: name}

    conn, _ := a.GetConn()
    ctx, cancel := a.GetCtx()
    defer cancel()

    c := pb.NewHomeClient(conn)
    resp, err := c.Hello(ctx, req)

    if err != nil {
        fmt.Println(err)
        return
    }
    fmt.Println("ok:", resp.Data)
}

// 采用 stream 流的方式
func (a *homeApi) HelloStream() {

	req := &pb.HelloRequest{Name: name}

	conn, _ := a.GetConn()
	ctx, cancel := a.GetCtx()
	defer cancel()

	c := pb.NewHomeClient(conn)
	stream, err := c.HelloStream(ctx, req)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	for {
		reply, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println(err)
			break
		}

		fmt.Println(reply)
	}
}
```


### 调用 api
```go
homeapi.Ins().Hello()
homeapi.Ins().HelloStream()

```

### grpc-Interceptor
RpcThird 使用了 grpc-Interceptor, 默认提供一个日志 interceptor，并对外提供了添加自定义 interceptor 的函数。
interceptor 支持添加多个，执行顺序为添加的顺序，默认的日志 interceptor 总是在最后执行。

* UnaryClientInterceptor 

```go
func (a *homeApi) RpcHello() {
	a.AddUnaryClientInterceptor(func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
        fmt.Println("method:", method, "client before")
        err := invoker(ctx, method, req, reply, cc, opts...)
        fmt.Println("method:", method, "client after")
        return err
    })
    // ...
}
```

* StreamClientInterceptor

```go
func (a *homeApi) RpcHelloStream() {
    a.AddStreamClientInterceptor(func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (stream grpc.ClientStream, e error) {
        fmt.Println("method:", method, "stream before")
        clientStream, err := streamer(ctx, desc, cc, method, opts...)
        return clientStream, err
    })
    // ...
}
```

* 关闭默认的日志 interceptor

```go
// 关闭 unary
a.DisableDefaultUnaryInterceptor()

// 关闭 stream
a.DisableDefaultStreamInterceptor()
```
