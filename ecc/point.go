package ecc

import (
	"errors"
	"fmt"
	"math/big"
)

type FieldInteger interface {
	Add(other interface{}) FieldInteger
	Sub(other interface{}) FieldInteger
	Mul(other interface{}) FieldInteger
	Div(other interface{}) FieldInteger
	Pow(exponent *big.Int) FieldInteger
	Rmul(coeff *big.Int) FieldInteger
	Eq(other interface{}) bool
	Ne(other interface{}) bool
}

type Point struct {
	X FieldInteger
	Y FieldInteger
	A FieldInteger
	B FieldInteger
}

func NewPoint(x FieldInteger, y FieldInteger, a FieldInteger, b FieldInteger) (*Point, error) {
	if x == nil && y == nil {
		return &Point{X: nil, Y: nil, A: a, B: b}, nil
	}
	if !y.Pow(big.NewInt(2)).Eq(x.Pow(big.NewInt(3)).Add(x.Mul(a)).Add(b)) {
		return nil, errors.New(fmt.Sprintf("(%d, %d) is not on the curve", x, y))
	}
	return &Point{X: x, Y: y, A: a, B: b}, nil
}

func (self *Point) String() string {
	if self.X == nil {
		return "Point(infinity)"
	} else {
		return fmt.Sprintf("Point(%v,%v)_%v_%v", self.X, self.Y, self.A, self.B)
	}
}

func (self *Point) Eq(other *Point) bool {
	return self.X.Eq(other.X) &&
		self.Y.Eq(other.Y) &&
		self.A.Eq(other.A) &&
		self.B.Eq(other.B)
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
