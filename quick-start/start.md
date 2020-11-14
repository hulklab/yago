# 创建你的 Hello World

## 用 yago 在当前目录创建你的项目 my_app

```bash
yago init -a my_app
```

## 进入目录初始化 mod

```bash
cd my_app/
go mod init
```

## 构建

```bash
go build
```

## 创建属于自己的配置文件，并启动

```bash
sh env.init.sh yourname
./my_app
```

### 除此之外，我们推荐使用 yago 热编译工具 `yago run` 来开发调试。

控制台输出

```bash
# http 服务注册信息
2019/10/30 19:04:17 [HTTP] /home/detail my_app/app/modules/home/homehttp.(*HomeHttp).DetailAction-fm
[GIN-debug] GET    /home/detail              --> github.com/hulklab/yago.(*App).loadHttpRouter.func2 (5 handlers)
2019/10/30 19:04:17 [HTTP] /home/update my_app/app/modules/home/homehttp.(*HomeHttp).UpdateAction-fm
[GIN-debug] POST   /home/update              --> github.com/hulklab/yago.(*App).loadHttpRouter.func2 (5 handlers)
2019/10/30 19:04:17 [HTTP] /home/list my_app/app/modules/home/homehttp.(*HomeHttp).ListAction-fm
[GIN-debug] POST   /home/list                --> github.com/hulklab/yago.(*App).loadHttpRouter.func2 (5 handlers)
2019/10/30 19:04:17 [HTTP] /home/upload my_app/app/modules/home/homehttp.(*HomeHttp).UploadAction-fm
[GIN-debug] POST   /home/upload              --> github.com/hulklab/yago.(*App).loadHttpRouter.func2 (5 handlers)
2019/10/30 19:04:17 [HTTP] /home/hello my_app/app/modules/home/homehttp.(*HomeHttp).HelloAction-fm
[GIN-debug] GET    /home/hello               --> github.com/hulklab/yago.(*App).loadHttpRouter.func2 (5 handlers)
2019/10/30 19:04:17 [HTTP] /home/add my_app/app/modules/home/homehttp.(*HomeHttp).AddAction-fm
[GIN-debug] POST   /home/add                 --> github.com/hulklab/yago.(*App).loadHttpRouter.func2 (5 handlers)
2019/10/30 19:04:17 [HTTP] /home/delete my_app/app/modules/home/homehttp.(*HomeHttp).DeleteAction-fm
[GIN-debug] POST   /home/delete              --> github.com/hulklab/yago.(*App).loadHttpRouter.func2 (5 handlers)

# rpc 服务注册信息
2019/10/30 19:04:17 [GRPC] app.homepb.Home Hello

# task 服务注册信息
2019/10/30 19:04:17 [TASK] @loop my_app/app/modules/home/hometask.(*HomeTask).HelloLoopAction-fm
2019/10/30 19:04:17 [TASK] 0 */1 * * * * my_app/app/modules/home/hometask.(*HomeTask).HelloSchduleAction-fm

# demo task 任务日志打印
2019/10/30 19:04:17 Start Task homeTask.HelloLoopAction
2019/10/30 19:04:17 Doing Task homeTask.HelloLoopAction
2019/10/30 19:04:22 End Task homeTask.HelloLoopAction
```

出现以上信息则表示启动成功

## 测试程序 http 服务是否启动成功

```bash
curl "http://localhost:8080/user/hello?name=world"
{"errno":0,"errmsg":"","data":"hello world"}
```