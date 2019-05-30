package homerpc

import (
	"context"
	"google.golang.org/grpc"
	"log"
	"os"
	"testing"
	"time"

	"github.com/hulklab/yago/example/app/modules/home/homerpc/protobuf/homepb"
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
	c := homepb.NewHomeClient(conn)

	// Contact the server and print out its response.
	name := defaultName
	if len(os.Args) > 1 {
		name = os.Args[1]
	}
	ctx, cancel := context.WithTimeout(context.Background(), 12*time.Second)
	defer cancel()
	r, err := c.Hello(ctx, &homepb.HelloRequest{Name: name})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Greeting: %s", r.Data)
}
