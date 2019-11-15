package date

import (
	"fmt"
	"testing"
)

func TestStrToTime(t *testing.T) {

	tm := StrToTime("2019-11-14 20:24:51", "Y-m-d H:i:s")
	fmt.Println(tm.String())

	ts := StrToTimestamp("2019-11-14 20:24:51", TimeFormat)
	fmt.Println(ts)
}
