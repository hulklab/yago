package date

import (
	"strings"
	"time"
)

const (
	TimeFormat         = "2006-01-02 15:04:05"
	TimeFormatYYYYMMDD = "20060102"
	TimeFormatUTC      = "2006-01-02T15:04:05Z"
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
 * Date("Y-m-d H:i:s",1588149134)
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

// StrToTime("2020-04-29 16:32:14",date.TimeFormat)
// Return Local Location Time
func StrToTime(datetime string, format string) time.Time {
	newFormat := format
	for k, v := range TimeFormatMap {
		newFormat = strings.Replace(newFormat, k, v, 1)
	}

	if newFormat == TimeFormat {
		theTime, _ := time.ParseInLocation(newFormat, datetime, time.Local)
		return theTime
	}

	// 优先判断是否为 UTC 时间
	if utcT, err := time.Parse(time.RFC3339, datetime); err == nil {
		return utcT.Local()
	}

	theTime, _ := time.ParseInLocation(newFormat, datetime, time.Local)
	return theTime
}

// StrToTimestamp("2020-04-29 16:32:14",date.TimeFormat)
func StrToTimestamp(datetime string, format string) int64 {
	theTime := StrToTime(datetime, format)
	return theTime.Unix()
}

func Now() string {
	return time.Now().Format(TimeFormat)
}
