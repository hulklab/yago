package basehttp

import (
	"github.com/gin-gonic/gin/binding"
	"github.com/hulklab/yago"
	"github.com/hulklab/yago/base/basemiddleware"
)

type BaseHttp struct{}

func init() {
	binding.Validator = &defaultValidator{}
	yago.GetHttpGlobalMiddleware().Use(basemiddleware.BizLog)
}
