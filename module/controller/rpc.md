# RPC 控制器

RPC 控制器使用的是 google 的 [grpc](https://grpc.io/docs/) 框架，[grpc](https://grpc.io/docs/) 框架使用了 [protobuf](https://github.com/golang/protobuf) 作为数据交换的格式，关于 [grpc](https://grpc.io/docs/)，[protobuf](https://github.com/golang/protobuf) 的使用可以点击链接查看官方文档，这里我们假设你已经安装好了所有的环境。


## Protobuf

proto 文件 home.proto

```protobuf
syntax = "proto3";

package app.homepb;

service Home {
    rpc Hello (HelloRequest) returns (HelloReply) {}
}

message HelloRequest {
    string name = 1;
}

message HelloReply {
    string data = 1;
}
```

关于在 yago 中 proto 文件的使用，请参考 [Protobuf 规范](protobuf.md)，根据要求生成 go 中间文件代码，以便在控制器中调用。


## 路由注册

RPC 控制器中，我们按照 proto 文件定义的接口名称写下我们的 Action 代码并且在 init 函数中完成路由的注册。

路由注册直接调用 protoc 生成的中间文件的函数 RegisterHomeServer。

```go
func init() {
	homeRpc := new(HomeRpc)
	pb.RegisterHomeServer(yago.RpcServer, homeRpc)
}
```


## RPCAction

```go

func (r *HomeRpc) Hello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	log.Printf("Received: %v", in.Name)
	return &pb.HelloReply{Data: "Hello " + in.Name}, nil
}

```

每个 Action 的参数和返回值都要遵循 grpc 的接口规范和 proto 的定义即可。