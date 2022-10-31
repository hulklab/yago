package str

import (
	"fmt"
	"testing"
)

func TestSplit(t *testing.T) {

	s := "a, b,c\n, d\t"
	ss := Split(s)
	for _, v := range ss {
		fmt.Printf("^%s$", v)
	}

}

func TestSubstr(t *testing.T) {

	var tests = []struct {
		in     string
		start  int
		length int
		out    string
	}{
		{"", 0, 1, ""},
		{"我是中国人", 0, 3, "我是中"},
		{"我是中国人", 0, 10, "我是中国人"},
		{"我是中国人", -3, 3, "中国人"},
		{"我是中国人", -8, 3, "我是中"},
		{"中国 is great", 3, 2, "is"},
		{"中国 is great", -5, 2, "gr"},
	}

	for _, v := range tests {
		get := Substr(v.in, v.start, v.length)
		if get != v.out {
			t.Errorf("Substr(%v,%v,%v) = %v ,expected %v\n",
				v.in, v.start, v.length, get, v.out)
		}

	}
}
