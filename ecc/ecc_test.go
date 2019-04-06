package ecc

import (
	"math/big"
	"testing"
)

func TestECC(t *testing.T) {
	var a int64 = 0
	var b int64 = 7
	var prime int64 = 223
	f223 := func(n interface{}) FieldInteger {
		return NewS256FieldFromInt64(n.(int64), prime)
	}

	t.Run("Test on curve", func(t *testing.T) {
		// tests the following points whether they are on the curve or not
		// on curve y^2=x^3-7 over F_223:
		// (192,105) (17,56) (200,119) (1,193) (42,99)
		validPoints := [][]int64{
			{192, 105}, {17, 56}, {1, 193},
		}
		invalidPoints := [][]int64{
			{200, 119}, {42, 99},
		}
		for _, item := range validPoints {
			p, err := NewPoint(item[0], item[1], a, b, f223)
			if p == nil || err != nil {
				t.Errorf("Unexpected error %v!", err)
			}
		}
		for _, item := range invalidPoints {
			p, err := NewPoint(item[0], item[1], a, b, f223)
			if p != nil || err == nil {
				t.Errorf("Point %v is invalid!", p)
			}
		}
	})

	t.Run("Test additions", func(t *testing.T) {
		// tests the following additions on curve y^2=x^3-7 over F_223:
		// (192,105) + (17,56)
		// (47,71) + (117,141)
		// (143,98) + (76,66)
		additions := [][]int64{
			{192, 105, 17, 56, 170, 142},
			{47, 71, 117, 141, 60, 139},
			{143, 98, 76, 66, 47, 71},
		}
		for _, item := range additions {
			p1, _ := NewPoint(item[0], item[1], a, b, f223)
			p2, _ := NewPoint(item[2], item[3], a, b, f223)
			expected, _ := NewPoint(item[4], item[5], a, b, f223)
			actual := p1.Add(p2)
			if !actual.Eq(expected) {
				t.Errorf("Expected %v, got %v", expected, actual)
			}
		}
	})

	t.Run("Test scalar multiplications", func(t *testing.T) {
		// tests the following scalar multiplications
		// 2*(192,105)
		// 2*(143,98)
		// 2*(47,71)
		// 4*(47,71)
		// 8*(47,71)
		// 21*(47,71)
		multiplications := [][]int64{
			{2, 192, 105, 49, 71},
			{2, 143, 98, 64, 168},
			{2, 47, 71, 36, 111},
			{4, 47, 71, 194, 51},
			{8, 47, 71, 116, 55},
		}
		for _, item := range multiplications {
			p1, _ := NewPoint(item[1], item[2], a, b, f223)
			actual, _ := NewPoint(item[3], item[4], a, b, f223)
			expected := p1.Rmul(big.NewInt(item[0]))
			if !actual.Eq(expected) {
				t.Errorf("Expected %v, got %v", expected, actual)
			}
		}
    // Test for infinity case.
    p1, _ := NewPoint(int64(47), int64(71), a, b, f223)
    actual := p1.Rmul(big.NewInt(21))
    expected, _ := NewPoint(nil, nil, a, b,f223)
    if !actual.Eq(expected) {
      t.Errorf("Expected %v, got %v", expected, actual)
    }
	})
}
