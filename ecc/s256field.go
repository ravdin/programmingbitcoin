package ecc

import (
	"encoding/hex"
	"fmt"
	"math/big"
)

// 2^256 - 2^32 - 2^9 - 2^8 - 2^7 - 2^6 - 2^4 - 1
const P string = "FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFEFFFFFC2F"

// Similar to FieldElement, but we will allow 256 bit integers.
type S256Field struct {
	Num   *big.Int
	Prime *big.Int
}

func NewS256Field(num *big.Int) *S256Field {
	decoded, _ := hex.DecodeString(P)
	var prime big.Int
	prime.SetBytes(decoded)
	if num.Cmp(big.NewInt(0)) < 0 || num.Cmp(&prime) >= 0 {
		panic(fmt.Sprintf("Num %d not in valid field range", num))
	}
	return &S256Field{Num: num, Prime: &prime}
}

func (self *S256Field) String() string {
	return fmt.Sprintf("S256Field(%d)", self.Num)
}

func (self *S256Field) Eq(other FieldInteger) bool {
	o := other.(*S256Field)
	return self.Num.Cmp(o.Num) == 0
}

func (self *S256Field) Ne(other FieldInteger) bool {
	o := other.(*S256Field)
	return self.Num.Cmp(o.Num) != 0
}

func (self *S256Field) Add(other FieldInteger) FieldInteger {
	o := other.(*S256Field)
	var num *big.Int
	num.Set(self.Num).Add(num, o.Num).Mod(num, self.Prime)
	return &S256Field{Num: num, Prime: self.Prime}
}

func (self *S256Field) Sub(other FieldInteger) FieldInteger {
	o := other.(*S256Field)
	var num *big.Int
	num.Set(self.Num).Sub(num, o.Num).Mod(num, self.Prime)
	return &S256Field{Num: num, Prime: self.Prime}
}

func (self *S256Field) Mul(other FieldInteger) FieldInteger {
	o := other.(*S256Field)
	var num *big.Int
	num.Set(self.Num).Mul(num, o.Num).Mod(num, self.Prime)
	return &S256Field{Num: num, Prime: self.Prime}
}

func (self *S256Field) Div(other FieldInteger) FieldInteger {
	o := other.(*S256Field)
	/*
	 * self.num and other.num are the actual values
	 * self.prime is what we need to mod against
	 * use fermat's little theorem:
	 * self.num**(p-1) % p == 1
	 * this means:
	 * 1/n == pow(n, p-2, p)
	 */
	var num, b, e *big.Int
	num.Set(self.Num)
	b.Set(o.Num)
	e.Set(self.Prime).Sub(e, big.NewInt(2))
	b.Exp(e, self.Prime, self.Prime)
	num.Mul(num, b).Mod(num, self.Prime)
	return &S256Field{Num: num, Prime: self.Prime}
}

func (self *S256Field) Pow(exponent *big.Int) FieldInteger {
	var num, m, n *big.Int
	m.Set(self.Prime).Sub(m, big.NewInt(1))
	n.Set(exponent).Add(n, self.Prime).Mod(n, m)
	num.Set(self.Num).Exp(num, n, self.Prime)
	return &S256Field{Num: num, Prime: self.Prime}
}

func (self *S256Field) Rmul(coeff *big.Int) FieldInteger {
	var num, c *big.Int
	c.Set(coeff).Mod(c, self.Prime)
	num.Set(self.Num).Mul(num, c).Mod(num, self.Prime)
	return &S256Field{Num: num, Prime: self.Prime}
}
