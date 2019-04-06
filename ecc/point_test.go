package ecc

import (
	"testing"
)

func withIntWrapper() FieldIntegerConverter {
	return func(n interface{}) FieldInteger {
		return newIntWrapper(n.(int))
	}
}

func TestPoint(t *testing.T) {
	option := withIntWrapper()

	t.Run("TestNe", func(t *testing.T) {
		a, _ := NewPoint(3, -7, 5, 7, option)
		b, _ := NewPoint(18, 77, 5, 7, option)
		if a.Eq(b) {
			t.Errorf("Expected a != b")
		}
		if !a.Eq(a) {
			t.Errorf("Expected a == a")
		}
	})

	t.Run("TestOnCurve", func(t *testing.T) {
		_, err := NewPoint(2, 4, 5, 7, option)
		if err == nil {
			t.Errorf("Point is not on curve, expected error!")
		}
		// These should not raise an error
		p1, err := NewPoint(3, -7, 5, 7, option)
		if p1 == nil || err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		p2, err := NewPoint(18, 77, 5, 7, option)
		if p2 == nil || err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
	})

	t.Run("TestAdd", func(t *testing.T) {
		a, _ := NewPoint(nil, nil, 5, 7, option)
		b, _ := NewPoint(2, 5, 5, 7, option)
		c, _ := NewPoint(2, -5, 5, 7, option)
		d, _ := NewPoint(3, 7, 5, 7, option)
		e, _ := NewPoint(-1, -1, 5, 7, option)
		f, _ := NewPoint(2, -5, 5, 7, option)
		g, _ := NewPoint(-1, 1, 5, 7, option)
		h, _ := NewPoint(18, -77, 5, 7, option)
		tests := [][]*Point{
			{
				a.Add(b), b,
			},
			{
				b.Add(a), b,
			},
			{
				b.Add(c), a,
			},
			{
				d.Add(e), f,
			},
			{
				g.Add(g), h,
			},
		}
		for _, test := range tests {
			actual := test[0]
			expected := test[1]
			if !actual.Eq(expected) {
				t.Errorf("Expected %v, got %v", expected, actual)
			}
		}
	})
}
