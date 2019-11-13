package homerpc

import (
	"context"
	"fmt"
	"log"

	"github.com/hulklab/yago"

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

func (r *HomeRpc) HelloStream(in *pb.HelloRequest, srv pb.Home_HelloStreamServer) error {
	for i := 0; i < 10; i++ {
		name := fmt.Sprintf("%d: %s", i, in.GetName())
		reply := &pb.HelloStreamReply{
			Data: name,
		}
		err := srv.Send(reply)
		if err != nil {
			return err
		}
	}
	return nil
}
