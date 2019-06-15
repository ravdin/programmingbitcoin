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
				a, b, b,
			},
			{
				b, a, b,
			},
			{
				b, c, a,
			},
			{
				d, e, f,
			},
			{
				g, g, h,
			},
		}
		for _, test := range tests {
			p1, p2, expected := test[0], test[1], test[2]
			actual := new(Point)
			actual.Add(p1, p2)
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

func (i *intWrapper) Eq(other FieldInteger) bool {
	o := other.(*intWrapper)
	return i.n == o.n
}

func (i *intWrapper) Ne(other FieldInteger) bool {
	o := other.(*intWrapper)
	return i.n != o.n
}

func (i *intWrapper) Add(x, y FieldInteger) FieldInteger {
	intx, inty := x.(*intWrapper), y.(*intWrapper)
	*i = intWrapper{n: intx.n + inty.n}
	return i
}

func (i *intWrapper) Sub(x, y FieldInteger) FieldInteger {
	intx, inty := x.(*intWrapper), y.(*intWrapper)
	*i = intWrapper{n: intx.n - inty.n}
	return i
}

func (i *intWrapper) Mul(x, y FieldInteger) FieldInteger {
	intx, inty := x.(*intWrapper), y.(*intWrapper)
	*i = intWrapper{n: intx.n * inty.n}
	return i
}

func (i *intWrapper) Div(x, y FieldInteger) FieldInteger {
	intx, inty := x.(*intWrapper), y.(*intWrapper)
	*i = intWrapper{n: intx.n / inty.n}
	return i
}

func (i *intWrapper) Pow(n FieldInteger, exponent *big.Int) FieldInteger {
	field := n.(*intWrapper)
	result := new(big.Int)
	result.Exp(big.NewInt(field.n), exponent, nil)
	*i = intWrapper{n: result.Int64()}
	return i
}

func (i *intWrapper) Cmul(n FieldInteger, coefficient *big.Int) FieldInteger {
	field := n.(*intWrapper)
	*i = intWrapper{n: coefficient.Int64() * field.n}
	return i
}

func (i *intWrapper) Copy() FieldInteger {
	return &intWrapper{n: i.n}
}

func (i *intWrapper) Set(n FieldInteger) FieldInteger {
	field := n.(*intWrapper)
	i.n = field.n
	return i
}
