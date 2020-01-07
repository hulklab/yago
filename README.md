# Yago Scafford - Web开发脚手架

![avatar](http://p0.qhimg.com/t0162ed78090852688f.png) 

## 目录

- [文档](#文档)
- [安装](#安装)
- [依赖](#依赖)
- [快速开始](#快速开始)
- [已知问题及解决方案](#已知问题及解决方案)
- [感谢](#感谢)

## 文档
[yago 文档](https://hulklab.github.io/yago/)

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



### Goland 配置 mod

1. Preferences -> Go -> Go modules(vgo)

![](http://p406.qhimgs4.com/t0100eba6c9f82cb921.png)

2. 如果还有标红的提示，点击 Sync packages

![](http://p406.qhimgs4.com/t019f0fcae328f7a0e0.png)


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
