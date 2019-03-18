package code_ch03

import (
	"fmt"
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

func (self *FieldElement) Eq(other interface{}) bool {
	o := other.(*FieldElement)
	return self.Num == o.Num && self.Prime == o.Prime
}

func (self *FieldElement) Ne(other interface{}) bool {
	o := other.(*FieldElement)
	return self.Num != o.Num || self.Prime != o.Prime
}

func (self *FieldElement) Add(other interface{}) FieldInteger {
	o := other.(*FieldElement)
	if self.Prime != o.Prime {
		panic("Cannot add two numbers in different Fields")
	}
	num := (self.Num + o.Num) % self.Prime
	return NewFieldElement(num, self.Prime)
}

func (self *FieldElement) Sub(other interface{}) FieldInteger {
	o := other.(*FieldElement)
	if self.Prime != o.Prime {
		panic("Cannot subtract two numbers in different Fields")
	}
	num := (self.Num - o.Num) % self.Prime
	return NewFieldElement(num, self.Prime)
}

func (self *FieldElement) Mul(other interface{}) FieldInteger {
	o := other.(*FieldElement)
	if self.Prime != o.Prime {
		panic("Cannot multiply two numbers in different Fields")
	}
	num := (self.Num * o.Num) % self.Prime
	return NewFieldElement(num, self.Prime)
}

func (self *FieldElement) Div(other interface{}) FieldInteger {
	o := other.(*FieldElement)
	if self.Prime != o.Prime {
		panic("Cannot divide two numbers in different Fields")
	}
	num := (self.Num * intPow(o.Num, self.Prime-2, self.Prime)) % self.Prime
	return NewFieldElement(num, self.Prime)
}

func (self *FieldElement) Pow(exponent int64) FieldInteger {
	n := exponent % (self.Prime - 1)
	num := intPow(self.Num, n, self.Prime)
	return NewFieldElement(num, self.Prime)
}

func (self *FieldElement) Rmul(coeff int64) FieldInteger {
	num := (coeff * self.Num) % self.Prime
	return NewFieldElement(num, self.Prime)
}

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