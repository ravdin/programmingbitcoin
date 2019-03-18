package code_ch03

import (
	"testing"
)

func TestOnCurve(t *testing.T) {
	var prime int64 = 223
	a := NewFieldElement(0, prime)
	b := NewFieldElement(7, prime)
	validPoints := [][]int64{
		[]int64{192, 105},
		[]int64{17, 56},
		[]int64{1, 193},
	}
	invalidPoints := [][]int64{
		[]int64{200, 119},
		[]int64{42, 99},
	}
	for _, validPoint := range validPoints {
		x := NewFieldElement(validPoint[0], prime)
		y := NewFieldElement(validPoint[1], prime)
		_, err := NewPoint(x, y, a, b)
		if err != nil {
			t.Errorf("%v\n", err)
		}
	}
	for _, invalidPoint := range invalidPoints {
		x := NewFieldElement(invalidPoint[0], prime)
		y := NewFieldElement(invalidPoint[1], prime)
		p, err := NewPoint(x, y, a, b)
		if p != nil || err == nil {
			t.Errorf("Expected error for invalid point %v", p)
		}
	}
}
