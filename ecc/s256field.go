package ecc

import (
	"fmt"
	"math/big"
)

// s256Field allows finite field math for 256 bit integers.
type s256Field struct {
	Num   *big.Int
	Prime *big.Int
}

func newS256FieldFromInt64(num int64, prime int64) *s256Field {
	return newS256Field(big.NewInt(num), big.NewInt(prime))
}

func newS256Field(num *big.Int, prime *big.Int) *s256Field {
	if num.Sign() < 0 || num.Cmp(prime) >= 0 {
		panic(fmt.Sprintf("Num %d not in valid field range", num))
	}
	return &s256Field{Num: num, Prime: prime}
}

func (field *s256Field) String() string {
	return fmt.Sprintf("s256Field(%d)(%d)", field.Num, field.Prime)
}

func (field *s256Field) Eq(other FieldInteger) bool {
	if other == nil {
		return false
	}
	o := other.(*s256Field)
	return field.Num.Cmp(o.Num) == 0 && field.Prime.Cmp(o.Prime) == 0
}

func (field *s256Field) Ne(other FieldInteger) bool {
	return !field.Eq(other)
}

func (field *s256Field) Add(x, y FieldInteger) FieldInteger {
	fx, fy := x.(*s256Field), y.(*s256Field)
	if fx.Prime.Cmp(fy.Prime) != 0 {
		panic("Cannot add two numbers in different Fields")
	}
	var num = new(big.Int)
	num.Add(fx.Num, fy.Num).Mod(num, fx.Prime)
	*field = s256Field{Num: num, Prime: fx.Prime}
	return field
}

func (field *s256Field) Sub(x, y FieldInteger) FieldInteger {
	fx, fy := x.(*s256Field), y.(*s256Field)
	if fx.Prime.Cmp(fy.Prime) != 0 {
		panic("Cannot subtract two numbers in different Fields")
	}
	var num = new(big.Int)
	num.Sub(fx.Num, fy.Num).Mod(num, fx.Prime)
	if num.Sign() < 0 {
		num.Add(num, fx.Prime)
	}
	*field = s256Field{Num: num, Prime: fx.Prime}
	return field
}

func (field *s256Field) Mul(x, y FieldInteger) FieldInteger {
	fx, fy := x.(*s256Field), y.(*s256Field)
	if fx.Prime.Cmp(fy.Prime) != 0 {
		panic("Cannot multiply two numbers in different Fields")
	}
	var num = new(big.Int)
	num.Mul(fx.Num, fy.Num).Mod(num, fx.Prime)
	*field = s256Field{Num: num, Prime: fx.Prime}
	return field
}

func (field *s256Field) Div(x, y FieldInteger) FieldInteger {
	fx, fy := x.(*s256Field), y.(*s256Field)
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
	*field = s256Field{Num: num, Prime: fx.Prime}
	return field
}

func (field *s256Field) Pow(n FieldInteger, exponent *big.Int) FieldInteger {
	f := n.(*s256Field)
	var num = new(big.Int)
	var e = new(big.Int)
	var m = new(big.Int)
	m.Sub(f.Prime, big.NewInt(1))
	e.Mod(exponent, m)
	num.Exp(f.Num, e, f.Prime)
	*field = s256Field{Num: num, Prime: f.Prime}
	return field
}

func (field *s256Field) Cmul(n FieldInteger, coefficient *big.Int) FieldInteger {
	var num = new(big.Int)
	var c = new(big.Int)
	f := n.(*s256Field)
	c.Mod(coefficient, f.Prime)
	num.Mul(f.Num, c).Mod(num, f.Prime)
	*field = s256Field{Num: num, Prime: f.Prime}
	return field
}

func (field *s256Field) Copy() FieldInteger {
	num := new(big.Int)
	num.Set(field.Num)
	return &s256Field{num, field.Prime}
}

func (field *s256Field) Set(n FieldInteger) FieldInteger {
	f := n.(*s256Field)
	field.Num.Set(f.Num)
	field.Prime.Set(f.Prime)
	return field
}

func (field *s256Field) Sqrt() *s256Field {
	e := new(big.Int)
	e.Add(field.Prime, big.NewInt(1))
	e.Div(e, big.NewInt(4))
	result := field.Copy()
	return result.Pow(result, e).(*s256Field)
}
