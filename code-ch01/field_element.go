package code_ch01

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

func (self *FieldElement) Eq(other *FieldElement) bool {
	return self.Num == other.Num && self.Prime == other.Prime
}

func (self *FieldElement) Ne(other *FieldElement) bool {
	panic("Not implemented")
}

func (self *FieldElement) Add(other *FieldElement) *FieldElement {
	if self.Prime != other.Prime {
		panic("Cannot add two numbers in different Fields")
	}
	num := (self.Num + other.Num) % self.Prime
	return NewFieldElement(num, self.Prime)
}

func (self *FieldElement) Sub(other *FieldElement) *FieldElement {
	if self.Prime != other.Prime {
		panic("Cannot subtract two numbers in different Fields")
	}
	// self.num and other.num are the actual values
	// self.prime is what we need to mod against
	// We return an element of the same class
	panic("Not implemented")
}

func (self *FieldElement) Pow(exponent int64) *FieldElement {
	n := exponent % (self.Prime - 1)
	num := intPow(self.Num, n, self.Prime)
	return NewFieldElement(num, self.Prime)
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
