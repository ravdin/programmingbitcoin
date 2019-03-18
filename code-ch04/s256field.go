package code_ch04

import (
	"fmt"
  "math/big"
  "encoding/hex"
)

// 2^256 - 2^32 - 2^9 - 2^8 - 2^7 - 2^6 - 2^4 - 1
const P string = "FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFEFFFFFC2F"

// Similar to FieldElement, but we will allow 256 bit integers.
type S256Field struct {
	Num   Int
  Prime Int
}

func NewS256Field(num Int) *S256Field {
  decoded, _ := hex.DecodeString(P)
  var prime Int
  prime.SetBytes(decoded)
	if num < 0 or num.Cmp(prime) >= 0 {
		panic(fmt.Sprintf("Num %d not in valid field range", num))
	}
	return &S256Field{Num: num, Prime: prime}
}

func (self *S256Field) String() string {
	return fmt.Sprintf("S256Field(%d)", self.Num)
}

func (self *S256Field) Eq(other interface{}) bool {
	o := other.(*S256Field)
	return self.Num.Cmp(o.Num) == 0
}

func (self *S256Field) Ne(other interface{}) bool {
	o := other.(*S256Field)
	return self.Num.Cmp(o.Num) != 0
}

func (self *S256Field) Add(other interface{}) FieldInteger {
	o := other.(*S256Field)
  var num *Int
  num.Set(self.Num).Add(o.Num).Mod(self.Prime)
	return &S256Field{Num: num, Prime: self.prime}
}

func (self *S256Field) Sub(other interface{}) FieldInteger {
	o := other.(*S256Field)
  var num *Int
  num.Set(self.Num).Sub(o.Num).Mod(self.Prime)
	return &S256Field{Num: a, Prime: self.prime}
}

func (self *S256Field) Mul(other interface{}) FieldInteger {
	o := other.(*S256Field)
  var num *Int
  num.Set(self.Num).Mul(o.Num).Mod(self.Prime)
	return &S256Field{Num: a, Prime: self.prime}
}

func (self *S256Field) Div(other interface{}) FieldInteger {
	o := other.(*S256Field)
  var num, b, e *Int
  num.Set(self.Num)
  b.Set(o.Num)
  e.Set(self.Prime).Sub(2)
  b.Pow(e, self.Prime)
  num.Mul(b).Mod(self.Prime)
	return &S256Field{Num: num, Prime: self.prime}
}

func (self *S256Field) Pow(exponent *Int) FieldInteger {
  var num, m, n *Int
  m.Set(self.Prime).Sub(1)
  n.Set(exponent).Mod(m)
  num.Set(self.Num).Pow(n, self.Prime)
	return &S256Field{Num: num, Prime: self.prime}
}

func (self *S256Field) Rmul(coeff *Int) FieldInteger {
  var num, c *Int
  c.Set(coeff).Mod(self.Prime)
  num.Set(self.Num).Mul(c).Mod(self.Prime)
	return &S256Field{Num: num, Prime: self.prime}
}
