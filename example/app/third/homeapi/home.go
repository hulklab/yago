package homeapi

import (
	"context"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/hulklab/yago"
	"github.com/hulklab/yago/common/basethird"
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
	api.Hostname = yago.Config.GetString("home_api.hostname")

	// rpc 配置
	api.Address = yago.Config.GetString("home_api.rpc_address")
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

	rep, err := a.Call(func(conn *grpc.ClientConn, ctx context.Context) (rep proto.Message, e error) {

		c := pb.NewHomeClient(conn)

		return c.Hello(ctx, &pb.HelloRequest{Name: name})

	}, name)

	if err != nil {
		fmt.Println(err)
		return
	}

	v, ok := rep.(*pb.HelloReply)
	if ok {
		fmt.Println("ok:", v.Data)
	} else {
		fmt.Println("not match", v)
	}
}
