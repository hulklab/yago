package ghttp

import (
	"github.com/hulklab/yago"
	"github.com/hulklab/yago/base/basehttp"
	"github.com/hulklab/yago/example/app/libs/trace"
)

type BaseHttp struct {
	basehttp.BaseHttp
}

func (h *BaseHttp) GetTraceCtx(c *yago.Ctx) *trace.Context {
	ctx := trace.NewWithCtx(c)
	return ctx
}

// 全局 routeGroup
var Root = yago.NewHttpGroupRouter("/", BizLog)
