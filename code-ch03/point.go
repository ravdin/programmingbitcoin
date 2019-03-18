package code_ch03

import (
	"errors"
	"fmt"
)

type FieldInteger interface {
	Add(other interface{}) FieldInteger
	Sub(other interface{}) FieldInteger
	Mul(other interface{}) FieldInteger
	Div(other interface{}) FieldInteger
	Pow(exponent int64) FieldInteger
	Rmul(coeff int64) FieldInteger
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
	if !y.Pow(2).Eq(x.Pow(3).Add(x.Mul(a)).Add(b)) {
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

	panic("Not implemented")
}
