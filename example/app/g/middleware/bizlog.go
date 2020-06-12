package middleware

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/hulklab/yago"
	"github.com/hulklab/yago/coms/logger"
	"github.com/sirupsen/logrus"
)

const (
	ctxParamsKey = "__PARAMS__"
)

func BizLog(c *yago.Ctx) {

	req := c.Request
	var paramKey = ctxParamsKey
	c.Set(paramKey, "")

	switch c.ContentType() {
	case gin.MIMEJSON:
		bodyBytes, err := ioutil.ReadAll(req.Body)
		if err != nil {
			log.Println("read body:", err.Error())
			return
		}

		err = req.Body.Close() //  must close
		if err != nil {
			log.Println("close body:", err.Error())
			return
		}
		req.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

		c.Set(paramKey, string(bodyBytes))
	case gin.MIMEPOSTForm:
		err := req.ParseForm()
		if err != nil {
			log.Println("parse form", err.Error())
			return
		}
		bs, err := json.Marshal(req.PostForm)
		if err != nil {
			log.Println("json encode err:", err.Error())
		}

		c.Set(paramKey, string(bs))

	case gin.MIMEMultipartPOSTForm:
		err := req.ParseMultipartForm(32 << 20)
		if err != nil {
			log.Println("parse multi form", err.Error())
			return
		} else if req.MultipartForm != nil {
			bs, err := json.Marshal(req.PostForm)
			if err != nil {
				log.Println("json encode err:", err.Error())
			}
			c.Set(paramKey, string(bs))
		}
	}

	c.Next()

	resp, ok := c.GetResponse()
	if !ok {
		return
	}

	params := c.GetString(ctxParamsKey)

	if resp.ErrNo != 0 {
		logger.Ins().Category("http.biz.error").WithFields(logrus.Fields{
			"url":     c.Request.URL.String(),
			"params":  params,
			"header":  c.Request.Header,
			"user_ip": c.ClientIP(),
		}).Error(c.GetError())
	} else {
		logger.Ins().Category("http.biz.info").WithFields(logrus.Fields{
			"url":     c.Request.URL.String(),
			"params":  params,
			"header":  c.Request.Header,
			"user_ip": c.ClientIP(),
			"resp":    resp,
		}).Debug()
	}
}
