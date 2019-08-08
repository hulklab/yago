package yago

import (
	"fmt"
	"strconv"
	"strings"
)

type Err string

const (
	// 1-1000 系统错误, 1000 - 9999 业务公共错误, 10000 - .... 业务错误
	OK           = Err("0")
	E            = Err("1=") // 自定义错误信息
	ErrParam     = Err("2=")
	ErrSign      = Err("3=sign failed")
	ErrAuth      = Err("4=auth failed")
	ErrForbidden = Err("5=forbidden")
	ErrNotLogin  = Err("6=user not login")
	ErrSystem    = Err("7=system error")
	ErrOperate   = Err("8=")
	ErrUnknown   = Err("9=unknown error")
)

func (e Err) Error() string {
	return string(e)
}

func (e Err) GetError() (int, string) {
	if e == OK {
		return 0, ""
	}

	err := strings.SplitN(e.Error(), "=", 2)
	if len(err) != 2 {
		return 1, fmt.Sprintf("Error 格式不正确: %s", e.Error())
	}
	code, _ := strconv.Atoi(err[0])
	return code, err[1]
}

func (e Err) HasErr() bool {
	errCode, _ := e.GetError()
	return errCode != 0
}

// 生成通用错误, 接受通用的 error 类型或者是 string 类型
// eg. yago.NewErr(errors.New("err occur"))
// eg. yago.NewErr("something is error")
// eg. yago.NewErr("%s is err","query")
func NewErr(err interface{}, args ...interface{}) Err {
	if err == nil {
		return OK
	}

	var s string
	switch e := err.(type) {
	case error:
		if e == nil {
			return OK
		} else {
			s = E.Error() + e.Error()
		}
	case string:
		if len(args) > 0 {
			s = E.Error() + fmt.Sprintf(e, args...)
		} else {
			s = E.Error() + e
		}
	}
	return Err(s)
}
