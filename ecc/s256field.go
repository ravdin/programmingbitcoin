package ecc

import (
	"fmt"
	"math/big"
)

// S256Field allows finite field math for 256 bit integers.
type S256Field struct {
	Num   *big.Int
	Prime *big.Int
}

func newS256FieldFromInt64(num int64, prime int64) *S256Field {
	return NewS256Field(big.NewInt(num), big.NewInt(prime))
}

func NewS256Field(num *big.Int, prime *big.Int) *S256Field {
	if num.Sign() < 0 || num.Cmp(prime) >= 0 {
		panic(fmt.Sprintf("Num %d not in valid field range", num))
	}
	return &S256Field{Num: num, Prime: prime}
}

func (field *S256Field) String() string {
	return fmt.Sprintf("S256Field(%d)(%d)", field.Num, field.Prime)
}

func (field *S256Field) Eq(other FieldInteger) bool {
	if other == nil {
		return false
	}
	o := other.(*S256Field)
	return field.Num.Cmp(o.Num) == 0 && field.Prime.Cmp(o.Prime) == 0
}

func (field *S256Field) Ne(other FieldInteger) bool {
	return !field.Eq(other)
}

func (field *S256Field) Add(x, y FieldInteger) FieldInteger {
	fx, fy := x.(*S256Field), y.(*S256Field)
	if fx.Prime.Cmp(fy.Prime) != 0 {
		panic("Cannot add two numbers in different Fields")
	}
	var num = new(big.Int)
	num.Add(fx.Num, fy.Num).Mod(num, fx.Prime)
	*field = S256Field{Num: num, Prime: fx.Prime}
	return field
}

func (field *S256Field) Sub(x, y FieldInteger) FieldInteger {
	fx, fy := x.(*S256Field), y.(*S256Field)
	if fx.Prime.Cmp(fy.Prime) != 0 {
		panic("Cannot subtract two numbers in different Fields")
	}
	var num = new(big.Int)
	num.Sub(fx.Num, fy.Num).Mod(num, fx.Prime)
	if num.Sign() < 0 {
		num.Add(num, fx.Prime)
	}
	*field = S256Field{Num: num, Prime: fx.Prime}
	return field
}

func (field *S256Field) Mul(x, y FieldInteger) FieldInteger {
	fx, fy := x.(*S256Field), y.(*S256Field)
	if fx.Prime.Cmp(fy.Prime) != 0 {
		panic("Cannot multiply two numbers in different Fields")
	}
	var num = new(big.Int)
	num.Mul(fx.Num, fy.Num).Mod(num, fx.Prime)
	*field = S256Field{Num: num, Prime: fx.Prime}
	return field
}

func (field *S256Field) Div(x, y FieldInteger) FieldInteger {
	fx, fy := x.(*S256Field), y.(*S256Field)
	if fx.Prime.Cmp(fy.Prime) != 0 {
		panic("Cannot divide two numbers in different Fields")
	}
	/*
	 * field.num and other.num are the actual values
	 * field.prime is what we need to mod against
	 * use fermat's little theorem:
	 * field.num**(p-1) % p == 1
	 * this means:
	 * 1/n == pow(n, p-2, p)
	 */
	var num = new(big.Int)
	var b = new(big.Int)
	var e = new(big.Int)
	e.Sub(fx.Prime, big.NewInt(2))
	b.Exp(fy.Num, e, fx.Prime)
	num.Mul(fx.Num, b).Mod(num, fx.Prime)
	*field = S256Field{Num: num, Prime: fx.Prime}
	return field
}

func (field *S256Field) Pow(n FieldInteger, exponent *big.Int) FieldInteger {
	f := n.(*S256Field)
	var num = new(big.Int)
	var e = new(big.Int)
	var m = new(big.Int)
	m.Sub(f.Prime, big.NewInt(1))
	e.Mod(exponent, m)
	num.Exp(f.Num, e, f.Prime)
	*field = S256Field{Num: num, Prime: field.Prime}
	return field
}

func (field *S256Field) Cmul(n FieldInteger, coefficient *big.Int) FieldInteger {
	var num = new(big.Int)
	var c = new(big.Int)
	f := n.(*S256Field)
	c.Mod(coefficient, f.Prime)
	num.Mul(f.Num, c).Mod(num, f.Prime)
	*field = S256Field{Num: num, Prime: f.Prime}
	return field
}

func (field *S256Field) Copy() FieldInteger {
	num := new(big.Int)
	num.Set(field.Num)
	return &S256Field{num, field.Prime}
}

func (field *S256Field) Set(n FieldInteger) FieldInteger {
	f := n.(*S256Field)
	field.Num.Set(f.Num)
	field.Prime.Set(f.Prime)
	return field
}

func (field *S256Field) Sqrt() *S256Field {
	e := new(big.Int)
	e.Add(field.Prime, big.NewInt(1))
	e.Div(e, big.NewInt(4))
	result := field.Copy()
	return result.Pow(result, e).(*S256Field)
}
