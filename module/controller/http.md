### Http控制器

http控制器我们也使用了一个web框架[gin](https://github.com/gin-gonic/gin)，至于为什么选用gin，是因为它够轻量，同时社区是用的人比较多。

我们看一下控制器的定义

```go
type HomeHttp struct {
	basehttp.BaseHttp
}
```

可以看到控制器的名字是HomeHttp，它继承自basehttp.BaseHttp，这里顺便提下yago/base目录，里面存放了一些结构体的父类定义，主要是为了以后有统一的功能扩展可以方便平滑的实现。这里basehttp里面实现两个钩子函数BeforeAction，AfterAction。如果需要在Action前后做一些处理的话，可以在自己的http控制器中覆写这两个方法。

举个例子，我们如果需要做auth认证我们就可以定一个auth控制器，然后在auth控制器中的BeforeAction函数里面实现这个逻辑，然后所有需要做auth的控制器都来继承这个auth的控制器就可以了。AfterAction同理。

#### 路由注册
init函数完成Action的路由注册，正如我们在[路由注册](../route/route.md)中提到的，这是路由注册的一个阶段。注册函数参数见下表

```go
func init() {
	homeHttp := new(HomeHttp)
	yago.AddHttpRouter("/home/hello", http.MethodGet, homeHttp.HelloAction, homeHttp)
	yago.AddHttpRouter("/home/add", http.MethodPost, homeHttp.AddAction, homeHttp)
}
```

AddHttpRouter参数说明

| 参数位置 | 参数类型 | 说明 |
| ------- | ------- | ------- |
| 1 | String | http请求路径 |
| 2 | String | 允许访问的http method |
| 3 | Func | http接口对应的Action Func |
| 4 | Struct | http控制器对象 |

#### HttpAction

Http Action接收一个参数c *yago.Ctx，它是gin.Ctx的扩展，主要用来获取参数和返回响应。

```go
func (h *HomeHttp) HelloAction(c *yago.Ctx) {
	name := c.RequestString("name")

	c.SetData("hello " + name)

	return
}

func (h *HomeHttp) AddAction(c *yago.Ctx) {
	name := c.RequestString("name")

	model := homemodel.NewHomeModel()
	id, err := model.Add(name, nil)
	if err.HasErr() {
		c.SetError(err)
		return
	}

	c.SetData(map[string]interface{}{"id": id})
	return
}
```

我们扩展了Request系列的方法，用来整合query_args和body_args，并且提供了类型转换。例如可以直接通过c.RequestSliceString("names")来获取一个逗号分隔的字符串类型的参数值并将其转换成切片返回。

Action内，可以通过c.SetData函数来返回正确的结果响应（json格式），或者c.SetError + return来返回错误信息（json）。c.SetError接收一个yago.Err类型的error，yago.Err定义来自 app/g/errors.go。 需要说明的是，c.SetError并不能阻止程序往下运行，如果需要接口中断，请加return。

#### Labels & Rules

Labels，Rules为验证器的两个函数，具体使用请看[validator](../library/validator.md)。

```go
func (h *HomeHttp) Labels() validator.Label {
	return map[string]string{
		"name":     "姓名",
	}
}

func (h *HomeHttp) Rules() []validator.Rule {
	return []validator.Rule{
		{
			Params: []string{"name"},
			Method: validator.Required,
			On:     []string{"add"},
		},
	}
}
```

#### WebSocket

如何在http控制器中使用websocket，这里给一个简短的示例。我们使用[Gorilla WebSocket](https://github.com/gorilla/websocket)一个websocket框架来完成。

在init函数中完成协议升级。

```go
type HelloHttp struct {
	basehttp.BaseHttp
	upGrader  *websocket.Upgrader
}

func init() {
	helloHttp := new(HelloHttp)
	helloHttp.upGrader = &websocket.Upgrader{
		HandshakeTimeout: 0,
		ReadBufferSize:   0,
		WriteBufferSize:  0,
		WriteBufferPool:  nil,
		Subprotocols:     nil,
		Error:            nil,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
		EnableCompression: false,
	}

	yago.AddHttpRouter("/hello", http.MethodGet, helloHttp.HelloAction, helloHttp)
}

func (h *HelloHttp) HelloAction(c *yago.Ctx) {
	ws, err := h.upGrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.SetError(yago.ErrSystem, err.Error())
		return
	}

	defer ws.Close()

	for {
		mt, message, err := ws.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		log.Printf("recv: %s", message)
		err = ws.WriteMessage(mt, message)
		if err != nil {
			log.Println("write:", err)
			break
		}
	}
}
```
