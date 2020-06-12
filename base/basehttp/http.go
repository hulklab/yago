package basehttp

import (
	"github.com/gin-gonic/gin/binding"
)

type BaseHttp struct{}

func init() {
	binding.Validator = &defaultValidator{}
}
