package ecc

import (
	"math/big"
	"testing"
)

type FieldIntegerOp func(FieldInteger, FieldInteger) bool

func fieldIntegerNeOp() FieldIntegerOp {
	return func(a FieldInteger, b FieldInteger) bool {
		return a.Ne(b)
	}
}

func assertEqual(a FieldInteger, b FieldInteger, t *testing.T) {
	if !a.Eq(b) {
		t.Errorf("%v is not equal to %v\n", a, b)
	}
}

func assertTrue(op FieldIntegerOp, a FieldInteger, b FieldInteger, t *testing.T) {
	if !op(a, b) {
		t.Errorf("Assertion failed!\n")
	}
}

func assertFalse(op FieldIntegerOp, a FieldInteger, b FieldInteger, t *testing.T) {
	if op(a, b) {
		t.Errorf("Assertion failed!\n")
	}
}

// Int wrapper for testing.
type intWrapper struct {
	n int
}

func newIntWrapper(n int) *intWrapper {
	return &intWrapper{n: n}
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
	base := big.NewInt(int64(self.n))
	result := base.Exp(base, exponent, nil)
	return newIntWrapper(int(result.Int64()))
}

func (self *intWrapper) Rmul(coeff *big.Int) FieldInteger {
	return newIntWrapper(int(coeff.Int64()) * self.n)
}
