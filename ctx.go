package yago

import (
	"errors"
	"fmt"
	"mime/multipart"
	"net/http"

	"github.com/gin-gonic/gin"
	validatorv10 "github.com/go-playground/validator/v10"
)

type Ctx struct {
	*gin.Context
}

const (
	ctxYagoKey  = "__YagoCtx__"
	ResponseKey = "__Resp__"
	ErrorKey    = "__Error__"
)

type ResponseBody struct {
	ErrNo  int         `json:"errno"`
	ErrMsg string      `json:"errmsg"`
	Data   interface{} `json:"data,omitempty"`
}

func newCtx(c *gin.Context) *Ctx {
	ctx := &Ctx{
		Context: c,
	}

	c.Set(ctxYagoKey, ctx)
	return ctx
}

func getCtxFromGin(c *gin.Context) (*Ctx, error) {
	v, ok := c.Get(ctxYagoKey)
	if !ok {
		ctx := newCtx(c)
		return ctx, nil
	}

	ctx, ok := v.(*Ctx)
	if !ok {
		return nil, fmt.Errorf("get yago ctx err, yago ctx type error")
	}

	return ctx, nil
}

func (c *Ctx) GetFileContent(key string) ([]byte, error) {
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
	resp := &ResponseBody{
		ErrNo:  ok.Code(),
		ErrMsg: ok.Error(),
		Data:   data,
	}
	c.Set(ResponseKey, resp)

	c.JSON(http.StatusOK, resp)
}

func (c *Ctx) setError(err Err) {
	resp := &ResponseBody{
		ErrNo:  err.Code(),
		ErrMsg: err.Error(),
		Data:   nil,
	}

	c.Set(ResponseKey, resp)

	c.JSON(http.StatusOK, resp)
}

func (c *Ctx) SetError(err interface{}) {
	c.Set(ErrorKey, err)

	switch v := err.(type) {
	case Err:
		c.setError(v)
	case validatorv10.ValidationErrors:
		for _, fieldErr := range v {
			e := ErrParam.String() + fieldErr.Translate(GetTranslator())
			c.setError(Err(e))
			return
		}
	// case json.UnmarshalTypeError:
	case error:
		var ye Err
		e := errors.As(v, &ye)
		if e {
			c.setError(ye)
		} else {
			c.setError(NewErr(v.Error()))
		}
	default:
		c.setError(ErrUnknown)
	}
}

func (c *Ctx) SetDataOrErr(data interface{}, err interface{}) {
	if err == nil {
		c.SetData(data)
		return
	}

	switch v := err.(type) {
	case Err:
		c.SetError(v)
		return
	case error:
		c.SetError(v)
		return
	default:
		c.setError(ErrUnknown)
		return
	}
}

// Abort in the middleware with yago error or error
func (c *Ctx) AbortWithE(err error) {
	c.Abort()
	c.SetError(err)
}

func (c *Ctx) GetError() error {
	err, exist := c.Get(ErrorKey)
	if !exist {
		return nil
	}
	return err.(error)
}

func (c *Ctx) GetResponse() (*ResponseBody, bool) {
	resp, exist := c.Get(ResponseKey)

	if !exist {
		return nil, false
	}

	if v, ok := resp.(*ResponseBody); !ok {
		return nil, false
	} else {
		return v, true
	}
}

func (c *Ctx) Copy() *Ctx {
	c.Context = c.Context.Copy()
	return c
}
