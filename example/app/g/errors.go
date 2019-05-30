package g

import "github.com/hulklab/yago"

const (
	// 1-9999 系统错误, 10000 - .... 业务错误
	ErrDemo = yago.Err("10001=错误样例")
)
