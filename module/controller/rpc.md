### Rpc控制器

rpc控制器使用的是google的[grpc](https://grpc.io/docs/)框架，grpc框架使用了[protobuf](https://github.com/golang/protobuf)作为数据交换的格式，关于[grpc](https://grpc.io/docs/)，[protobuf](https://github.com/golang/protobuf)的使用可以点击链接查看官方文档，这里我们假设你已经安装好了所有的环境。

proto文件，home.proto
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

关于在Yago中proto文件的使用，请参考[protobuf规范](protobuf.md)，根据要求生成golang中间文件代码，以便在控制器中调用。

```go
package homerpc

import (
	"context"
	"github.com/hulklab/yago"
	"log"

	pb "github.com/hulklab/yago/example/app/modules/home/homerpc/homepb"
)

type HomeRpc struct {
}

func init() {
	homeRpc := new(HomeRpc)
	pb.RegisterHomeServer(yago.RpcServer, homeRpc)
}

func (r *HomeRpc) Hello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	log.Printf("Received: %v", in.Name)
	return &pb.HelloReply{Data: "Hello " + in.Name}, nil
}

```

rpc控制器中，我们按照proto文件定义的接口名称写下我们的Action代码并且在 init函数中完成路由的注册。

路由注册直接调用protoc生成的中间文件的函数RegisterHomeServer。

每个Action的参数和返回值都要遵循grpc的接口规范和proto的定义。