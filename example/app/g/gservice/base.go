package gservice

import (
	"github.com/hulklab/yago/example/app/libs/trace"
	"github.com/sirupsen/logrus"
)

type BaseService struct {
	Ctx *trace.Context
}

func (s *BaseService) Init(ctx *trace.Context) {
	s.Ctx = ctx
}

func (s *BaseService) Logger() *logrus.Entry {
	return s.Ctx.Logger()
}
