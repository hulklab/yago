package ghttp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hulklab/yago"
	"github.com/hulklab/yago/coms/logger"
	"github.com/sirupsen/logrus"
)

const (
	ctxParamsKey = "__PARAMS__"
)

func CheckUserName(c *yago.Ctx) {
	name := c.Param("name")
	if name == "devil" {
		c.SetError(yago.NewErr("path param name can not be devil"))
		c.Abort()
	}
}

func Compute(c *yago.Ctx) {
	// before request
	c.Set("number", 1)

	c.Next()

	// after request
	number := c.GetInt("number")

	c.SetData(fmt.Sprintf("the number is %d", number))
}

func BizLog(c *yago.Ctx) {

	setParam(c)

	c.Next()

	go func(c *yago.Ctx) {
		resp, ok := c.GetResponse()
		if !ok {
			return
		}
		var w http.Header
		params := c.GetString(ctxParamsKey)
		if resp.ErrNo != 0 {
			logger.Ins().Category("http.biz.error").WithFields(logrus.Fields{
				"url":             c.Request.URL.String(),
				"params":          params,
				"header":          c.Request.Header,
				"response_header": w,
				"user_ip":         c.ClientIP(),
			}).Error(c.GetError())
		} else {
			logger.Ins().Category("http.biz.info").WithFields(logrus.Fields{
				"url":             c.Request.URL.String(),
				"params":          params,
				"header":          c.Request.Header,
				"response_header": w,
				"user_ip":         c.ClientIP(),
				"resp":            resp,
			}).Debug()
		}
	}(c.Copy())
}

func setParam(c *yago.Ctx) {
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

}
