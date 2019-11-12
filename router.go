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
}

type CmdArg = CmdStringArg

type CmdStringArg struct {
	Name      string
	Shorthand string
	Usage     string
	Required  bool
	Value     string
}

func markFlagRequired(required bool, cmd *cobra.Command, name string) {
	if required {
		if err := cmd.MarkFlagRequired(name); err != nil {
			log.Printf("cmd arg %s mark flag failed: %s", name, err.Error())
		}
	}
}

func (c CmdStringArg) SetFlag(cmd *cobra.Command) {
	cmd.Flags().StringP(c.Name, c.Shorthand, c.Value, c.Usage)
	markFlagRequired(c.Required, cmd, c.Name)
}

type CmdStringSliceArg struct {
	Name      string
	Shorthand string
	Usage     string
	Required  bool
	Value     []string
}

func (c CmdStringSliceArg) SetFlag(cmd *cobra.Command) {
	cmd.Flags().StringSliceP(c.Name, c.Shorthand, c.Value, c.Usage)
	markFlagRequired(c.Required, cmd, c.Name)
}

type CmdBoolArg struct {
	Name      string
	Shorthand string
	Usage     string
	Required  bool
	Value     bool
}

func (c CmdBoolArg) SetFlag(cmd *cobra.Command) {
	cmd.Flags().BoolP(c.Name, c.Shorthand, c.Value, c.Usage)
	markFlagRequired(c.Required, cmd, c.Name)
}

type CmdIntArg struct {
	Name      string
	Shorthand string
	Usage     string
	Required  bool
	Value     int
}

func (c CmdIntArg) SetFlag(cmd *cobra.Command) {
	cmd.Flags().IntP(c.Name, c.Shorthand, c.Value, c.Usage)
	markFlagRequired(c.Required, cmd, c.Name)
}

type CmdIntSliceArg struct {
	Name      string
	Shorthand string
	Usage     string
	Required  bool
	Value     []int
}

func (c CmdIntSliceArg) SetFlag(cmd *cobra.Command) {
	cmd.Flags().IntSliceP(c.Name, c.Shorthand, c.Value, c.Usage)
	markFlagRequired(c.Required, cmd, c.Name)
}

type CmdInt64Arg struct {
	Name      string
	Shorthand string
	Usage     string
	Required  bool
	Value     int64
}

func (c CmdInt64Arg) SetFlag(cmd *cobra.Command) {
	cmd.Flags().Int64P(c.Name, c.Shorthand, c.Value, c.Usage)
	markFlagRequired(c.Required, cmd, c.Name)
}

type CmdDurationArg struct {
	Name      string
	Shorthand string
	Usage     string
	Required  bool
	Value     time.Duration
}

func (c CmdDurationArg) SetFlag(cmd *cobra.Command) {
	cmd.Flags().DurationP(c.Name, c.Shorthand, c.Value, c.Usage)
	markFlagRequired(c.Required, cmd, c.Name)
}

type CmdFloat64Arg struct {
	Name      string
	Shorthand string
	Usage     string
	Required  bool
	Value     float64
}

func (c CmdFloat64Arg) SetFlag(cmd *cobra.Command) {
	cmd.Flags().Float64P(c.Name, c.Shorthand, c.Value, c.Usage)
	markFlagRequired(c.Required, cmd, c.Name)
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
