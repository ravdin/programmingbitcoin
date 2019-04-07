package ecc

import (
	"fmt"
	"math/big"
)

var G *S256Point
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

type S256Point struct {
	X *S256Field
	Y *S256Field
}

func NewS256Point(x *big.Int, y *big.Int) (*S256Point, error) {
	withBigInt := func(bigint interface{}) FieldInteger {
		return NewS256Field(bigint.(*big.Int), P)
	}
	p, err := NewPoint(x, y, A, B, withBigInt)
	if err != nil {
		return nil, err
	}
	return &S256Point{X: p.X.(*S256Field), Y: p.Y.(*S256Field)}, nil
}

func (self *S256Point) Point() *Point {
	var s256FieldConverter = func(num interface{}) FieldInteger {
		switch num.(type) {
		case *big.Int:
			return NewS256Field(num.(*big.Int), P)
		case *S256Field:
			return num.(*S256Field)
		default:
			panic("Unsupported type!")
		}
	}
	if result, err := NewPoint(self.X, self.Y, A, B, s256FieldConverter); err == nil {
		return result
	}
	panic("Error casting to Point!")
}

func (self *S256Point) String() string {
	if self.X == nil {
		return "Point(infinity)"
	} else {
		return fmt.Sprintf("Point(%v,%v)", self.X.Num, self.Y.Num)
	}
}

func (self *S256Point) Eq(other *S256Point) bool {
	if self.X == nil {
		return other.X == nil
	}
	return self.X.Eq(other.X) && self.Y.Eq(other.Y)
}

func (self *S256Point) Ne(other *S256Point) bool {
	return !self.Eq(other)
}

func (self *S256Point) Add(other *S256Point) *S256Point {
	p1 := self.Point()
	p2 := other.Point()
	result := p1.Add(p2)
	return &S256Point{X: result.X.(*S256Field), Y: result.Y.(*S256Field)}
}

func (self *S256Point) Rmul(coefficient *big.Int) *S256Point {
	var coef *big.Int = new(big.Int)
	coef.Mod(coefficient, N)
	result := self.Point().Rmul(coef)
	if result.X == nil && result.Y == nil {
		return &S256Point{X: nil, Y: nil}
	}
	return &S256Point{X: result.X.(*S256Field), Y: result.Y.(*S256Field)}
}
