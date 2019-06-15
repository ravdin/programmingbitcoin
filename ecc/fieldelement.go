package ecc

import (
	"fmt"
	"math/big"
)

// Finite field impmentation for testing. This type has a limit of 64 bits.
type fieldElement struct {
	Num   int64
	Prime int64
}

func newFieldElement(num int64, prime int64) *fieldElement {
	if num >= prime || num < 0 {
		panic(fmt.Sprintf("Num %d not in field range 0 to %d", num, prime-1))
	}
	return &fieldElement{Num: num, Prime: prime}
}

func (elem *fieldElement) String() string {
	return fmt.Sprintf("FieldElement_%d(%d)", elem.Prime, elem.Num)
}

func (elem *fieldElement) Eq(other FieldInteger) bool {
	o := other.(*fieldElement)
	return elem.Num == o.Num && elem.Prime == o.Prime
}

func (elem *fieldElement) Ne(other FieldInteger) bool {
	o := other.(*fieldElement)
	return elem.Num != o.Num || elem.Prime != o.Prime
}

func (elem *fieldElement) Add(x, y FieldInteger) FieldInteger {
	elemx, elemy := x.(*fieldElement), y.(*fieldElement)
	if elemx.Prime != elemy.Prime {
		panic("Cannot add two numbers in different Fields")
	}
	num := (elemx.Num + elemy.Num) % elemx.Prime
	*elem = fieldElement{Num: num, Prime: elemx.Prime}
	return elem
}

func (elem *fieldElement) Sub(x, y FieldInteger) FieldInteger {
	elemx, elemy := x.(*fieldElement), y.(*fieldElement)
	if elemx.Prime != elemy.Prime {
		panic("Cannot subtract two numbers in different Fields")
	}
	num := (elemx.Num - elemy.Num + elemx.Prime) % elemx.Prime
	*elem = fieldElement{Num: num, Prime: elemx.Prime}
	return elem
}

func (elem *fieldElement) Mul(x, y FieldInteger) FieldInteger {
	elemx, elemy := x.(*fieldElement), y.(*fieldElement)
	if elemx.Prime != elemy.Prime {
		panic("Cannot multiply two numbers in different Fields")
	}
	num := (elemx.Num * elemy.Num) % elemx.Prime
	*elem = fieldElement{Num: num, Prime: elemx.Prime}
	return elem
}

func (elem *fieldElement) Div(x, y FieldInteger) FieldInteger {
	elemx, elemy := x.(*fieldElement), y.(*fieldElement)
	if elemx.Prime != elemy.Prime {
		panic("Cannot divide two numbers in different Fields")
	}
	num := (elemx.Num * intPow(elemy.Num, elemx.Prime-2, elemx.Prime)) % elemx.Prime
	*elem = fieldElement{Num: num, Prime: elemx.Prime}
	return elem
}

func (elem *fieldElement) Pow(n FieldInteger, exponent *big.Int) FieldInteger {
	field := n.(*fieldElement)
	e := (exponent.Int64() + field.Prime - 1) % (field.Prime - 1)
	num := intPow(field.Num, e, field.Prime)
	*elem = fieldElement{Num: num, Prime: field.Prime}
	return elem
}

func (elem *fieldElement) Cmul(n FieldInteger, coefficient *big.Int) FieldInteger {
	panic("Not implemented")
}

func (elem *fieldElement) Copy() FieldInteger {
	panic("Not implemented")
}

func (elem *fieldElement) Set(n FieldInteger) FieldInteger {
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
