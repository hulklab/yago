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
