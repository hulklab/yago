# Http 控制器

http 控制器我们依赖于 web 框架 [gin](https://github.com/gin-gonic/gin)，至于为什么选用 gin，是因为它够轻量，同时社区是用的人比较多。

我们看一下控制器的定义

```go
type HomeHttp struct {
	basehttp.BaseHttp
}
```

可以看到控制器的名字是 HomeHttp，它继承自 basehttp.BaseHttp，这里顺便提下 yago/base 目录，里面存放了一些结构体的父类定义，主要是为了以后有统一的功能扩展可以方便平滑的实现。这里 basehttp 里面实现两个钩子函数 BeforeAction，AfterAction。如果需要在 Action 前后做一些处理的话，可以在自己的 http 控制器中覆写这两个方法。

举个例子，我们如果需要做 auth 认证我们就可以定一个 auth 控制器，然后在 auth 控制器中的 BeforeAction 函数里面实现这个逻辑，然后所有需要做 auth 的控制器都来继承这个 auth 的控制器就可以了。AfterAction 同理。

## 路由注册

init 函数完成 Action 的路由注册，正如我们在[路由注册](/route/route.md)中提到的，这是路由注册的一个阶段。注册函数参数见下表

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
| 1 | String | http 请求路径 |
| 2 | String | 允许访问的 http method |
| 3 | Func | http 接口对应的 Action Func |
| 4 | Struct | http 控制器对象 |

## HttpAction

Http Action 接收一个参数 c *yago.Ctx，它是 gin.Ctx 的扩展，主要用来获取参数和返回响应。

```go
func (h *HomeHttp) HelloAction(c *yago.Ctx) {
    var p struct {
        Name string `json:"name" validate:"omitempty,max=20" form:"name"`
    }
    err := c.ShouldBind(&p)
    if err != nil {
        c.SetError(err)
        return
    }

	c.SetData("hello " + name)
	return
}

func (h *HomeHttp) AddAction(c *yago.Ctx) {
    var p struct {
        Name string `json:"name" validate:"required,max=20" form:"name" label:"姓名"`
    }
    
    err := c.ShouldBind(&p)
    if err != nil {
        c.SetError(err)
        return
    }

	model := homemodel.NewHomeModel()
	id, e := model.Add(p.Name, nil)
	if e.HasErr() {
		c.SetError(e)
		return
	}

	c.SetData(map[string]interface{}{"id": id})
	return
}
```

Action 内，可以通过 c.SetData 函数来返回正确的结果响应（json 格式），或者 c.SetError + return 来返回错误信息（json）。c.SetError 接收一个 yago.Err 类型的 error，yago.Err 定义来自 app/g/errors.go。 需要说明的是，c.SetError 并不能阻止程序往下运行，如果需要接口中断，请加 return。

## Validate

validate 采用 gin 的 validator，具体使用请看 [validator](../library/validator.md)。


## WebSocket 服务

如何在 http 控制器中使用 websocket，这里给一个简短的服务端示例。我们使用 [Gorilla WebSocket](https://github.com/gorilla/websocket)  websocket 框架来完成。

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
