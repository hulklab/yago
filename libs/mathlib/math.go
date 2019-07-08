package mathlib

import "math"

/**
 * 类似php的round()方法
 *
 * @param f float64
 * @param n int 取小数点后几位
 * @return float64
 *
 */
func Round(f float64, n int) float64 {
	pow10N := math.Pow10(n)
	return math.Trunc((f+0.5/pow10N)*pow10N) / pow10N
}
