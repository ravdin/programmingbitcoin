package ecc

import (
	"fmt"
	"math/big"
)

// Similar to FieldElement, but we will allow 256 bit integers.
type S256Field struct {
	Num   *big.Int
	Prime *big.Int
}

func NewS256FieldFromInt64(num int64, prime int64) *S256Field {
	return NewS256Field(big.NewInt(num), big.NewInt(prime))
}

func NewS256Field(num *big.Int, prime *big.Int) *S256Field {
	if num.Sign() < 0 || num.Cmp(prime) >= 0 {
		panic(fmt.Sprintf("Num %d not in valid field range", num))
	}
	return &S256Field{Num: num, Prime: prime}
}

func (self *S256Field) String() string {
	return fmt.Sprintf("S256Field(%d)(%d)", self.Num, self.Prime)
}

func (self *S256Field) Eq(other FieldInteger) bool {
	if other == nil {
		return false
	}
	o := other.(*S256Field)
	return self.Num.Cmp(o.Num) == 0 && self.Prime.Cmp(o.Prime) == 0
}

func (self *S256Field) Ne(other FieldInteger) bool {
	return !self.Eq(other)
}

func (z *S256Field) Add(x, y FieldInteger) FieldInteger {
	fx, fy := x.(*S256Field), y.(*S256Field)
	if fx.Prime.Cmp(fy.Prime) != 0 {
		panic("Cannot add two numbers in different Fields")
	}
	var num = new(big.Int)
	num.Add(fx.Num, fy.Num).Mod(num, fx.Prime)
	*z = S256Field{Num: num, Prime: fx.Prime}
	return z
}

func (z *S256Field) Sub(x, y FieldInteger) FieldInteger {
	fx, fy := x.(*S256Field), y.(*S256Field)
	if fx.Prime.Cmp(fy.Prime) != 0 {
		panic("Cannot subtract two numbers in different Fields")
	}
	var num = new(big.Int)
	num.Sub(fx.Num, fy.Num).Mod(num, fx.Prime)
	if num.Sign() < 0 {
		num.Add(num, fx.Prime)
	}
	*z = S256Field{Num: num, Prime: fx.Prime}
	return z
}

func (z *S256Field) Mul(x, y FieldInteger) FieldInteger {
	fx, fy := x.(*S256Field), y.(*S256Field)
	if fx.Prime.Cmp(fy.Prime) != 0 {
		panic("Cannot multiply two numbers in different Fields")
	}
	var num = new(big.Int)
	num.Mul(fx.Num, fy.Num).Mod(num, fx.Prime)
	*z = S256Field{Num: num, Prime: fx.Prime}
	return z
}

func (z *S256Field) Div(x, y FieldInteger) FieldInteger {
	fx, fy := x.(*S256Field), y.(*S256Field)
	if fx.Prime.Cmp(fy.Prime) != 0 {
		panic("Cannot divide two numbers in different Fields")
	}
	/*
	 * self.num and other.num are the actual values
	 * self.prime is what we need to mod against
	 * use fermat's little theorem:
	 * self.num**(p-1) % p == 1
	 * this means:
	 * 1/n == pow(n, p-2, p)
	 */
	var num = new(big.Int)
	var b = new(big.Int)
	var e = new(big.Int)
	e.Sub(fx.Prime, big.NewInt(2))
	b.Exp(fy.Num, e, fx.Prime)
	num.Mul(fx.Num, b).Mod(num, fx.Prime)
	*z = S256Field{Num: num, Prime: fx.Prime}
	return z
}

func (z *S256Field) Pow(n FieldInteger, exponent *big.Int) FieldInteger {
	field := n.(*S256Field)
	var num = new(big.Int)
	var e = new(big.Int)
	var m = new(big.Int)
	m.Sub(field.Prime, big.NewInt(1))
	e.Mod(exponent, m)
	num.Exp(field.Num, e, field.Prime)
	*z = S256Field{Num: num, Prime: field.Prime}
	return z
}

func (z *S256Field) Cmul(n FieldInteger, coefficient *big.Int) FieldInteger {
	var num = new(big.Int)
	var c = new(big.Int)
	field := n.(*S256Field)
	c.Mod(coefficient, field.Prime)
	num.Mul(field.Num, c).Mod(num, field.Prime)
	*z = S256Field{Num: num, Prime: field.Prime}
	return z
}

func (z *S256Field) Copy() FieldInteger {
	num := new(big.Int)
	num.Set(z.Num)
	return &S256Field{num, z.Prime}
}

func (z *S256Field) Set(n FieldInteger) FieldInteger {
	field := n.(*S256Field)
	z.Num.Set(field.Num)
	z.Prime.Set(field.Prime)
	return z
}

func (self *S256Field) Sqrt() *S256Field {
	var e *big.Int = new(big.Int)
	e.Add(self.Prime, big.NewInt(1))
	e.Div(e, big.NewInt(4))
	result := self.Copy()
	return result.Pow(result, e).(*S256Field)
}
