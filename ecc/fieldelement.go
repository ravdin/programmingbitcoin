package ecc

import (
	"fmt"
	"math/big"
)

// Finite field impmentation for testing. This type has a limit of 64 bits.
type FieldElement struct {
	Num   int64
	Prime int64
}

func NewFieldElement(num int64, prime int64) *FieldElement {
	if num >= prime || num < 0 {
		panic(fmt.Sprintf("Num %d not in field range 0 to %d", num, prime-1))
	}
	return &FieldElement{Num: num, Prime: prime}
}

func (self *FieldElement) String() string {
	return fmt.Sprintf("FieldElement_%d(%d)", self.Prime, self.Num)
}

func (self *FieldElement) Eq(other FieldInteger) bool {
	o := other.(*FieldElement)
	return self.Num == o.Num && self.Prime == o.Prime
}

func (self *FieldElement) Ne(other FieldInteger) bool {
	o := other.(*FieldElement)
	return self.Num != o.Num || self.Prime != o.Prime
}

func (z *FieldElement) Add(x, y FieldInteger) FieldInteger {
	elemx, elemy := x.(*FieldElement), y.(*FieldElement)
	if elemx.Prime != elemy.Prime {
		panic("Cannot add two numbers in different Fields")
	}
	num := (elemx.Num + elemy.Num) % elemx.Prime
	*z = FieldElement{Num: num, Prime: elemx.Prime}
	return z
}

func (z *FieldElement) Sub(x, y FieldInteger) FieldInteger {
	elemx, elemy := x.(*FieldElement), y.(*FieldElement)
	if elemx.Prime != elemy.Prime {
		panic("Cannot subtract two numbers in different Fields")
	}
	num := (elemx.Num - elemy.Num + elemx.Prime) % elemx.Prime
	*z = FieldElement{Num: num, Prime: elemx.Prime}
	return z
}

func (z *FieldElement) Mul(x, y FieldInteger) FieldInteger {
	elemx, elemy := x.(*FieldElement), y.(*FieldElement)
	if elemx.Prime != elemy.Prime {
		panic("Cannot multiply two numbers in different Fields")
	}
	num := (elemx.Num * elemy.Num) % elemx.Prime
	*z = FieldElement{Num: num, Prime: elemx.Prime}
	return z
}

func (z *FieldElement) Div(x, y FieldInteger) FieldInteger {
	elemx, elemy := x.(*FieldElement), y.(*FieldElement)
	if elemx.Prime != elemy.Prime {
		panic("Cannot divide two numbers in different Fields")
	}
	num := (elemx.Num * intPow(elemy.Num, elemx.Prime-2, elemx.Prime)) % elemx.Prime
	*z = FieldElement{Num: num, Prime: elemx.Prime}
	return z
}

func (z *FieldElement) Pow(n FieldInteger, exponent *big.Int) FieldInteger {
	field := n.(*FieldElement)
	e := (exponent.Int64() + field.Prime - 1) % (field.Prime - 1)
	num := intPow(field.Num, e, field.Prime)
	*z = FieldElement{Num: num, Prime: field.Prime}
	return z
}

func (z *FieldElement) Cmul(n FieldInteger, coefficient *big.Int) FieldInteger {
	panic("Not implemented")
}

func (z *FieldElement) Copy() FieldInteger {
	panic("Not implemented")
}

func (z *FieldElement) Set(n FieldInteger) FieldInteger {
	panic("Not implemented")
}

// Integer exponent (doesn't exist in golang's math package).
func intPow(num int64, exponent int64, mod int64) int64 {
	if exponent < 0 {
		panic("Exponent cannot be negative")
	}
	var result int64 = 1
	for exponent > 0 {
		if exponent&1 == 1 {
			result = (result * num) % mod
		}
		num = (num * num) % mod
		if num == 1 {
			break
		}
		exponent >>= 1
	}
	return result
}
