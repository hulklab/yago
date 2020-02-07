package main

var HttpTemplate = `package {{PACKAGE}} 

import (
	"net/http"

	"github.com/hulklab/yago"
	"github.com/hulklab/yago/base/basehttp"
)

type {{NAME}}Http struct {
	basehttp.BaseHttp
}

func init() {
	h := new({{NAME}}Http)
	yago.AddHttpRouter("", http.MethodPost, h.ListAction, h)
}

func (h *{{NAME}}Http) ListAction(c *yago.Ctx) {
	return
}
`

var RpcTemplate = `package {{PACKAGE}} 

import (
	"context"
	"fmt"
	"log"

	"github.com/hulklab/yago"

	pb "github.com/hulklab/yago/example/app/modules/home/homerpc/homepb"
)

type {{NAME}}Rpc struct {
}

func init() {
	h := new({{NAME}}Rpc)
	pb.RegisterHomeServer(yago.RpcServer, h)
}

func (r *{{NAME}}Rpc) Hello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	log.Printf("Received: %v", in.Name)
	return &pb.HelloReply{Data: "Hello " + in.Name}, nil
}
`

var CmdTemplate = `package {{PACKAGE}} 

import (
	"github.com/hulklab/yago"
	"github.com/hulklab/yago/base/basecmd"
	"github.com/spf13/cobra"
)

type {{NAME}}Cmd struct {
	basecmd.BaseCmd
}

func init() {
	c := new({{NAME}}Cmd)
	// 注册路由
	yago.AddCmdRouter("demo", "Demo action", c.DemoAction, yago.CmdStringArg{
		Name: "arg", Shorthand: "a", Value: "value", Usage: "参数", Required: true,
	})
}

func (c *{{NAME}}Cmd) DemoAction(cmd *cobra.Command, args []string) {

}

`

var TaskTemplate = `package {{PACKAGE}} 

import (
	"github.com/hulklab/yago"
	"github.com/hulklab/yago/base/basetask"
)

type {{NAME}}Task struct {
	basetask.BaseTask
}

func init() {
	t := new({{NAME}}Task)
	yago.AddTaskRouter("@loop", t.HelloAction)
	yago.AddTaskRouter("0 */1 * * * *", t.HelloAction)
}

func (t *{{NAME}}Task) HelloAction() {
	//t.RunLoop(func() {
	//})
}
`

var ModelTemplate = `package {{PACKAGE}} 

import (
	"github.com/hulklab/yago"
)

type {{LNAME}}Model struct {
}

func New{{NAME}}Model() *{{LNAME}}Model {
	return &{{LNAME}}Model{}
}

`
