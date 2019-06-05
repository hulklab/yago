package basehttp

import (
	"github.com/hulklab/yago"
	"github.com/hulklab/yago/libs/validator"
)

type BaseHttp struct {
}

func (h *BaseHttp) Rules() []validator.Rule {
	return nil
}

func (h *BaseHttp) Labels() validator.Label {
	return nil
}

func (h *BaseHttp) BeforeAction(c *yago.Ctx) yago.Err {
	return yago.OK
}
