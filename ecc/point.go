package ecc

import (
	"fmt"
	"math/big"
)

// Point represents a point in an elliptic curve.
type Point struct {
	X FieldInteger
	Y FieldInteger
	A FieldInteger
	B FieldInteger
}

// NewPoint returns a Point.
func NewPoint(x FieldInteger, y FieldInteger, a FieldInteger, b FieldInteger) (*Point, error) {
	if x == nil && y == nil {
		return &Point{X: nil, Y: nil, A: a, B: b}, nil
	}
	// Verify that y^2 = x^3 + ax + b
	var left, right FieldInteger = y.Copy(), x.Copy()
	left.Pow(left, big.NewInt(2))
	right.Pow(right, big.NewInt(3)).Add(right, x.Copy().Mul(x, a)).Add(right, b)
	if !left.Eq(right) {
		return nil, fmt.Errorf("(%d, %d) is not on the curve", x, y)
	}
	return &Point{X: x, Y: y, A: a, B: b}, nil
}

func (p *Point) String() string {
	if p.X == nil {
		return "Point(infinity)"
	}
	return fmt.Sprintf("Point(%v,%v)_%v_%v", p.X, p.Y, p.A, p.B)
}

// Eq returns true if the points are equal.
func (p *Point) Eq(other *Point) bool {
	if p.X == nil {
		return other.X == nil
	}
	return p.X.Eq(other.X) &&
		p.Y.Eq(other.Y) &&
		p.A.Eq(other.A) &&
		p.B.Eq(other.B)
}

// Ne returns true if the points are not equal.
func (p *Point) Ne(other *Point) bool {
	return !p.Eq(other)
}

// Add p1 + p2 and return p.
func (p *Point) Add(p1, p2 *Point) *Point {
	if p1.A.Ne(p2.A) || p1.B.Ne(p2.B) {
		panic(fmt.Sprintf("Points %v, %v are not on the same curve", p1, p2))
	}

	a, b := p1.A, p1.B
	x1, y1, x2, y2 := p1.X, p1.Y, p2.X, p2.Y
	two := big.NewInt(2)

	if p1.X == nil {
		*p = Point{X: p2.X, Y: p2.Y, A: a, B: b}
		return p
	}

	if p2.X == nil {
		*p = Point{X: p1.X, Y: p1.Y, A: a, B: b}
		return p
	}

	// Case 1: p1.x == p2.x, p1.y != p2.y
	// Result is point at infinity
	if p1.X.Eq(p2.X) && p1.Y.Ne(p2.Y) {
		*p = Point{X: nil, Y: nil, A: a, B: b}
		return p
	}

	// Case 2: p1.x ≠ p2.x
	// Formula (x3,y3)==(x1,y1)+(x2,y2)
	// s=(y2-y1)/(x2-x1)
	// x3=s**2-x1-x2
	// y3=s*(x1-x3)-y1
	if x1.Ne(x2) {
		s := y2.Copy()
		s.Sub(s, y1)
		tmp := x2.Copy()
		tmp.Sub(tmp, x1)
		s.Div(s, tmp)
		x3 := s.Copy()
		x3.Pow(x3, two).Sub(x3, x1).Sub(x3, x2)
		y3 := tmp.Sub(x1, x3)
		y3.Mul(y3, s).Sub(y3, y1)
		*p = Point{X: x3, Y: y3, A: a, B: b}
		return p
	}

	// Case 4: if we are tangent to the vertical line,
	// we return the point at infinity
	// note instead of figuring out what 0 is for each type
	// we just use 0 * p.x
	zero := x1.Copy()
	zero.Cmul(zero, big.NewInt(0))
	if p1.Eq(p2) && p1.Y.Eq(zero) {
		*p = Point{X: nil, Y: nil, A: a, B: b}
		return p
	}

	// Case 3: p1 == p2
	// Formula (x3,y3)=(x1,y1)+(x1,y1)
	// s=(3*x1**2+a)/(2*y1)
	// x3=s**2-2*x1
	// y3=s*(x1-x3)-y1
	s := x1.Copy()
	s.Pow(s, two).Cmul(s, big.NewInt(3)).Add(s, a)
	tmp := y1.Copy()
	tmp.Cmul(tmp, two)
	s.Div(s, tmp)
	tmp.Cmul(x1, two)
	x3 := s.Copy()
	x3.Pow(x3, two).Sub(x3, tmp)
	y3 := s.Copy()
	tmp.Sub(x1, x3)
	y3.Mul(s, tmp).Sub(y3, y1)
	*p = Point{X: x3, Y: y3, A: a, B: b}
	return p
}

// Cmul sets p to c * r and return p.
func (p *Point) Cmul(r *Point, coefficient *big.Int) *Point {
	coef := new(big.Int)
	coef.Set(coefficient)
	current := &Point{X: r.X, Y: r.Y, A: r.A, B: r.B}
	result := &Point{X: nil, Y: nil, A: r.A, B: r.B}
	for coef.Sign() > 0 {
		if coef.Bit(0) == 1 {
			result.Add(result, current)
		}
		current.Add(current, current)
		coef.Rsh(coef, 1)
	}
	*p = *result
	return p
}
