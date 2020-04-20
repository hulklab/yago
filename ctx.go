package yago

import (
	"errors"
	"mime/multipart"
	"net/http"

	"github.com/gin-gonic/gin"
	validatorv10 "github.com/go-playground/validator/v10"
)

type Ctx struct {
	*gin.Context
	resp *ResponseBody
	err  error
}

const CtxParamsKey = "__PARAMS__"

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
	c.resp = &ResponseBody{
		ErrNo:  ok.Code(),
		ErrMsg: ok.Error(),
		Data:   data,
	}

	c.JSON(http.StatusOK, c.resp)
}

func (c *Ctx) setError(err Err) {
	c.resp = &ResponseBody{
		ErrNo:  err.Code(),
		ErrMsg: err.Error(),
		Data:   map[string]interface{}{},
	}

	c.JSON(http.StatusOK, c.resp)
}

func (c *Ctx) GetError() error {
	return c.err
}

func (c *Ctx) SetError(err interface{}) {

	switch v := err.(type) {
	case Err:
		c.err = v
		c.setError(v)
	case validatorv10.ValidationErrors:
		for _, fieldErr := range v {
			e := ErrParam.String() + fieldErr.Translate(GetTranslator())
			c.err = Err(e)
			c.setError(Err(e))
			return
		}
	//case json.UnmarshalTypeError:
	case error:
		e := errors.Unwrap(v.(error))
		c.err = v
		if er, ok := e.(Err); ok {
			c.setError(er)
		} else {
			c.setError(NewErr(v.Error()))
		}

	default:
		c.err = ErrUnknown
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
		if v.HasErr() {
			c.SetError(v)
			return
		}
	case error:
		if v != nil {
			c.SetError(v)
			return
		}
	default:
		c.setError(ErrUnknown)
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

func (c *Ctx) Copy() *Ctx {
	c.Context = c.Context.Copy()
	return c
}
