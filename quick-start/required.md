# 依赖
```
go >= 1.11 （由于使用了 go mod 管理版本依赖）
```

> 如果想在 GOPATH 下用 mod，请设置 GO111MODULE=on 则在 GOPATH/src 目录下使用 go get 时也默认采用 go mod

```bash
export GO111MODULE=on (通常我们会将这条命令加入到 ~/.bashrc 文件中)
```
