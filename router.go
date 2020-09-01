package yago

import (
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
)

// http
type HttpHandlerFunc func(c *Ctx)

type HttpRouter struct {
	Group    *HttpGroupRouter
	Path     string
	Method   string
	Actions  []HttpHandlerFunc
	Metadata interface{}
}

func (h *HttpRouter) WithMetadata(md interface{}) *HttpRouter {
	h.Metadata = md
	return h
}

func SetHttpNoRouter(action HttpHandlerFunc) {
	httpNoRouterHandler = action
}

type HttpGlobalMiddleware []HttpHandlerFunc

func HttpUse(middleware ...HttpHandlerFunc) {
	httpGlobalMiddleware.Use(middleware...)
}

func (r *HttpGlobalMiddleware) Use(middleware ...HttpHandlerFunc) {
	*r = append(*r, middleware...)
}

var (
	httpGlobalMiddleware HttpGlobalMiddleware
	httpGroupRouterMap   = make(map[string]*HttpGroupRouter)
	httpNoRouterHandler  HttpHandlerFunc
)

func (h *HttpRouter) Url() string {
	url := h.Path
	p := h.Group
	for p != nil {
		if p.Prefix != "/" {
			url = p.Prefix + url
		}
		p = p.Parent
	}
	return url
}

// http group router
type HttpGroupRouter struct {
	Prefix         string
	GinGroup       *gin.RouterGroup
	Middlewares    []HttpHandlerFunc
	HttpRouterList []*HttpRouter
	Parent         *HttpGroupRouter
	Children       map[string]*HttpGroupRouter
}

func NewHttpGroupRouter(prefix string) *HttpGroupRouter {
	if len(prefix) == 0 {
		log.Panic("http group router name can not be empty")
	}

	if group, ok := httpGroupRouterMap[prefix]; ok {
		return group
	}

	httpGroupRouterMap[prefix] = &HttpGroupRouter{
		Prefix: prefix,
	}

	return httpGroupRouterMap[prefix]
}

func (g *HttpGroupRouter) Group(prefix string) *HttpGroupRouter {
	if len(prefix) == 0 {
		log.Panic("http sub group router name can not be empty")
	}

	if group, ok := g.Children[prefix]; ok {
		return group
	}

	group := &HttpGroupRouter{
		Prefix: prefix,
		Parent: g,
	}

	if g.Children == nil {
		g.Children = make(map[string]*HttpGroupRouter)
	}

	g.Children[prefix] = group

	return group
}

func getHttpRoutersWithSubGroup(g *HttpGroupRouter) []*HttpRouter {
	var subRouterList []*HttpRouter
	subRouterList = append(subRouterList, g.HttpRouterList...)
	if len(g.Children) > 0 {
		for _, sub := range g.Children {
			subRouterList = append(subRouterList, getHttpRoutersWithSubGroup(sub)...)
		}
	}
	return subRouterList
}

func GetHttpRouters() []*HttpRouter {
	var routerList []*HttpRouter
	for _, v := range httpGroupRouterMap {
		routerList = append(routerList, getHttpRoutersWithSubGroup(v)...)
	}
	return routerList
}

func (g *HttpGroupRouter) Use(middlewares ...HttpHandlerFunc) {
	g.Middlewares = append(g.Middlewares, middlewares...)
}

func (g *HttpGroupRouter) addHttpRouter(url, method string, actions ...HttpHandlerFunc) *HttpRouter {

	router := &HttpRouter{
		Path:    url,
		Method:  method,
		Group:   g,
		Actions: actions,
	}

	g.HttpRouterList = append(g.HttpRouterList, router)

	return router
}

func (g *HttpGroupRouter) Get(url string, actions ...HttpHandlerFunc) *HttpRouter {
	return g.addHttpRouter(url, http.MethodGet, actions...)
}

func (g *HttpGroupRouter) Post(url string, actions ...HttpHandlerFunc) *HttpRouter {
	return g.addHttpRouter(url, http.MethodPost, actions...)
}

func (g *HttpGroupRouter) Put(url string, actions ...HttpHandlerFunc) *HttpRouter {
	return g.addHttpRouter(url, http.MethodPut, actions...)
}

func (g *HttpGroupRouter) Delete(url string, actions ...HttpHandlerFunc) *HttpRouter {
	return g.addHttpRouter(url, http.MethodDelete, actions...)
}

func (g *HttpGroupRouter) Patch(url string, actions ...HttpHandlerFunc) *HttpRouter {
	return g.addHttpRouter(url, http.MethodPatch, actions...)
}

func (g *HttpGroupRouter) Head(url string, actions ...HttpHandlerFunc) *HttpRouter {
	return g.addHttpRouter(url, http.MethodHead, actions...)
}

func (g *HttpGroupRouter) Options(url string, actions ...HttpHandlerFunc) *HttpRouter {
	return g.addHttpRouter(url, http.MethodOptions, actions...)
}

func (g *HttpGroupRouter) Any(url string, actions ...HttpHandlerFunc) *HttpRouter {
	return g.addHttpRouter(url, "Any", actions...)
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
