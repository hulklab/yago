# Yago Scafford

![avatar](http://p0.qhimg.com/t0162ed78090852688f.png)

### Yago web开发脚手架

## 目录

- [安装](#安装)
- [依赖](#依赖)
- [快速开始](#快速开始)
- [使用说明](#使用说明)
	- [文件结构](#文件结构)
	- [路由注册](#路由注册)
	- [配置](#配置)
	- [组件](#组件)
	- [模块](#模块)
	- [错误处理](#错误处理)
	- [第三方API调用](#第三方API调用)
	- [Goland配置mod](#Goland配置mod)
- [常用类库推荐](#常用类库推荐])
- [已知问题及解决方案](#已知问题及解决方案)
- [感谢](#感谢)

## 安装

##### 首先你需要安装 [Go](https://golang.org/) (**version 1.11+**), 然后执行go get 安装yago

```bash
 go get github.com/hulklab/yago/yago
```
## 依赖
go >= 1.11(由于使用了 go mod 管理版本依赖)

##### 如果想在GOPATH下用mod, 请设置 GO111MODULE=on 则在 GOPATH/src 目录下使用 go get 时也默认采用 go mod
```bash
export GO111MODULE=on
```

## 快速开始

##### 1. 用 yago 在当前目录创建你的项目 myapp
```bash
yago init -a myapp
```

##### 2. 进入目录初始化

```bash
cd myapp/
go mod init
```

##### 3. 构建
```bash
go build
```
> [如 go build 遇报错，请看解决方案](#已知问题及解决方案)

##### 4. 创建属于自己的配置文件，并启动
```bash
sh env.init.sh yourname
./myapp
```

##### 5. 控制是否需要在此机器上开启 task 任务，有两种方式

* 修改配置文件中的 app.task_enable，默认为开启
* 修改环境变量 export {{配置文件中的app_name}}_APP_TASK_ENABLE=1, 1 表示开启，0 表示关闭，配置文件与环境变量同时存在时环境变量生效

## 使用说明

### 文件结构
```
├── app
│   ├── g
│   │   └── errors.go //系统级错误定义处
│   ├── modules //模块目录
│   │   └── home // 样例模块
│   │       ├── homecmd // 命令行 控制器
│   │       │   └── home.go
│   │       ├── homedao // 数据库访问层
│   │       │   └── home.go
│   │       ├── homehttp // http 控制器
│   │       │   └── home.go
│   │       ├── homemodel // 业务逻辑层
│   │       │   └── home.go
│   │       ├── homerpc // grpc 控制器
│   │       │   ├── home.go
│   │       │   ├── home_test.go
│   │       │   │   └── homepb
│   │       │   │       ├── home.pb.go
│   │       │   │       └── home.proto // proto文件
│   │       │   └── README.md
│   │       └── hometask // 常驻进程和定时任务控制器
│   │           └── home.go
│   ├── route // 路由管理目录
│   │   ├── route.go // 路由控制文件
│   └── third // 第三方api调用目录
│       └── homeapi 
│           ├── home.go // http rpc 客户端接口封装
│           └── protobuf
│               └── homepb
│                   ├── home.pb.go
│                   └── home.proto
├── conf // 配置文件目录
│   └── app.toml
├── main.go // 程序总入口
└── tools // 构建工具
    └── build.sh
```

### 路由注册

yago 提供了四种控制器入口，分别是 命令行(cmd)，http，grpc, 常驻进程和定时任务(task)。各个控制器的路由注册就在各个控制器层的 init函数中完成。

请分别参考样例
#### 1. http 路由 

@reference example/app/modules/home/homehttp/home.go


#### 2. cmd 路由

@reference example/app/modules/home/homecmd/home.go


#### 3. task 路由

@reference example/app/modules/home/hometask/home.go

#### 4. grpc 路由

@reference example/app/modules/home/homerpc/home.go


### 配置

配置文件解析完成后存储在全局yago.Config中，yago.Config是 https://github.com/spf13/viper 的扩展。原生采用viper的方法来获取配置文件的值即可。

### 组件

yago集成的组件包括 redis，mysql orm，logger，kafka，组建统一实现了Ins 单例方法，需要使用组件可以参考各组件的test文件。

### 模块

##### 1. 创建新模块

在项目根目录下，使用yago创建模块。新创建的模块会自动将路由加载到myapp/routes中，模块内容的编写可以参考样例 home模块

```
cd myapp
yago new -m newmodule
```

### 错误处理

1. 系统级错误定义处 `./app/g/error.go`
2. 使用 `@reference example/app/modules/homehttp/home.go::AddAction`

### 第三方API调用

第三方API调用目前支持 http和grpc

1. 目录规范 `@see example/app/third`
2. 使用样例 `@reference example/app/third/homeapi/home.go`

### Goland配置mod

1. Preferences -> Go -> Go modules(vgo)

![](http://p406.qhimgs4.com/t0100eba6c9f82cb921.png)

2. 如果还有标红的提示，点击 Sync packages

![](http://p406.qhimgs4.com/t019f0fcae328f7a0e0.png)

## 常用类库推荐
1. 内存 cache 
    >https://github.com/patrickmn/go-cache
2. LRU cache
    >https://github.com/golang/groupcache
3. Converts a mysql table into a golang struct
    >https://github.com/Shelnutt2/db2struct
4. Struct copy to struct, map convert into struct
    >https://github.com/mitchellh/mapstructure
5. MessagePack encoding for Golang
    >https://github.com/vmihailenco/msgpack

## 已知问题及解决方案

1.  unknown import path "github.com/ugorji/go/codec": ambiguous import: found github.com/ugorji/go/codec in multiple modules 模块冲突问题

	原因参考：
	
	> https://cloud.tencent.com/developer/article/1417112
	
	解决方案：
	在go.mod文件最下面添加如下代码
	```go
	replace github.com/ugorji/go v1.1.4 => github.com/ugorji/go/codec v0.0.0-20190204201341-e444a5086c43
	```


## 感谢
[gin](https://github.com/gin-gonic/gin)
[cron](https://github.com/robfig/cron)
[cobra](https://github.com/spf13/cobra)
[xorm](http://github.com/go-xorm/xorm)
[logrus](https://github.com/sirupsen/logrus)
[beego](https://github.com/astaxie/beego)