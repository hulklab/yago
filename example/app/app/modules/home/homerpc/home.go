package homerpc

import (
	"context"
	"github.com/hulklab/yago"
	"log"

	pb "github.com/hulklab/yago/example/app/app/modules/home/homerpc/protobuf/homepb"
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
