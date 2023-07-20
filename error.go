package yago

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type Err string

var (
	// You can replace yago Err variable in your application,eg. yago.ErrParam = Err("4000=Param err")
	ok           = Err("0")
	E            = Err("1=") // custom error
	ErrParam     = Err("2=")
	ErrSign      = Err("3=Sign failed")
	ErrAuth      = Err("4=Auth failed")
	ErrForbidden = Err("5=Forbidden")
	ErrNotLogin  = Err("6=User not login")
	ErrSystem    = Err("7=System error")
	ErrOperate   = Err("8=")
	ErrUnknown   = Err("9=Unknown error")
)

func (e Err) Error() string {
	_, err := e.getError()
	return err
}

func (e Err) String() string {
	return string(e)
}

func (e Err) Code() int {
	code, _ := e.getError()
	return code
}

func (e Err) getError() (int, string) {
	if e == ok || e == "" {
		return 0, ""
	}

	err := strings.SplitN(e.String(), "=", 2)
	if len(err) != 2 {
		return 1, fmt.Sprintf("Error 格式不正确: %s", e.String())
	}
	code, _ := strconv.Atoi(err[0])
	return code, err[1]
}

// Deprecated
func (e Err) HasErr() bool {
	return e.Code() != 0
}

// 生成通用错误, 接受通用的 error 类型或者是 string 类型
// eg. yago.NewErr(404, "record is not exists")
// eg. yago.NewErr(yago.Forbidden, "you are not permitted")
// eg. yago.NewErr(errors.New("err occur"))
// eg. yago.NewErr("something is error")
// eg. yago.NewErr("%s is err","query")
func NewErr(err interface{}, args ...interface{}) Err {
	if err == nil {
		return ok
	}

	var s string
	switch e := err.(type) {
	case Err:
		if len(args) > 0 {
			if len(e.Error()) == 0 {
				s = fmt.Sprintf("%s%s", e.String(), args[0])
			} else {
				s = fmt.Sprintf("%s: %s", e.String(), args[0])
			}
		} else {
			return e
		}
	case int:
		if e == 0 {
			return ok
		}

		if len(args) > 0 {
			s = fmt.Sprintf("%d=%s",e,args[0])
		} else {
			s = fmt.Sprintf("%d=",e)
		}
	case error:
		if e == nil {
			return ok
		} else {
			s = E.String() + e.Error()
		}
	case string:
		if len(args) > 0 {
			s = E.String() + fmt.Sprintf(e, args...)
		} else {
			s = E.String() + e
		}
	}
	return Err(s)
}

// 返回业务报错（业务报错给接口展示），包裹系统错误（系统错误转到日志记录）
func WrapErr(ye Err, err error) error {
	if err == nil {
		return NewErr("err can not be nil when use yago.WrapErr()")
	}
	if ye == ok {
		return NewErr("ye can not be OK when use yago.WrapErr()")
	}
	return fmt.Errorf("%w: %s", ye, err)
}

type HTTPCodeError int

func (e HTTPCodeError) Code() int {
	return int(e)
}

func (e HTTPCodeError) Error() string {
	code := e.Code()
	return fmt.Sprintf("%d: %s", code, http.StatusText(code))
}

// check if err is http code error
func AsHTTPCodeError(err error) (b bool,e HTTPCodeError) {
	b = errors.As(err, &e)
	return
}
