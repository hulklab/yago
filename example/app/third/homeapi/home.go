package homeapi

import (
	"context"
	"fmt"
	"io"

	"google.golang.org/grpc"

	"github.com/hulklab/yago"
	"github.com/hulklab/yago/base/basethird"
	pb "github.com/hulklab/yago/example/app/third/homeapi/homepb"
)

type HomeApi struct {
	basethird.HttpThird
	basethird.RpcThird
}

// Usage: New().GetUserById()
func New() *HomeApi {

	api := new(HomeApi)

	// http 配置
	api.Domain = yago.Config.GetString("home_api.domain")
	api.HttpThird.Hostname = yago.Config.GetString("home_api.hostname")

	// rpc 配置
	api.Address = yago.Config.GetString("home_api.rpc_address")
	api.SslOn = yago.Config.GetBool("home_api.ssl_on")
	api.CertFile = yago.Config.GetString("home_api.cert_file")
	api.RpcThird.Hostname = yago.Config.GetString("home_api.hostname")
	api.Timeout = yago.Config.GetInt("home_api.timeout")
	api.MaxRecvMsgsizeMb = yago.Config.GetInt("home_api.max_recv_msgsize_mb")
	api.MaxSendMsgsizeMb = yago.Config.GetInt("home_api.max_send_msgsize_mb")

	return api
}

func (a *HomeApi) GetUserById(id int64) string {

	params := map[string]interface{}{
		"id": id,
	}

	req, err := a.Get("/home/user/detail", params)
	if err != nil {
		return err.Error()
	} else {
		s, _ := req.String()
		return s
	}
}

func (a *HomeApi) UploadFile(filepath string) string {

	params := map[string]interface{}{
		"file": basethird.PostFile(filepath),
	}

	req, err := a.Post("/home/user/upload", params)
	if err != nil {
		return err.Error()
	} else {
		s, _ := req.String()
		return s
	}
}

// eg. homeapi.New().RpcHello()
func (a *HomeApi) RpcHello() {
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

func (a *HomeApi) RpcHelloStream() {

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
