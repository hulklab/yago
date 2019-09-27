package yago

import (
	"fmt"
	"log"
	"github.com/gin-gonic/gin"
	"github.com/hulklab/yago/libs/validator"
	"mime/multipart"
	"reflect"
	"net/http"
	"strconv"
	"strings"
)

type Ctx struct {
	*gin.Context
	resp *ResponseBody
}

type ResponseBody struct {
	ErrNo  int         `json:"errno"`
	ErrMsg string      `json:"errmsg"`
	Data   interface{} `json:"data"`
}

func NewCtx(c *gin.Context) *Ctx {
	return &Ctx{
		Context: c,
	}
}

func (c *Ctx) Validate() error {
	url := c.Request.URL.Path
	paths := strings.Split(url, "/")
	action := paths[len(paths)-1]
	if router, ok := HttpRouterMap[url]; ok {
		rules := router.h.Rules()
		labels := router.h.Labels()
		check, err := ValidateHttp(c, action, labels, rules)
		if !check {
			return err
		}
	}
	return nil
}

func (c *Ctx) RequestString(key string, def ...string) string {
	val := c.GetString(key)
	if len(def) > 0 && val == "" {
		val = def[0]
	}
	return val
}

func (c *Ctx) RequestSliceString(key string, def ...string) []string {
	val := c.GetString(key)
	if val == "" {
		if len(def) > 0 {
			return def
		}
		return nil
	}
	return strings.Split(val, ",")
}

func (c *Ctx) RequestInt(key string, def ...int) (int, error) {
	val := c.GetString(key)
	if len(def) > 0 && val == "" {
		return def[0], nil
	}
	return strconv.Atoi(val)
}

func (c *Ctx) RequestInt64(key string, def ...int64) (int64, error) {
	val := c.GetString(key)
	if len(def) > 0 && val == "" {
		return def[0], nil
	}
	return strconv.ParseInt(val, 10, 64)
}

func (c *Ctx) RequestFloat64(key string, def ...float64) (float64, error) {
	val := c.GetString(key)
	if len(def) > 0 && val == "" {
		return def[0], nil
	}
	return strconv.ParseFloat(val, 64)
}

func (c *Ctx) RequestSliceInt(key string, def ...int) []int {
	val := c.GetString(key)
	if val == "" {
		if len(def) > 0 {
			return def
		}
		return nil
	}
	slice := strings.Split(val, ",")
	sliceInt := make([]int, len(slice))
	for k, v := range slice {
		vInt, _ := strconv.Atoi(v)
		sliceInt[k] = vInt
	}
	return sliceInt
}

func (c *Ctx) RequestSliceInt64(key string, def ...int64) []int64 {
	val := c.GetString(key)
	if val == "" {
		if len(def) > 0 {
			return def
		}
		return nil
	}
	slice := strings.Split(val, ",")
	sliceInt64 := make([]int64, len(slice))
	for k, v := range slice {
		vInt64, _ := strconv.ParseInt(v, 10, 64)
		sliceInt64[k] = vInt64
	}
	return sliceInt64
}

func (c *Ctx) RequestBool(key string, def ...bool) bool {
	val := c.GetString(key)
	if len(def) > 0 && val == "" {
		return def[0]
	}
	switch val {
	case "1", "true", "TRUE", "True":
		return true
	case "0", "false", "FALSE", "False":
		return false
	}
	return false
}

func (c *Ctx) RequestFileContent(key string) ([]byte, error) {
	formFile, err := c.FormFile(key)
	if err != nil {
		return nil, err
	}
	var file multipart.File
	file, err = formFile.Open()
	if err != nil {
		return nil, err
	}
	content := make([]byte, formFile.Size)
	_, err = file.Read(content)
	if err != nil {
		return nil, err
	}
	return content, nil
}

func (c *Ctx) SetData(data interface{}) {
	c.resp = &ResponseBody{
		ErrNo:  OK.Code(),
		ErrMsg: OK.Error(),
		Data:   data,
	}

	c.JSON(http.StatusOK, c.resp)
}

func (c *Ctx) SetError(err Err, msgEx ...string) {
	errMsg := err.Error()
	if len(msgEx) > 0 {
		if errMsg == "" {
			errMsg = msgEx[0]
		} else {
			errMsg = errMsg + ": " + msgEx[0]
		}
	}

	c.resp = &ResponseBody{
		ErrNo:  err.Code(),
		ErrMsg: errMsg,
		Data:   map[string]interface{}{},
	}

	c.JSON(http.StatusOK, c.resp)
}

func (c *Ctx) SetDataOrErr(data interface{}, err Err) {
	if err.HasErr() {
		c.SetError(err)
		return
	}

	c.SetData(data)
}

func (c *Ctx) GetResponse() (*ResponseBody, bool) {
	if c.resp == nil {
		return c.resp, false
	}

	return c.resp, true
}

func ValidateHttp(c *Ctx, action string, labels validator.Label, rules []validator.Rule) (bool, error) {
	type CustomValidatorFunc = func(c *Ctx, p string) (valid bool, err error)

	for _, rule := range rules {
		actionMatch := false
		if len(rule.On) == 0 {
			actionMatch = true
		} else {
			for _, a := range rule.On {
				if a == action {
					actionMatch = true
					break
				}
			}
		}

		if actionMatch {
			switch method := rule.Method.(type) {
			case int:
				return validateByRule(c, labels, rule, method)
			case CustomValidatorFunc:
				for _, p := range rule.Params {
					_, exist := c.Get(p)
					if !exist {
						continue
					}
					if valid, err := method(c, p); !valid {
						return false, err
					}
				}
			default:
				log.Fatalf("not support method: %s", reflect.TypeOf(rule.Method))
			}
		}
	}
	return true, nil
}

func validateByRule(c *Ctx, labels validator.Label, rule validator.Rule, method int) (bool, error) {
	switch method {
	case validator.Required:
		for _, p := range rule.Params {
			pv, exist := c.Get(p)
			if !exist {
				return false, fmt.Errorf("%s 不存在", labels.Get(p))
			}
			if valid, err := (validator.RequiredValidator{}).Check(pv); !valid {
				return false, getErr(labels.Get(p), err, rule.Message)
			}
		}
	case validator.String:
		for _, p := range rule.Params {
			pv, exist := c.Get(p)
			if !exist {
				continue
			}
			if valid, err := (validator.StringValidator{Min: int(rule.Min), Max: int(rule.Max)}).Check(pv); !valid {
				return false, getErr(labels.Get(p), err, rule.Message)
			}
		}
	case validator.Int:
		for _, p := range rule.Params {
			pv, exist := c.Get(p)
			if !exist {
				continue
			}
			pvInt, err := strconv.Atoi(pv.(string))
			if err != nil {
				return false, fmt.Errorf("%s 不是个整数", labels.Get(p))
			}
			if valid, err := (validator.IntValidator{Min: int(rule.Min), Max: int(rule.Max)}).Check(pvInt); !valid {
				return false, getErr(labels.Get(p), err, rule.Message)
			}
		}
	case validator.Float:
		for _, p := range rule.Params {
			pv, exist := c.Get(p)
			if !exist {
				continue
			}

			pvFloat, err := strconv.ParseFloat(pv.(string), 64)
			if err != nil {
				return false, fmt.Errorf("%s 不是个浮点数", labels.Get(p))
			}
			if valid, err := (validator.FloatValidator{Min: rule.Min, Max: rule.Max}).Check(pvFloat); !valid {
				return false, getErr(labels.Get(p), err, rule.Message)
			}
		}
	case validator.JSON:
		for _, p := range rule.Params {
			pv, exist := c.Get(p)
			if !exist {
				continue
			}
			if valid, err := (validator.JSONValidator{}).Check(pv); !valid {
				return false, getErr(labels.Get(p), err, rule.Message)
			}
		}
	case validator.IP:
		for _, p := range rule.Params {
			pv, exist := c.Get(p)
			if !exist {
				continue
			}
			if valid, err := (validator.IPValidator{}).Check(pv); !valid {
				return false, getErr(labels.Get(p), err, rule.Message)
			}
		}
	case validator.Match:
		for _, p := range rule.Params {
			pv, exist := c.Get(p)
			if !exist {
				continue
			}
			if valid, err := (validator.MatchValidator{Pattern: rule.Pattern}).Check(pv); !valid {
				return false, getErr(labels.Get(p), err, rule.Message)
			}
		}
	}
	return true, nil
}

func getErr(label string, err error, message string) error {
	if message == "" {
		return fmt.Errorf("%s %s", label, err)
	}
	return fmt.Errorf("%s %s", label, message)
}
