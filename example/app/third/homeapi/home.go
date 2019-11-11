package homeapi

import (
	"context"
	"errors"
	"fmt"

	"github.com/golang/protobuf/proto"
	"github.com/hulklab/yago"
	"github.com/hulklab/yago/base/basethird"
	"google.golang.org/grpc"

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

	var name = "zhangsan"

	req := &pb.HelloRequest{Name: name}
	resp, err := a.Call(func(conn *grpc.ClientConn, ctx context.Context, req proto.Message) (resp proto.Message, e error) {

		c := pb.NewHomeClient(conn)

		v, ok := req.(*pb.HelloRequest)

		if ok {
			return c.Hello(ctx, v)
		}
		return nil, errors.New("req is not the type of HelloRequest")
	}, req)

	if err != nil {
		fmt.Println(err)
		return
	}

	v, ok := resp.(*pb.HelloReply)
	if ok {
		fmt.Println("ok:", v.Data)
	} else {
		fmt.Println("not match", v)
	}
}
