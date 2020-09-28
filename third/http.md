# HttpThird

`yago/base/basethird/http.go` 是对开源包 `github.com/levigross/grequests` 的再次封装，
主要的作用是统一规范第三方接口的调用方式，简化业务层的调用，统一记录调用的日志。

## 如何使用
### 定义第三方调用的 ThirdApi
每个第三方单独定义一个目录（package），统一放置于 `app/third` 目录下。
本文以 yago 中 example 里面的 homeapi 为例，homeapi 即为一个第三方调用，
它具体的 homeApi 结构体定义于 homeapi/home.go 中。

```go
type homeApi struct {
	basethird.HttpThird
}
```

由于 homeApi 组合了 basethird.HttpThird，即拥有了所有 HttpThird 的调用方法，下文中给出调用样例。

### 配置第三方 http 接口
```toml
[home_api]
domain = "http://127.0.0.1:8080"
# Host 配置，如果域名已解析， hostname 可以设置为空串
hostname = "localhost"
# 读写超时时间，单位 s
timeout = 10
```
我们在模版 app.toml 中给出了配置 homeapi 的样例。

### 实现 api 的实例化 Ins 方法
```go
func Ins() *homeApi {
	name := "home_api"
	v := yago.Component.Ins(name, func() interface{} {
		api := new(homeApi)
        
        err := api.InitConfig(name)
        if err != nil {
            log.Fatalf("init %s api failed",name)
        }

		return api
	})
	
	return v.(*homeApi)
}
```
可以使用组件的方式来实例化 ThirdApi，`grequests` 的连接采用的是连接池。

### 实现接口调用的方法
通常第三方的每个接口参数和返回值都不一样，我们需要在 ThirdApi 中为每个接口定义一个方法，下面给出示例。

* 普通的 post, get 请求

```go
func (a *homeApi) GetUserById(id int64) (*basethird.Response,error) {
	params := map[string]interface{}{
		"id": id,
	}

	resp, err := a.Get("/home/user/detail", params)
	return resp, err
}
```

* 文件上传的请求

```go
func (a *homeApi) UploadFile(filepath string) (*basethird.Response, error){

	params := map[string]interface{}{
		"file": basethird.PostFile(filepath),
	}

	resp, err := a.Post("/home/user/upload", params)
	return resp, err
}
```

* 直接传 body(Content-type: application/json 的请求)

```go
func (a *homeApi) AddUser(id int64,name string) (*basethird.Response, error){

    u := g.Hash{"id":id, "name":name}

    resp, err := a.Post("/home/user/add", nil, &grequests.RequestOptions{
        JSON: u,
    })
    return resp, err
}
```

### 调用 api
```go
// 普通的请求
resp, err := homeapi.Ins().GetUserById(1)

// 上传文件的请求
resp, err := homeapi.Ins().UploadFile("/tmp/test.jpeg")

// 直接请求 Body
resp, err := homeapi.Ins().AddUser(1,"zhangsan")
```

### http-Interceptor
HttpThird 实现了 Interceptor 中间件机制, 默认提供一个日志 interceptor，并对外提供了添加自定义 interceptor 的函数。
interceptor 支持添加多个，执行顺序为添加的顺序。

* HttpInterceptor 

```go
// 添加中间件
api.AddInterceptor(func(method, uri string, ro *grequests.RequestOptions, call basethird.Caller) (response *basethird.Response, e error) {
    fmt.Println("before caller....", uri, method)

    resp, err := call(method, uri, ro)

    fmt.Println("after caller....", resp.StatusCode)

    return resp, err
})
```

* 关闭默认的日志 interceptor

```go
// 关闭
api.DisableDefaultInterceptor()
```
