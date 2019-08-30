package basehttp

import (
	"github.com/hulklab/yago"
	"github.com/hulklab/yago/coms/logger"
	"github.com/hulklab/yago/libs/validator"
	"github.com/sirupsen/logrus"
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

func (h *BaseHttp) AfterAction(c *yago.Ctx) {
	resp, ok := c.GetResponse()
	if !ok {
		return
	}

	if resp.ErrNo != 0 {
		logger.Ins().Category("http.biz.error").WithFields(logrus.Fields{
			"url":     c.Request.URL.Path,
			"params":  c.Keys,
			"user_ip": c.ClientIP(),
		}).Error(resp.ErrMsg)

	} else {
		logger.Ins().Category("http.biz.info").WithFields(logrus.Fields{
			"url":     c.Request.URL.Path,
			"params":  c.Keys,
			"user_ip": c.ClientIP(),
			"resp":    resp,
		}).Debug()

	}
}
