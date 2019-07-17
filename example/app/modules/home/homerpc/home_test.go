package homerpc

import (
	"context"
	"google.golang.org/grpc"
	"log"
	"testing"
	"time"

	pb "github.com/hulklab/yago/example/app/modules/home/homerpc/homepb"
)

const (
	address     = ":50051"
	defaultName = "world"
)

func TestHomeRpc_Hello(t *testing.T) {
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewHomeClient(conn)

	// Contact the server and print out its response.
	name := defaultName
	ctx, cancel := context.WithTimeout(context.Background(), 12*time.Second)
	defer cancel()
	r, err := c.Hello(ctx, &pb.HelloRequest{Name: name})
	if err != nil {
		t.Fatalf("could not greet: %v", err)
	}
	t.Logf("Greeting: %s", r.Data)
}
