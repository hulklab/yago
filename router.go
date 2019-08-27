package yago

import (
	"github.com/hulklab/yago/libs/validator"
	"github.com/spf13/cobra"
	"log"
	"strings"
)

type HttpHandlerFunc func(c *Ctx)

// http
type HttpRouter struct {
	Url    string
	Method string
	Action HttpHandlerFunc
	h      HttpInterface
}

var HttpRouterMap = make(map[string]*HttpRouter)

type HttpInterface interface {
	Labels() validator.Label
	Rules() []validator.Rule
	BeforeAction(c *Ctx) Err
	AfterAction(c *Ctx)
}

func AddHttpRouter(url, method string, action HttpHandlerFunc, h HttpInterface) {
	if _, ok := HttpRouterMap[url]; ok {
		log.Panicf("http router duplicate : %s", url)
	}
	HttpRouterMap[url] = &HttpRouter{url, method, action, h}
}

//  task
type TaskHandlerFunc func()

type TaskRouter struct {
	Spec   string
	Action TaskHandlerFunc
}

var TaskRouterList []*TaskRouter

func AddTaskRouter(spec string, action TaskHandlerFunc) {
	TaskRouterList = append(TaskRouterList, &TaskRouter{spec, action})
}

//  cmd
type CmdHandlerFunc func(cmd *cobra.Command, args []string)

type CmdArg struct {
	Name      string
	Shorthand string
	Value     string
	Usage     string
	Required  bool
}

type CmdRouter struct {
	Use    string
	Short  string
	Action CmdHandlerFunc
	Args   []CmdArg
}

var CmdRouterMap = make(map[string]*CmdRouter)

func AddCmdRouter(use, short string, action CmdHandlerFunc, args ...CmdArg) {
	cmdSlice := strings.Split(use, "/")
	if len(cmdSlice) == 0 {
		return
	}

	if _, ok := CmdRouterMap[use]; ok {
		log.Panicf("http router duplicate : %s", use)
	}

	CmdRouterMap[use] = &CmdRouter{use, short, action, args}
}
