package yago

import (
	"strconv"
	"strings"
)

type Err string

const (
	// 1-1000 系统错误, 1000 - 9999 业务公共错误, 10000 - .... 业务错误
	OK          = Err("0")
	E           = Err("1=") // 自定义错误信息
	ErrParam    = Err("2=参数错误")
	ErrSign     = Err("3=签名认证失败")
	ErrAuth     = Err("4=权限不足")
	ErrNotLogin = Err("5=用户未登录")
	ErrNotFound = Err("6=记录不存在")
	ErrSystem   = Err("7=系统错误")
	ErrOperate  = Err("8=操作失败")
	ErrUnknown  = Err("9=未知错误")
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
		return 1, "Error 格式不正确"
	}
	code, _ := strconv.Atoi(err[0])
	return code, err[1]
}

func (e Err) HasErr() bool {
	errCode, _ := e.GetError()
	return errCode != 0
}

// 生成通用错误, 接受通用的 error 类型或者是 string 类型
func NewErr(err interface{}) Err {
	var s string
	switch e := err.(type) {
	case error:
		if e == nil {
			return OK
		} else {
			s = E.Error() + e.Error()
		}
	case string:
		s = E.Error() + e
	}
	return Err(s)
}
