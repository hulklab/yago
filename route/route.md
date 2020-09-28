# 路由 Router

## Yago 路由注册分为两个阶段：

* step1. 在各个控制器的 init 函数中完成一个 Controller 内具体的 Action 级别的路由注册，详细使用请参考 [模块控制器](/module/controller.md)

```go
func init() {
	userHttp := new(UserHttp)
    
    userGroup := yago.NewHttpGroupRouter("/user") // 创建路由分组
    userGroup.Get("/hello", userHttp.HelloAction)
    userGroup.Post("/add", userHttp.AddAction)
    userGroup.Post("/delete", userHttp.DeleteAction)
    userGroup.Get("/detail", userHttp.DetailAction)
    userGroup.Post("/update", userHttp.UpdateAction)
    userGroup.Post("/list", userHttp.ListAction)
    userGroup.Post("/base-list", userHttp.BaseListAction)
    userGroup.Post("/upload", userHttp.UploadAction)
    userGroup.Get("/user/:name", userHttp.Hello2Action)
    userGroup.Get("/cookie", userHttp.CookieAction)
    userGroup.Get("/metadata", userHttp.MetadataAction).WithMetadata(HttpMetadata{
        Label: "自定义HTTP名称",
    }) // 注册 API metadata 信息

    memberGroup := yago.NewHttpGroupRouter("/user/member", homemiddleware.CheckUserName) // 对路由分组使用中间件
    {
        memberGroup.Post("/:name", userHttp.UserSetAction)
        memberGroup.Get("/:name", userHttp.UserGetAction)
        memberGroup.Put("/:name", userHttp.UserUpdateAction)
        memberGroup.Delete("/:name", userHttp.UserDeleteAction)

        consumeSubGroup := memberGroup.Group("/plus")
        consumeSubGroup.Patch("/number/:number", homemiddleware.Compute, userHttp.PlusAction) // 对单个 API 使用中间件
    }

    yago.SetHttpNoRouter(userHttp.NoRouterAction) // 注册 404 页面
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
