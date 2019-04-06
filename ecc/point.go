package ecc

import (
	"errors"
	"fmt"
	"math/big"
)

var G *Point
var A *big.Int
var B *big.Int
var N *big.Int
var P *big.Int

func init() {
	A = big.NewInt(0)
	B = big.NewInt(7)
	N = new(big.Int)
	N.SetString("fffffffffffffffffffffffffffffffebaaedce6af48a03bbfd25e8cd0364141", 16)

	// 2^256 - 2^32 - 2^9 - 2^8 - 2^7 - 2^6 - 2^4 - 1
	P = new(big.Int)
	P.SetString("fffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f", 16)

	Gx := "79be667ef9dcbbac55a06295ce870b07029bfcdb2dce28d959f2815b16f81798"
	Gy := "483ada7726a3c4655da4fbfc0e1108a8fd17b448a68554199c47d08ffb10d4b8"
	var x, y big.Int
	x.SetString(Gx, 16)
	y.SetString(Gy, 16)
	G, _ = NewS256Point(&x, &y)
}

type Point struct {
	X FieldInteger
	Y FieldInteger
	A FieldInteger
	B FieldInteger
}

func NewPoint(x_arg interface{}, y_arg interface{}, a_arg interface{}, b_arg interface{}, option FieldIntegerConverter) (*Point, error) {
	a := option(a_arg)
	b := option(b_arg)
	if x_arg == nil && y_arg == nil {
		return &Point{X: nil, Y: nil, A: a, B: b}, nil
	}
	x := option(x_arg)
	y := option(y_arg)
	if !y.Pow(big.NewInt(2)).Eq(x.Pow(big.NewInt(3)).Add(x.Mul(a)).Add(b)) {
		return nil, errors.New(fmt.Sprintf("(%d, %d) is not on the curve", x, y))
	}
	return &Point{X: x, Y: y, A: a, B: b}, nil
}

func NewS256Point(x *big.Int, y *big.Int) (*Point, error) {
	withBigInt := func(bigint interface{}) FieldInteger {
		return NewS256Field(bigint.(*big.Int), P)
	}
	result, err := NewPoint(x, y, A, B, withBigInt)
	return result, err
}

func (self *Point) String() string {
	if self.X == nil {
		return "Point(infinity)"
	} else {
		return fmt.Sprintf("Point(%v,%v)_%v_%v", self.X, self.Y, self.A, self.B)
	}
}

func (self *Point) Eq(other *Point) bool {
	if self.X == nil {
		return other.X == nil
	}
	return self.X.Eq(other.X) &&
		self.Y.Eq(other.Y) &&
		self.A.Eq(other.A) &&
		self.B.Eq(other.B)
}

func (self *Point) Ne(other *Point) bool {
	return !self.Eq(other)
}

func (self *Point) Add(other *Point) *Point {
	if self.A.Ne(other.A) || self.B.Ne(other.B) {
		panic(fmt.Sprintf("Points %v, %v are not on the same curve", self, other))
	}

	if self.X == nil {
		return other
	}

	if other.X == nil {
		return self
	}

	// Case 1: self.x == other.x, self.y != other.y
	// Result is point at infinity
	if self.X.Eq(other.X) && self.Y.Ne(other.Y) {
		return &Point{X: nil, Y: nil, A: self.A, B: self.B}
	}

	// Case 2: self.x ≠ other.x
	// Formula (x3,y3)==(x1,y1)+(x2,y2)
	// s=(y2-y1)/(x2-x1)
	// x3=s**2-x1-x2
	// y3=s*(x1-x3)-y1
	if self.X.Ne(other.X) {
		s := (other.Y.Sub(self.Y)).Div(other.X.Sub(self.X))
		x := s.Pow(big.NewInt(2)).Sub(self.X).Sub(other.X)
		y := s.Mul(self.X.Sub(x)).Sub(self.Y)
		return &Point{X: x, Y: y, A: self.A, B: self.B}
	}

	// Case 4: if we are tangent to the vertical line,
	// we return the point at infinity
	// note instead of figuring out what 0 is for each type
	// we just use 0 * self.x
	if self.Eq(other) && self.Y.Eq(self.X.Rmul(big.NewInt(0))) {
		return &Point{X: nil, Y: nil, A: self.A, B: self.B}
	}

	// Case 3: self == other
	// Formula (x3,y3)=(x1,y1)+(x1,y1)
	// s=(3*x1**2+a)/(2*y1)
	// x3=s**2-2*x1
	// y3=s*(x1-x3)-y1
	s := self.X.Pow(big.NewInt(2)).Rmul(big.NewInt(3)).Add(self.A).Div(self.Y.Rmul(big.NewInt(2)))
	x := s.Pow(big.NewInt(2)).Sub(self.X.Rmul(big.NewInt(2)))
	y := s.Mul(self.X.Sub(x)).Sub(self.Y)
	return &Point{X: x, Y: y, A: self.A, B: self.B}
}

func (self *Point) Rmul(coefficient *big.Int) *Point {
	var coef *big.Int = new(big.Int)
	var z *big.Int = new(big.Int)
	coef.Set(coefficient)
	var current *Point = &Point{X: self.X, Y: self.Y, A: self.A, B: self.B}
	var result *Point = &Point{X: nil, Y: nil, A: self.A, B: self.B}
	for coef.Cmp(big.NewInt(0)) != 0 {
		z.And(coef, big.NewInt(1))
		if z.Cmp(big.NewInt(1)) == 0 {
			result = result.Add(current)
		}
		current = current.Add(current)
		coef.Rsh(coef, 1)
	}
	return result
}
