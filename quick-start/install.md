# 安装

首先你需要安装 [Go](https://golang.org/) (**version 1.13+**), 然后执行 go get 安装 yago

```bash
 go get github.com/hulklab/yago/yago
```

建议添加环境变量 GOPROXY=[https://goproxy.cn](https://github.com/goproxy/goproxy.cn/blob/master/README.zh-CN.md)，解决网络访问的问题

```bash
export GOPROXY=https://goproxy.cn,direct
```

## 依赖
go >= 1.13(由于使用了 go mod 管理版本依赖)

如果想在GOPATH下用mod, 请设置 GO111MODULE=on 则在 GOPATH/src 目录下使用 go get 时也默认采用 go mod

```bash
export GO111MODULE=on
```

## 环境变量参考

```bash
# go sdk
export GOROOT=/usr/local/go
# go path，根据各人习惯，一般都放在各人的home目录下面
export GOPATH=~/Workspace/go
# go 的可执行文件
export GOBIN=$GOPATH/bin
# 添加 go 的可执行文件到全局path 
export PATH=$PATH:$GOROOT/bin:$GOBIN
# go 1.13 开始支持 go get 时哪些仓库绕过代理，多用于私有仓库
export GOPRIVATE=*.private.repo
# go proxy 设置 go get 时的代理，direct 用来表示 go get 时如果遇到404，则直接走直连
export GOPROXY=https://goproxy.cn,direct
# 开启 go mod 管理依赖，默认为 auto
export GO111MODULE=on
```