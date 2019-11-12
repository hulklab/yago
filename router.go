package yago

import (
	"log"
	"strings"
	"time"

	"github.com/hulklab/yago/libs/validator"
	"github.com/spf13/cobra"
)

// http
type HttpHandlerFunc func(c *Ctx)

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

// task
type TaskHandlerFunc func()

type TaskRouter struct {
	Spec   string
	Action TaskHandlerFunc
}

var TaskRouterList []*TaskRouter

func AddTaskRouter(spec string, action TaskHandlerFunc) {
	TaskRouterList = append(TaskRouterList, &TaskRouter{spec, action})
}

// cmd
type CmdHandlerFunc func(cmd *cobra.Command, args []string)

type ICmdArg interface {
	SetFlag(cmd *cobra.Command)
	MarkRequired(cmd *cobra.Command, use string)
}

type baseCmdArg struct {
	Name      string
	Shorthand string
	Usage     string
	Required  bool
}

func (c baseCmdArg) MarkRequired(cmd *cobra.Command, use string) {
	if c.Required {
		if err := cmd.MarkFlagRequired(c.Name); err != nil {
			log.Printf("cmd %s mark flag failed: %s", use, err.Error())
		}
	}
}

type CmdArg = CmdStringArg

type CmdStringArg struct {
	baseCmdArg
	Value string
}

func (c CmdStringArg) SetFlag(cmd *cobra.Command) {
	cmd.Flags().StringP(c.Name, c.Shorthand, c.Value, c.Usage)
}

type CmdIntArg struct {
	baseCmdArg
	Value int
}

func (c *CmdIntArg) SetFlag(cmd *cobra.Command) {
	cmd.Flags().IntP(c.Name, c.Shorthand, c.Value, c.Usage)
}

type CmdInt64Arg struct {
	baseCmdArg
	Value int64
}

func (c *CmdInt64Arg) SetFlag(cmd *cobra.Command) {
	cmd.Flags().Int64P(c.Name, c.Shorthand, c.Value, c.Usage)
}

type CmdDurationArg struct {
	baseCmdArg
	Value time.Duration
}

func (c *CmdDurationArg) SetFlag(cmd *cobra.Command) {
	cmd.Flags().DurationP(c.Name, c.Shorthand, c.Value, c.Usage)
}

type CmdFloat64Arg struct {
	baseCmdArg
	Value float64
}

func (c *CmdFloat64Arg) SetFlag(cmd *cobra.Command) {
	cmd.Flags().Float64P(c.Name, c.Shorthand, c.Value, c.Usage)
}

type CmdRouter struct {
	Use    string
	Short  string
	Action CmdHandlerFunc
	Args   []ICmdArg
}

var CmdRouterMap = make(map[string]*CmdRouter)

func AddCmdRouter(use, short string, action CmdHandlerFunc, args ...ICmdArg) {
	cmdSlice := strings.Split(use, "/")
	if len(cmdSlice) == 0 {
		return
	}

	if _, ok := CmdRouterMap[use]; ok {
		log.Panicf("http router duplicate : %s", use)
	}

	CmdRouterMap[use] = &CmdRouter{use, short, action, args}
}
