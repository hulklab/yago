package basethird

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
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
	client              *http.Client
	Domain              string
	Hostname            string
	ConnectTimeout      int
	ReadWriteTimeout    int
	headers             map[string]string
	username            string
	password            string
	tlsCfg              *tls.Config
	logInfoOff          bool
	once                sync.Once
	maxIdleConnsPerHost int
	maxConnsPerHost     int
}

func (a *HttpThird) SetMaxIdleConnsPerHost(num int) {
	a.maxIdleConnsPerHost = num
}

func (a *HttpThird) getMaxIdleConnsPerHost() int {
	if a.maxIdleConnsPerHost == 0 {
		return 20
	}
	return a.maxIdleConnsPerHost
}

func (a *HttpThird) SetMaxConnsPerHost(num int) {
	a.maxConnsPerHost = num
}

func (a *HttpThird) getMaxConnsPerHost() int {
	if a.maxConnsPerHost == 0 {
		return 500
	}
	return a.maxConnsPerHost
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
	a.once.Do(func() {
		a.client = &http.Client{
			Transport: &http.Transport{
				Proxy:               http.ProxyFromEnvironment,
				MaxIdleConnsPerHost: a.getMaxIdleConnsPerHost(),
				MaxConnsPerHost:     a.getMaxConnsPerHost(),
				IdleConnTimeout:     90 * time.Second,
				TLSClientConfig:     a.tlsCfg,
				DialContext: (&net.Dialer{
					Timeout:   a.getConnTimeout(),
					KeepAlive: 30 * time.Second,
				}).DialContext,
				ExpectContinueTimeout: 1 * time.Second,
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
		// 如果请求的 api 里面带有 http 协议，则保留不拼接
		u, err := url.Parse(api)
		if err == nil && (u.Scheme == "http" || u.Scheme == "https") {
			uri = api
		} else {
			uri = strings.TrimRight(a.Domain, "/") + "/" + strings.TrimLeft(api, "/")
		}
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

func mergeOptions(ro *grequests.RequestOptions, opts ...*grequests.RequestOptions) {
	if len(opts) == 0 {
		return
	}
	opt := opts[0]

	if opt.Context != nil {
		ro.Context = opt.Context
	}

	if opt.Cookies != nil {
		ro.Cookies = opt.Cookies
	}

	if opt.Data != nil {
		ro.Data = opt.Data
	}

	if opt.Auth != nil {
		ro.Auth = opt.Auth
	}

	if opt.Headers != nil {
		ro.Headers = opt.Headers
	}

	if opt.JSON != nil {
		ro.JSON = opt.JSON
	}

	if opt.CookieJar != nil {
		ro.CookieJar = opt.CookieJar
	}

	if opt.UserAgent != "" {
		ro.UserAgent = opt.UserAgent
	}

	if opt.BeforeRequest != nil {
		ro.BeforeRequest = opt.BeforeRequest
	}
}

func (a *HttpThird) call(method string, api string, params map[string]interface{}, opts ...*grequests.RequestOptions) (*Response, error) {
	log := logger.Ins().Category("third.http")

	ro := a.newRo()

	logParams := make(map[string]interface{})
	dataParams := make(map[string]string)

	for k, v := range params {
		logParams[k] = v
		switch val := v.(type) {
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
		if method == http.MethodGet {
			ro.Params = dataParams
		} else {
			ro.Data = dataParams
		}
	}

	mergeOptions(ro, opts...)

	//fmt.Printf("%+v", ro)

	uri := a.genUri(api)

	begin := time.Now()

	res, err := grequests.Req(method, uri, ro)

	end := time.Now()
	consume := end.Sub(begin).Nanoseconds() / 1e6

	logInfo := logrus.Fields{
		"url":            uri,
		"hostname":       a.Hostname,
		"params":         logParams,
		"consume":        consume,
		"request_header": ro.Headers,
	}

	if ro.JSON != nil {
		var requestBody string
		switch ro.JSON.(type) {
		case string:
			requestBody = ro.JSON.(string)
		case []byte:
			requestBody = string(ro.JSON.([]byte))
		default:
			byteSlice, err := json.Marshal(ro.JSON)
			if err == nil {
				requestBody = string(byteSlice)
			}
		}
		logInfo["request_body"] = requestBody
	}

	if err != nil {
		logInfo["hint"] = err.Error()

		log.WithFields(logInfo).Error()

		return nil, err
	}

	retStr := res.String()

	// 默认是日志没关
	if !a.logInfoOff {
		logInfo["result"] = retStr
	}

	if !res.Ok {
		logInfo["hint"] = fmt.Sprintf("http status err,code:%d", res.StatusCode)

		log.WithFields(logInfo).Error()

		return nil, fmt.Errorf("http status error: %d", res.StatusCode)
	}

	log.WithFields(logInfo).Info()

	return &Response{res}, nil
}

func (a *HttpThird) Post(api string, params map[string]interface{}, opts ...*grequests.RequestOptions) (*Response, error) {
	//a.Req.Header("Expect", "")
	return a.call(http.MethodPost, api, params, opts...)
}

func (a *HttpThird) Get(api string, params map[string]interface{}, opts ...*grequests.RequestOptions) (*Response, error) {
	return a.call(http.MethodGet, api, params, opts...)
}

func (a *HttpThird) Put(api string, params map[string]interface{}, opts ...*grequests.RequestOptions) (*Response, error) {
	return a.call(http.MethodPut, api, params, opts...)
}

func (a *HttpThird) Patch(api string, params map[string]interface{}, opts ...*grequests.RequestOptions) (*Response, error) {
	return a.call(http.MethodPatch, api, params, opts...)
}

func (a *HttpThird) Delete(api string, params map[string]interface{}, opts ...*grequests.RequestOptions) (*Response, error) {
	return a.call(http.MethodDelete, api, params, opts...)
}

func (a *HttpThird) Head(api string, opts ...*grequests.RequestOptions) (*Response, error) {
	return a.call(http.MethodHead, api, nil, opts...)
}

func (a *HttpThird) Options(api string, opts ...*grequests.RequestOptions) (*Response, error) {
	return a.call(http.MethodOptions, api, nil, opts...)
}
