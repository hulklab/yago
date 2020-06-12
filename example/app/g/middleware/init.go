package middleware

import "github.com/hulklab/yago"

func init() {
	yago.HttpUse(BizLog)
}
