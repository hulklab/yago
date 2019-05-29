
* 生成proto文件以及 server 和 client 文件
```
protoc -I protobuf homepb/home.proto --go_out=plugins=grpc:protobuf
```
* 修改服务端代码

* 原生 client 样例参考 app/modules/home/homerpc/home_test.go

* 封装 client 样例参考 app/third/homeapi/home.go::Hello