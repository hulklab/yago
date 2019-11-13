package basethird

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/hulklab/yago/coms/logger"
	"github.com/levigross/grequests"
	"github.com/sirupsen/logrus"
)

type PostFile string

func (p PostFile) Value(name string) (f grequests.FileUpload, e error) {
	fd, err := os.Open(string(p))

	if err != nil {
		return f, err
	}

	return grequests.FileUpload{FileContents: fd, FileName: string(p), FieldName: name}, nil
}

type Body string

func (b Body) Value() io.Reader {
	bf := bytes.NewBufferString(string(b))
	return ioutil.NopCloser(bf)
}

type Response struct {
	*grequests.Response
}

// rewrite ToJSON then you can use ToJSON many times
func (r *Response) ToJSON(v interface{}) error {
	if r.Error != nil {
		return r.Error
	}

	var reader io.Reader

	reader = bytes.NewBuffer(r.Bytes())

	jsonDecoder := json.NewDecoder(reader)

	defer r.Close()

	return jsonDecoder.Decode(&v)
}

func (r *Response) JSON(v interface{}) error {
	return r.ToJSON(v)
}

func (r *Response) String() (string, error) {
	return r.Response.String(), r.Error
}

// 封装 http 接口的基础类
type HttpThird struct {
	client           *http.Client
	Domain           string
	Hostname         string
	ConnectTimeout   int
	ReadWriteTimeout int
	headers          map[string]string
	username         string
	password         string
	tlsCfg           *tls.Config
	logInfoOff       bool
	once             sync.Once
}

func (a *HttpThird) getConnTimeout() time.Duration {
	ctimeout := 3
	if a.ConnectTimeout != 0 {
		ctimeout = a.ConnectTimeout
	}
	return time.Duration(ctimeout) * time.Second
}

func (a *HttpThird) getRequestTimeout() time.Duration {
	wtimeout := 20
	if a.ReadWriteTimeout != 0 {
		wtimeout = a.ReadWriteTimeout
	}
	return time.Duration(wtimeout) * time.Second
}

func (a *HttpThird) getClient() *http.Client {
	// @todo proxy
	a.once.Do(func() {
		a.client = &http.Client{
			Transport: &http.Transport{
				MaxIdleConnsPerHost: 100,
				TLSClientConfig:     a.tlsCfg,
				DialContext: (&net.Dialer{
					Timeout: a.getConnTimeout(),
					//KeepAlive: 30* time.Second,
				}).DialContext,
				//Proxy:               Proxy,
			},
			Timeout: a.getRequestTimeout(),
		}
	})

	return a.client
}

func (a *HttpThird) newRo() *grequests.RequestOptions {

	ro := &grequests.RequestOptions{
		HTTPClient: a.getClient(),
		UserAgent:  "github.com-hulklab-yago",
	}

	if a.Hostname != "" {
		ro.Host = a.Hostname
	}

	if a.username != "" && a.password != "" {
		ro.Auth = []string{a.username, a.password}
	}

	if len(a.headers) > 0 {
		ro.Headers = a.headers
	}

	return ro
}

func (a *HttpThird) genUri(api string) string {
	var uri string
	if a.Domain != "" {
		uri = strings.TrimRight(a.Domain, "/") + "/" + strings.TrimLeft(api, "/")
	} else {
		uri = api
	}

	return uri
}

func (a *HttpThird) SetBaseAuth(username, password string) {
	a.username = username
	a.password = password
}

func (a *HttpThird) SetHeader(headers map[string]string) {
	a.headers = headers
}

func (a *HttpThird) SetTLSClientConfig(cfg *tls.Config) {
	a.tlsCfg = cfg
}

// 设置是否要关闭 info 日志
func (a *HttpThird) SetLogInfoFlag(on bool) {
	if on {
		a.logInfoOff = false
	} else {
		a.logInfoOff = true
	}
}

func (a *HttpThird) call(method string, api string, params map[string]interface{}) (*Response, error) {
	log := logger.Ins().Category("third.http")

	//a.newRequest(method, api)
	ro := a.newRo()
	logParams := make(map[string]interface{})
	dataParams := make(map[string]string)

	for k, v := range params {
		logParams[k] = v
		switch val := v.(type) {
		case Body: // 原始 body, k 随意
			ro.RequestBody = val.Value()
			//a.req.Body(v)
		case PostFile: // 文件上传
			uf, err := val.Value(k)
			if err != nil {
				log.Error("post file params err:", err)
				continue
			}
			ro.Files = append(ro.Files, uf)
		case string:
			dataParams[k] = val
			if len(val) > 1000 {
				logParams[k] = val[:1000] + "..."
			}
		case int64:
			dataParams[k] = strconv.Itoa(int(val))
		case int:
			dataParams[k] = strconv.Itoa(val)
		case uint64:
			dataParams[k] = strconv.Itoa(int(val))
		case uint:
			dataParams[k] = strconv.Itoa(int(val))
		case float64:
			dataParams[k] = fmt.Sprintf("%v", val)
		case []byte:
			dataParams[k] = string(val)
		default:
			return nil, errors.New("unsupported type" + fmt.Sprintf("%T", val))
		}
	}

	if len(dataParams) > 0 {
		if strings.ToUpper(method) == "GET" {
			ro.Params = dataParams
		} else {
			ro.Data = dataParams
		}
	}

	//fmt.Printf("%+v", ro)

	uri := a.genUri(api)

	begin := time.Now()

	res, err := grequests.Req(method, uri, ro)

	end := time.Now()
	consume := end.Sub(begin).Nanoseconds() / 1e6

	retStr := res.String()

	logInfo := logrus.Fields{
		"url":            uri,
		"hostname":       a.Hostname,
		"params":         logParams,
		"consume":        consume,
		"request_header": ro.Headers,
	}

	if err != nil {
		logInfo["hint"] = err.Error()

		log.WithFields(logInfo).Error()

		return nil, err

	} else if !res.Ok {

		logInfo["hint"] = fmt.Sprintf("http status err,code:%d", res.StatusCode)

		log.WithFields(logInfo).Error()

		return nil, errors.New("http status error")
	}

	// 默认是日志没关
	if !a.logInfoOff {
		logInfo["result"] = retStr
	}
	log.WithFields(logInfo).Info()

	return &Response{res}, nil
}

func (a *HttpThird) Post(api string, params map[string]interface{}) (*Response, error) {

	//a.Req.Header("Expect", "")

	return a.call("POST", api, params)
}

func (a *HttpThird) Get(api string, params map[string]interface{}) (*Response, error) {

	return a.call("GET", api, params)

}

func (a *HttpThird) Put(api string, params map[string]interface{}) (*Response, error) {

	return a.call("PUT", api, params)
}

func (a *HttpThird) Delete(api string, params map[string]interface{}) (*Response, error) {

	return a.call("DELETE", api, params)
}

func (a *HttpThird) Head(api string, params map[string]interface{}) (*Response, error) {

	return a.call("HEAD", api, params)
}
