package ecc

import (
	"fmt"
	"math/big"
)

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

func (self *FieldElement) Add(other FieldInteger) FieldInteger {
	o := other.(*FieldElement)
	if self.Prime != o.Prime {
		panic("Cannot add two numbers in different Fields")
	}
	num := (self.Num + o.Num) % self.Prime
	return NewFieldElement(num, self.Prime)
}

func (self *FieldElement) Sub(other FieldInteger) FieldInteger {
	o := other.(*FieldElement)
	if self.Prime != o.Prime {
		panic("Cannot subtract two numbers in different Fields")
	}
	num := (self.Num - o.Num + self.Prime) % self.Prime
	return NewFieldElement(num, self.Prime)
}

func (self *FieldElement) Mul(other FieldInteger) FieldInteger {
	o := other.(*FieldElement)
	if self.Prime != o.Prime {
		panic("Cannot multiply two numbers in different Fields")
	}
	num := (self.Num * o.Num) % self.Prime
	return NewFieldElement(num, self.Prime)
}

func (self *FieldElement) Div(other FieldInteger) FieldInteger {
	o := other.(*FieldElement)
	if self.Prime != o.Prime {
		panic("Cannot divide two numbers in different Fields")
	}
	num := (self.Num * intPow(o.Num, self.Prime-2, self.Prime)) % self.Prime
	return NewFieldElement(num, self.Prime)
}

func (self *FieldElement) Pow(exponent *big.Int) FieldInteger {
	n := (exponent.Int64() + self.Prime - 1) % (self.Prime - 1)
	num := intPow(self.Num, n, self.Prime)
	return NewFieldElement(num, self.Prime)
}

func (self *FieldElement) Rmul(coeff *big.Int) FieldInteger {
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
