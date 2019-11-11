# Yago程序的目录结构

```bash
tree
.
├── app
│   ├── g
│   │   ├── errors.go
│   │   └── type.go
│   ├── modules
│   │   └── home
│   │       ├── homecmd
│   │       │   └── home.go
│   │       ├── homedao
│   │       │   └── home.go
│   │       ├── homehttp
│   │       │   └── home.go
│   │       ├── homemodel
│   │       │   └── home.go
│   │       ├── homerpc
│   │       │   ├── README.md
│   │       │   ├── home.go
│   │       │   ├── home_test.go
│   │       │   └── homepb
│   │       │       ├── home.pb.go
│   │       │       └── home.proto
│   │       └── hometask
│   │           └── home.go
│   ├── route
│   │   └── route.go
│   └── third
│       └── homeapi
│           ├── home.go
│           └── homepb
│               ├── home.pb.go
│               └── home.proto
├── app.toml -> conf/app.toml
├── conf
│   └── app.toml
├── env.init.sh
├── go.mod
├── go.sum
├── logs
│   └── app.log
├── main.go
└── tools
    └── build.sh
```

| 路径 | 说明 |
| ---------- | --------- |
| app | 程序的主体代码目录 |
| app.toml | 程序的配置文件，通常为软链的conf目录中的某一个文件 |
| env.init.sh | 用于初始化环境和配置文件的脚本 |
| go.mod, go.sum | go mod 依赖所需文件 |
| logs | 默认的日志文件目录 |
| main.go | 程序的总入口 |
| tools | 存放脚本工具的目录 |
| app/g | 存放全局类型定义或变量的目录 |
| app/g/errors.go | 全局自定义错误信息及错误码 |
| app/g/type.go | 全局自定义类型 |
| app/modules | 模块目录，yago的程序都是按模块来编写的 |
| app/modules/home | 默认的home模块目录 |
| app/modules/homecmd | home模块的cmd（命令行）服务控制器目录 |
| app/modules/homehttp | home模块的http（web）服务控制器目录 |
| app/modules/homerpc | home模块的rpc（远程调用）服务控制器目录 |
| app/modules/hometask | home模块的task（定时任务，常驻进程）服务控制器目录 |
| app/modules/homedao | home模块的数据库映射和操作目录 |
| app/modules/homemodel | home模块的主要业务逻辑处理目录 |
| app/route | 路由注册开关目录 |
| app/third | 第三方调用api目录 |