package basethird

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/hulklab/yago"

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

func ErrResponse(err error) *Response {
	res := &grequests.Response{Error: err}

	return &Response{res}
}

type Caller func(method, uri string, ro *grequests.RequestOptions) (resp *Response, err error)

type HttpInterceptor func(method, uri string, ro *grequests.RequestOptions, call Caller) (*Response, error)

// 封装 http 接口的基础类
type HttpThird struct {
	client                    *http.Client
	Address                   string
	Hostname                  string
	ConnectTimeout            int
	ReadWriteTimeout          int
	SslOn                     bool
	CertFile                  string
	headers                   map[string]string
	username                  string
	password                  string
	tlsCfg                    *tls.Config
	logInfoOff                bool
	once                      sync.Once
	maxIdleConnsPerHost       int
	maxConnsPerHost           int
	interceptors              []HttpInterceptor
	onceInter                 sync.Once
	disableDefaultInterceptor bool
}

func (a *HttpThird) InitConfig(configSection string) error {
	if !yago.Config.IsSet(configSection) {
		return fmt.Errorf("config section %s is not exists", configSection)
	}

	a.Address = yago.Config.GetString(configSection + ".address")
	a.Hostname = yago.Config.GetString(configSection + ".hostname")
	a.ReadWriteTimeout = yago.Config.GetInt(configSection + ".timeout")
	a.ConnectTimeout = yago.Config.GetInt(configSection + ".conn_timeout")
	a.SslOn = yago.Config.GetBool(configSection + ".ssl_on")
	a.CertFile = yago.Config.GetString(configSection + ".cert_file")
	if a.SslOn {
		if a.CertFile == "" {
			return fmt.Errorf("cert file is required in config section %s", configSection)
		}

		roots := x509.NewCertPool()

		pem, err := ioutil.ReadFile(a.CertFile)
		if err != nil {
			return fmt.Errorf("read cert file %s error", a.CertFile)
		}

		isOk := roots.AppendCertsFromPEM(pem)
		if !isOk {
			return fmt.Errorf("load cert file %s error", a.CertFile)
		}

		sslConf := &tls.Config{
			RootCAs: roots,
		}

		a.SetTLSClientConfig(sslConf)
	}

	return nil
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
	if a.Address != "" {
		// 如果请求的 api 里面带有 http 协议，则保留不拼接
		u, err := url.Parse(api)
		if err == nil && (u.Scheme == "http" || u.Scheme == "https") {
			uri = api
		} else {
			uri = strings.TrimRight(a.Address, "/") + "/" + strings.TrimLeft(api, "/")
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

func (a *HttpThird) SetTLSInsecure() {
	a.tlsCfg = &tls.Config{
		InsecureSkipVerify: true,
	}
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

	if opt.RequestBody != nil {
		ro.RequestBody = opt.RequestBody
	}

	if opt.RedirectLimit != 0 {
		ro.RedirectLimit = opt.RedirectLimit
	}

	if opt.Proxies != nil {
		ro.Proxies = opt.Proxies
	}

	if opt.LocalAddr != nil {
		ro.LocalAddr = opt.LocalAddr
	}

	if opt.BeforeRequest != nil {
		ro.BeforeRequest = opt.BeforeRequest
	}

	if opt.Files != nil {
		ro.Files = opt.Files
	}
}

func (a *HttpThird) call(method string, api string, params map[string]interface{}, opts ...*grequests.RequestOptions) (*Response, error) {
	log := logger.Ins().Category("third.http")

	ro := a.newRo()

	dataParams := make(map[string]string)

	for k, v := range params {
		//logParams[k] = v
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
		case bool:
			dataParams[k] = strconv.FormatBool(val)
		default:
			err := errors.New("unsupported type" + fmt.Sprintf("%T", val))
			return ErrResponse(err), err
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

	uri := a.genUri(api)

	chainCaller := a.chainCaller()

	return chainCaller(method, uri, ro)
}

func (a *HttpThird) AddInterceptor(hi HttpInterceptor) {
	if a.interceptors == nil {
		a.interceptors = make([]HttpInterceptor, 0)
	}

	a.interceptors = append(a.interceptors, hi)
}

func (a *HttpThird) getInterceptors() []HttpInterceptor {

	a.onceInter.Do(func() {
		// 注册日志插件(放到最后)
		if !a.disableDefaultInterceptor {
			a.AddInterceptor(a.logInterceptor)
		}
	})

	return a.interceptors
}

// build caller chain
// TODO cache?
func (a *HttpThird) chainCaller() Caller {
	innerCaller := func(method, uri string, ro *grequests.RequestOptions) (*Response, error) {
		res, err := grequests.Req(method, uri, ro)
		return &Response{res}, err
	}

	chainWrap := func(currentInter HttpInterceptor, currentCaller Caller) Caller {
		return func(method, uri string, ro *grequests.RequestOptions) (*Response, error) {
			return currentInter(method, uri, ro, currentCaller)
		}
	}

	chainedCaller := innerCaller

	interceptors := a.getInterceptors()
	n := len(interceptors)

	if n >= 1 {
		for i := n - 1; i >= 0; i-- {
			chainedCaller = chainWrap(interceptors[i], chainedCaller)
		}
	}
	return chainedCaller
}

func (a *HttpThird) DisableDefaultInterceptor() {
	a.disableDefaultInterceptor = true
}

func (a *HttpThird) logInterceptor(method, uri string, ro *grequests.RequestOptions, call Caller) (*Response, error) {
	log := logger.Ins().Category("third.http")

	var dataParams map[string]string
	logParams := make(map[string]string)

	if method == http.MethodGet {
		dataParams = ro.Params
	} else {
		dataParams = ro.Data
	}

	if len(dataParams) > 0 {
		for k, val := range dataParams {
			if len(val) > 1000 {
				logParams[k] = val[:1000] + "..."
			} else {
				logParams[k] = val
			}
		}
	}

	//log.Printf("before invoker. method: %+v, request:%+v", method, req)
	begin := time.Now()

	resp, err := call(method, uri, ro)

	end := time.Now()
	consume := end.Sub(begin).Nanoseconds() / 1e6
	logInfo := logrus.Fields{
		"url":            uri,
		"hostname":       a.Hostname,
		"method":         method,
		"params":         logParams,
		"consume":        consume,
		"request_header": ro.Headers,
		//"response_header": resp.Header,
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

		return resp, err
	}

	if resp == nil {
		log.WithFields(logInfo).Error()
		return resp, err
	}

	logInfo["response_header"] = resp.Header
	logInfo["status_code"] = resp.StatusCode

	retStr, _ := resp.String()

	// 默认是日志没关
	if !a.logInfoOff {
		logInfo["result"] = retStr
	}

	if !resp.Ok {
		logInfo["hint"] = fmt.Sprintf("http status err,code:%d,body:%s", resp.StatusCode, retStr)

		log.WithFields(logInfo).Error()

		return resp, yago.HTTPCodeError(resp.StatusCode)
	}

	log.WithFields(logInfo).Info()

	return resp, nil
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
