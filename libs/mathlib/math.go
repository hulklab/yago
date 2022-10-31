package mathlib

import (
	"fmt"
	"math"
	"strings"
)

/**
 * 四舍五入
 * @param f float64
 * @param n int 取小数点后几位
 * @return float64
 *
 */
func Round(f float64, n int) float64 {
	pow10N := math.Pow10(n)
	return math.Trunc((f+0.5/pow10N)*pow10N) / pow10N
}

// 向下取整, 并保留 n 位小数
func Floor(f float64, n int) float64 {
	if n == 0 {
		return math.Floor(f)
	}

	a := fmt.Sprintf("%g", f)

	ss := strings.Split(a, ".")
	if len(ss) < 2 {
		return math.Floor(f)
	}

	// 如果小数点后位数等于截取长度，直接返回
	if len(ss[1]) == n {
		return f
	}

	pow10N := math.Pow10(n)
	return math.Trunc(f*pow10N) / pow10N
}
