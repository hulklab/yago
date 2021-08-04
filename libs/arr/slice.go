package arr

import "reflect"

// 判断 needle 是否存在于 haystack 中
func InArray(needle interface{}, haystack interface{}) bool {
	targetValue := reflect.ValueOf(haystack)
	switch reflect.TypeOf(haystack).Kind() {
	case reflect.Slice, reflect.Array:
		for i := 0; i < targetValue.Len(); i++ {
			if targetValue.Index(i).Interface() == needle {
				return true
			}
		}
	}

	return false
}
