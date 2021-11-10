#### 生成 proto 文件以及 server 和 client 文件

```
cd app/modules/demo/demorpc

protoc -I demopb home.proto --go_out=plugins=grpc:demopb
```

#### protoc 命令说明

* -I 参数指定 proto 文件所在的包，上面命令中会去 demopb 目录搜索 home.proto 文件
* --go_out 参数里冒号后面执行的目录，为生成的 go 文件目录放置的位置，上面命令中会将生成的 home.pb.go 放入 demopb 目录中

#### pb 规范说明：

* proto 文件的包名：pb 中的包名是全局的，我们建议采用 ${app}.${module}pb 的规则来命令，例如 home.proto 的 package 为 app.demopb
* 每个 module 的 rpc 目录下创建有且一个 ${module}pb 的目录，如：demopb
* 每个 service(可以理解为 controller) 定义一个 proto 文件，如 home.proto，不同的 service 的 proto 文件统一放入当前模块的 pb 目录中
* 每个 service 在 rpc 目录下有一个对应的 .go 文件，如：home.go

#### 样例参考

* 原生 client 样例参考 app/modules/demo/demorpc/home_test.go
* 封装 client 样例参考 app/third/homeapi/home.go::Hello

#### SSL

```
// 服务端生成证书
// openssl genrsa -out server.key 2048
// openssl req -new -x509 -sha256 -key server.key -out server.pem -days 3650
```
