package mathlib

import (
	"testing"
)

func TestRound(t *testing.T) {
	var tests = []struct {
		in     float64
		length int
		out    float64
	}{
		{1.16, 1, 1.2},
		{1.11, 1, 1.1},
		{1.15, 2, 1.15},
	}

	for _, v := range tests {
		get := Round(v.in, v.length)
		if get != v.out {
			t.Errorf("Floor(%v,%v) = %v ,expected %v\n",
				v.in, v.length, get, v.out)
		}

	}
}

func TestFloor(t *testing.T) {
	var tests = []struct {
		in     float64
		length int
		out    float64
	}{
		{1.16, 1, 1.1},
		{1.11, 1, 1.1},
		{1.15, 2, 1.15},
		{1.19, 2, 1.19},
		{1.199999, 1, 1.1},
		{1.15, 0, 1},
		{1.15, 3, 1.15},
	}

	for _, v := range tests {
		get := Floor(v.in, v.length)
		if get != v.out {
			t.Errorf("Floor(%v,%v) = %v ,expected %v\n",
				v.in, v.length, get, v.out)
		}

	}
}
