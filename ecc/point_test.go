package ecc

import (
	"math/big"
	"testing"
)

func TestPoint(t *testing.T) {
	t.Run("TestNe", func(t *testing.T) {
		a, _ := newPointFromInts(3, -7, 5, 7)
		b, _ := newPointFromInts(18, 77, 5, 7)
		if a.Eq(b) {
			t.Errorf("Expected a != b")
		}
		if !a.Eq(a) {
			t.Errorf("Expected a == a")
		}
	})

	t.Run("TestOnCurve", func(t *testing.T) {
		_, err := newPointFromInts(2, 4, 5, 7)
		if err == nil {
			t.Errorf("Point is not on curve, expected error!")
		}
		// These should not raise an error
		p1, err := newPointFromInts(3, -7, 5, 7)
		if p1 == nil || err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		p2, err := newPointFromInts(18, 77, 5, 7)
		if p2 == nil || err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
	})

	t.Run("TestAdd", func(t *testing.T) {
		a, _ := newPointAtInfinity(5, 7)
		b, _ := newPointFromInts(2, 5, 5, 7)
		c, _ := newPointFromInts(2, -5, 5, 7)
		d, _ := newPointFromInts(3, 7, 5, 7)
		e, _ := newPointFromInts(-1, -1, 5, 7)
		f, _ := newPointFromInts(2, -5, 5, 7)
		g, _ := newPointFromInts(-1, 1, 5, 7)
		h, _ := newPointFromInts(18, -77, 5, 7)
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

// Int wrapper for testing.
type intWrapper struct {
	n int64
}

func newIntWrapper(n int64) *intWrapper {
	return &intWrapper{n: n}
}

func newPointFromInts(x int64, y int64, a int64, b int64) (*Point, error) {
	return NewPoint(newIntWrapper(x), newIntWrapper(y), newIntWrapper(a), newIntWrapper(b))
}

func newPointAtInfinity(a int64, b int64) (*Point, error) {
	return NewPoint(nil, nil, newIntWrapper(a), newIntWrapper(b))
}

func (self *intWrapper) Eq(other FieldInteger) bool {
	o := other.(*intWrapper)
	return self.n == o.n
}

func (self *intWrapper) Ne(other FieldInteger) bool {
	o := other.(*intWrapper)
	return self.n != o.n
}

func (self *intWrapper) Add(other FieldInteger) FieldInteger {
	o := other.(*intWrapper)
	return newIntWrapper(self.n + o.n)
}

func (self *intWrapper) Sub(other FieldInteger) FieldInteger {
	o := other.(*intWrapper)
	return newIntWrapper(self.n - o.n)
}

func (self *intWrapper) Mul(other FieldInteger) FieldInteger {
	o := other.(*intWrapper)
	return newIntWrapper(self.n * o.n)
}

func (self *intWrapper) Div(other FieldInteger) FieldInteger {
	o := other.(*intWrapper)
	return newIntWrapper(self.n / o.n)
}

func (self *intWrapper) Pow(exponent *big.Int) FieldInteger {
	result := new(big.Int)
	result.Exp(big.NewInt(self.n), exponent, nil)
	return newIntWrapper(result.Int64())
}

func (self *intWrapper) Rmul(coeff *big.Int) FieldInteger {
	return newIntWrapper(coeff.Int64() * self.n)
}
