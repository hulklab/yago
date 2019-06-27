package yago

import (
	"github.com/gin-gonic/gin"
	"github.com/hulklab/yago/libs/validator"
	"net/http"
	"strconv"
	"strings"
)

type Ctx struct {
	*gin.Context
}

type ResponseBody struct {
	ErrNo  int         `json:"errno"`
	ErrMsg string      `json:"errmsg"`
	Data   interface{} `json:"data"`
}

func NewCtx(c *gin.Context) (*Ctx, error) {
	ctx := &Ctx{c}

	if err := ctx.Validate(); err != nil {
		return ctx, err
	}

	return ctx, nil
}

func (c *Ctx) Validate() error {
	url := c.Request.URL.Path
	paths := strings.Split(url, "/")
	action := paths[len(paths)-1]
	if router, ok := HttpRouterMap[url]; ok {
		rules := router.h.Rules()
		labels := router.h.Labels()
		check, err := validator.ValidateHttp(c.Context, action, labels, rules)
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
	for _, v := range slice {
		vInt, _ := strconv.Atoi(v)
		sliceInt = append(sliceInt, vInt)
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
	for _, v := range slice {
		vInt64, _ := strconv.ParseInt(v, 10, 64)
		sliceInt64 = append(sliceInt64, vInt64)
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

func (c *Ctx) SetData(data interface{}) {
	errCode, errMsg := OK.GetError()
	c.JSON(http.StatusOK, ResponseBody{
		ErrNo:  errCode,
		ErrMsg: errMsg,
		Data:   data,
	})
}

func (c *Ctx) SetError(err Err, msgEx ...string) {
	errCode, errMsg := err.GetError()
	if len(msgEx) > 0 {
		if errMsg == "" {
			errMsg = msgEx[0]
		} else {
			errMsg = errMsg + ": " + msgEx[0]
		}
	}
	c.JSON(http.StatusOK, ResponseBody{
		ErrNo:  errCode,
		ErrMsg: errMsg,
		Data:   nil,
	})
}

func (c *Ctx) SetDataOrErr(data interface{}, err Err) {

	if err.HasErr() {
		c.SetError(err)
		return
	}

	c.SetData(data)
}
