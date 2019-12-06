package homeapi

import (
	"context"
	"fmt"
	"io"
	"log"

	"github.com/hulklab/yago"

	"google.golang.org/grpc"

	"github.com/hulklab/yago/base/basethird"
	pb "github.com/hulklab/yago/example/app/third/homeapi/homepb"
)

type homeRpcApi struct {
	basethird.RpcThird
}

// Usage: InsRpc().Hello()
func InsRpc() *homeRpcApi {
	name := "home_rpc_api"
	v := yago.Component.Ins(name, func() interface{} {
		api := new(homeRpcApi)

		// rpc 配置
		err := api.InitConfig(name)
		if err != nil {
			log.Fatal("init rpc api config error")
		}
		return api
	})
	return v.(*homeRpcApi)
}

// eg. homeapi.InsRpc().Hello()
func (a *homeRpcApi) Hello() {
	a.SetBeforeUnaryClientInterceptor(func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, opts ...grpc.CallOption) error {
		fmt.Println("method:", method, "before")
		return nil
	})
	a.SetAfterUnaryClientInterceptor(func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, opts ...grpc.CallOption) error {
		fmt.Println("method:", method, "after")
		return nil
	})

	var name = "zhangsan"

	req := &pb.HelloRequest{Name: name}

	conn, _ := a.GetConn()
	ctx, cancel := a.GetCtx()
	defer cancel()

	c := pb.NewHomeClient(conn)
	resp, err := c.Hello(ctx, req)

	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("ok:", resp.Data)
}

func (a *homeRpcApi) HelloStream() {

	var name = "zhangsan"
	req := &pb.HelloRequest{Name: name}

	a.SetBeforeStreamClientInterceptor(func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, opts ...grpc.CallOption) error {
		fmt.Println("method:", method, "stop")
		return nil
	})

	conn, _ := a.GetConn()
	ctx, cancel := a.GetCtx()
	defer cancel()

	c := pb.NewHomeClient(conn)
	stream, err := c.HelloStream(ctx, req)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	for {
		reply, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println(err)
			break
		}

		fmt.Println(reply)
	}
}
