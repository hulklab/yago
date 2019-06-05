package basethird

import (
	"errors"
	"fmt"
	"github.com/astaxie/beego/httplib"
	"github.com/hulklab/yago/libs/logger"
	"github.com/sirupsen/logrus"
	"strconv"
	"strings"
	"time"
)

type PostFile string
type Body string

type Request struct {
	*httplib.BeegoHTTPRequest
}

// 封装 http 接口的基础类

type HttpThird struct {
	req              *Request
	Domain           string
	Hostname         string
	ConnectTimeout   int
	ReadWriteTimeout int
	headers          map[string]string
	username         string
	password         string
}

func (a *HttpThird) newRequest(method string, api string) *Request {
	var uri string
	if a.Domain != "" {
		uri = strings.TrimRight(a.Domain, "/") + "/" + strings.TrimLeft(api, "/")
	} else {
		uri = api
	}

	switch strings.ToLower(method) {
	case "post":
		a.req = &Request{
			httplib.Post(uri),
		}
	case "get":
		a.req = &Request{
			httplib.Get(uri),
		}
	case "put":
		a.req = &Request{
			httplib.Put(uri),
		}
	case "delete":
		a.req = &Request{
			httplib.Delete(uri),
		}
	case "head":
		a.req = &Request{
			httplib.Head(uri),
		}
	}

	if a.Hostname != "" {
		a.req.SetHost(a.Hostname)
	}

	ctimeout := 3
	wtimeout := 20

	if a.ConnectTimeout != 0 {
		ctimeout = a.ConnectTimeout
	}

	if a.ReadWriteTimeout != 0 {
		wtimeout = a.ReadWriteTimeout
	}

	// 设置超时时间
	a.req.SetTimeout(time.Duration(ctimeout)*time.Second, time.Duration(wtimeout)*time.Second)

	if a.username != "" && a.password != "" {
		a.req.SetBasicAuth(a.username, a.password)
	}

	if len(a.headers) > 0 {
		for k, v := range a.headers {
			a.req.Header(k, v)
		}
	}

	return a.req
}

func (a *HttpThird) SetBaseAuth(username, password string) {
	a.username = username
	a.password = password
}

func (a *HttpThird) SetHeader(headers map[string]string) {
	a.headers = headers
}

func (a *HttpThird) call(method string, api string, params map[string]interface{}) error {
	a.newRequest(method, api)

	logParams := make(map[string]interface{})

	for k, v := range params {
		logParams[k] = v
		switch val := v.(type) {
		case Body: // 原始 body, k 随意
			a.req.Body(v)
		case PostFile: // 文件上传
			a.req.PostFile(k, string(val))
		case string:
			a.req.Param(k, val)
			if len(val) > 1000 {
				logParams[k] = val[:1000] + "..."
			}
		case int64:
			a.req.Param(k, strconv.Itoa(int(val)))
		case int:
			a.req.Param(k, strconv.Itoa(val))
		case uint64:
			a.req.Param(k, strconv.Itoa(int(val)))
		case uint:
			a.req.Param(k, strconv.Itoa(int(val)))
		case float64:
			a.req.Param(k, fmt.Sprintf("%v", val))
		case []byte:
			a.req.Param(k, string(val))
		default:
			return errors.New("unsupported type" + fmt.Sprintf("%T", val))
		}
	}

	begin := time.Now()

	res, err := a.req.Response()

	end := time.Now()
	consume := end.Sub(begin).Nanoseconds() / 1e6

	retStr, _ := a.req.String()

	urlInfo := a.req.GetRequest().URL

	logInfo := logrus.Fields{
		"url":            fmt.Sprintf("%s://%s/%s", urlInfo.Scheme, urlInfo.Host, strings.TrimLeft(urlInfo.Path, "/")),
		"hostname":       a.Hostname,
		"params":         logParams,
		"consume(ms)":    consume,
		"request_header": a.req.GetRequest().Header,
		"result":         retStr,
		"error":          "",
		"category":       "third.http",
	}

	if err != nil {
		logInfo["error"] = err.Error()

		logger.Ins().WithFields(logInfo).Error()

		return errors.New("system err")

	} else if res.StatusCode < 200 || res.StatusCode >= 300 {

		logInfo["error"] = fmt.Sprintf("http status err,code:%d,status:%s", res.StatusCode, res.Status)

		logger.Ins().WithFields(logInfo).Error()

		return errors.New("http status error")
	}

	logger.Ins().WithFields(logInfo).Info()

	return nil
}

func (a *HttpThird) Post(api string, params map[string]interface{}) (*Request, error) {

	//a.Req.Header("Expect", "")

	err := a.call("post", api, params)

	return a.req, err
}

func (a *HttpThird) Get(api string, params map[string]interface{}) (*Request, error) {

	err := a.call("get", api, params)

	return a.req, err
}

func (a *HttpThird) Put(api string, params map[string]interface{}) (*Request, error) {

	err := a.call("put", api, params)

	return a.req, err
}

// @todo 放在 url 上的参数
func (a *HttpThird) Delete(api string, params map[string]interface{}) (*Request, error) {

	err := a.call("delete", api, params)

	return a.req, err
}

func (a *HttpThird) Head(api string, params map[string]interface{}) (*Request, error) {

	err := a.call("head", api, params)

	return a.req, err
}
