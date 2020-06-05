package date

import (
	"testing"
)

func TestStrToTime(t *testing.T) {

	tm := StrToTime("2019-11-14 20:24:51", "Y-m-d H:i:s")
	if tm.String() != "2019-11-14 20:24:51 +0800 CST" {
		t.Error("test StrToTime err", tm.String())
	}

	ts := StrToTimestamp("2019-11-14 20:24:51", TimeFormat)
	if ts != 1573734291 {
		t.Error("test StrToTimestamp err", ts)
	}

	st := StrToTime("2019-11-14T12:24:51Z", "Y-m-dTH:i:sZ")
	if st.Format(TimeFormat) != "2019-11-14 20:24:51" {
		t.Error("test StrToTime with UTC str err", st.Format(TimeFormat))
	}
}
