package yago

import (
	"mime/multipart"
	"net/http"

	"github.com/gin-gonic/gin"
	validatorv10 "github.com/go-playground/validator/v10"
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
		ErrNo:  OK.Code(),
		ErrMsg: OK.Error(),
		Data:   data,
	}

	c.JSON(http.StatusOK, c.resp)
}

func (c *Ctx) setError(err Err, msgEx ...string) {
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

func (c *Ctx) SetError(err interface{}, msgEx ...string) {

	switch v := err.(type) {
	case Err:
		c.setError(v, msgEx...)
	case validatorv10.ValidationErrors:
		for _, fieldErr := range v {
			e := ErrParam.String() + fieldErr.Translate(GetTranslator())
			c.setError(Err(e))
			return
		}
	//case json.UnmarshalTypeError:
	case error:
		c.setError(NewErr(v.Error()))
	default:
		c.setError(ErrUnknown)
	}
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
