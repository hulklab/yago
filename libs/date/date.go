package date

import (
	"strings"
	"time"
)

const (
	TimeFormat         = "2006-01-02 15:04:05"
	TimeFormatYYYYMMDD = "20060102"
)

// 用来格式化时间
var TimeFormatMap = map[string]string{
	"Y": "2006",
	"n": "1",
	"m": "01",
	"j": "2",
	"d": "02",
	"H": "15",
	"i": "04",
	"s": "05",
	"J": "05.000",       // 毫秒
	"Q": "05.000000",    // 微秒
	"K": "05.000000000", // 纳秒
}

/**
 * 类似php的date()函数
 *
 * @param
 * @return
 *
 */
func Date(format string, timestamp ...int64) string {
	newFormat := format
	for k, v := range TimeFormatMap {
		newFormat = strings.Replace(newFormat, k, v, 1)
	}
	var tm time.Time
	if len(timestamp) > 0 {
		tm = time.Unix(timestamp[0], 0)
	} else {
		tm = time.Now()
	}
	return tm.Format(newFormat)
}

func Strtotime(datetime string, format string) int64 {
	newFormat := format
	for k, v := range TimeFormatMap {
		newFormat = strings.Replace(newFormat, k, v, 1)
	}
	local, _ := time.LoadLocation("Local")
	theTime, _ := time.ParseInLocation(newFormat, datetime, local)
	return theTime.Unix()
}

func Now() string {
	return time.Now().Format(TimeFormat)
}
