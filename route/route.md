# 路由 Router

## Yago 路由注册分为两个阶段：

* step1. 在各个控制器的 init 函数中完成一个 Controller 内具体的 Action 级别的路由注册，详细使用请参考 [模块控制器](/module/controller.md)

```go
func init() {
	homeHttp := new(HomeHttp)
	yago.AddHttpRouter("/home/hello", http.MethodGet, homeHttp.HelloAction, homeHttp)
	yago.AddHttpRouter("/home/add", http.MethodPost, homeHttp.AddAction, homeHttp)
	yago.AddHttpRouter("/home/delete", http.MethodPost, homeHttp.DeleteAction, homeHttp)
	yago.AddHttpRouter("/home/detail", http.MethodGet, homeHttp.DetailAction, homeHttp)
	yago.AddHttpRouter("/home/update", http.MethodPost, homeHttp.UpdateAction, homeHttp)
	yago.AddHttpRouter("/home/list", http.MethodPost, homeHttp.ListAction, homeHttp)
	yago.AddHttpRouter("/home/upload", http.MethodPost, homeHttp.UploadAction, homeHttp)
}
```

* step2. 在 app/route/route.go import 函数中完成各个模块的 Controller 级别的路由注册

```go
package route

import (
	_ "github.com/hulklab/yago/example/app/modules/home/homecmd"
	_ "github.com/hulklab/yago/example/app/modules/home/homehttp"
	_ "github.com/hulklab/yago/example/app/modules/home/homerpc"
	_ "github.com/hulklab/yago/example/app/modules/home/hometask"
)

```
>注：使用 `yago new -m ${module}` 创建新模块时，会自动加载到 app/route/route.go 文件中，不需要手动添加。
