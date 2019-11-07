## 路由 Router

Yago路由注册分为两个阶段：

1. 在各个控制器的init函数中完成一个Controller内具体的Action级别的路由注册

2. 在app/route/route.go import函数中完成各个模块的Controller级别的路由注册
