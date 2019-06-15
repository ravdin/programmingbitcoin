package ecc

import (
	"fmt"
	"math/big"

	"github.com/ravdin/programmingbitcoin/util"
)

var (
	G *S256Point
	A *s256Field
	B *s256Field
	N *big.Int
	P *big.Int
)

func init() {
	N = util.HexStringToBigInt("fffffffffffffffffffffffffffffffebaaedce6af48a03bbfd25e8cd0364141")

	// 2^256 - 2^32 - 2^9 - 2^8 - 2^7 - 2^6 - 2^4 - 1
	P = util.HexStringToBigInt("fffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f")
	A = newS256Field(big.NewInt(0), P)
	B = newS256Field(big.NewInt(7), P)

	x := util.HexStringToBigInt("79be667ef9dcbbac55a06295ce870b07029bfcdb2dce28d959f2815b16f81798")
	y := util.HexStringToBigInt("483ada7726a3c4655da4fbfc0e1108a8fd17b448a68554199c47d08ffb10d4b8")
	G, _ = NewS256Point(x, y)
}

type S256Point struct {
	X *s256Field
	Y *s256Field
}

func NewS256Point(x *big.Int, y *big.Int) (*S256Point, error) {
	p, err := NewPoint(newS256Field(x, P), newS256Field(y, P), A, B)
	if err != nil {
		return nil, err
	}
	return &S256Point{X: p.X.(*s256Field), Y: p.Y.(*s256Field)}, nil
}

// ParseS256Point returns a Point object from a SEC binary (not hex)
func ParseS256Point(secBin []byte) *S256Point {
	if secBin[0] == 4 {
		var x, y *big.Int = new(big.Int), new(big.Int)
		x.SetBytes(secBin[1:33])
		y.SetBytes(secBin[33:65])
		result, _ := NewS256Point(x, y)
		return result
	}
	isEven := secBin[0] == 2
	xval := new(big.Int)
	xval.SetBytes(secBin[1:])
	x := newS256Field(xval, P)
	// right side of the equation y^2 = x^3 + 7
	alpha := new(s256Field)
	alpha.Pow(x, big.NewInt(3))
	alpha.Add(alpha, B)
	// solve for left side
	beta := alpha.Sqrt()
	var evenBeta, oddBeta *s256Field
	betaOffset := new(big.Int)
	betaOffset.Sub(P, beta.Num)
	if beta.Num.Bit(0) == 0 {
		evenBeta = beta
		oddBeta = newS256Field(betaOffset, P)
	} else {
		evenBeta = newS256Field(betaOffset, P)
		oddBeta = beta
	}
	if isEven {
		return &S256Point{X: x, Y: evenBeta}
	}
	return &S256Point{X: x, Y: oddBeta}
}

func (p *S256Point) point() *Point {
	if result, err := NewPoint(p.X, p.Y, A, B); err == nil {
		return result
	}
	panic("Error casting to Point!")
}

func (p *S256Point) String() string {
	if p.X == nil {
		return "Point(infinity)"
	} else {
		return fmt.Sprintf("Point(%v,%v)", p.X.Num, p.Y.Num)
	}
}

func (p *S256Point) Eq(other *S256Point) bool {
	if p.X == nil {
		return other.X == nil
	}
	return p.X.Eq(other.X) && p.Y.Eq(other.Y)
}

func (p *S256Point) Ne(other *S256Point) bool {
	return !p.Eq(other)
}

// Set z to p1 + p2 and return z.
func (p *S256Point) Add(p1, p2 *S256Point) *S256Point {
	result := new(Point)
	result.Add(p1.point(), p2.point())
	if result.X == nil && result.Y == nil {
		*p = S256Point{X: nil, Y: nil}
	} else {
		*p = S256Point{X: result.X.(*s256Field), Y: result.Y.(*s256Field)}
	}
	return p
}

// Set z to c * r and return p.
func (p *S256Point) Cmul(r *S256Point, coefficient *big.Int) *S256Point {
	coef := new(big.Int)
	coef.Mod(coefficient, N)
	result := new(Point)
	result.Cmul(r.point(), coef)
	if result.X == nil && result.Y == nil {
		*p = S256Point{X: nil, Y: nil}
	} else {
		*p = S256Point{X: result.X.(*s256Field), Y: result.Y.(*s256Field)}
	}
	return p
}

func (p *S256Point) Verify(z *big.Int, sig *Signature) bool {
	// By Fermat's Little Theorem, 1/s = pow(s, N-2, N)
	sInv := new(big.Int)
	e := new(big.Int)
	e.Sub(N, big.NewInt(2))
	sInv.Exp(sig.s, e, N)
	// u = z / s
	u := new(big.Int)
	u.Mul(z, sInv).Mod(u, N)
	// v = r / s
	v := new(big.Int)
	v.Mul(sig.r, sInv).Mod(v, N)
	// u*G + v*P should have as the x coordinate, r
	total := new(S256Point)
	total.Cmul(G, u)
	total.Add(total, new(S256Point).Cmul(p, v))
	return total.X.Num.Cmp(sig.r) == 0
}

// returns the binary version of the SEC format
func (p *S256Point) Sec(compressed bool) []byte {
	x := util.IntToBytes(p.X.Num, 32)
	y := util.IntToBytes(p.Y.Num, 32)
	var result []byte
	if compressed {
		result = make([]byte, 33)
		copy(result[1:], x)
		if y[31]%2 == 0 {
			result[0] = 2
		} else {
			result[0] = 3
		}
	} else {
		result = make([]byte, 65)
		copy(result[1:], x)
		copy(result[33:], y)
		result[0] = 4
	}
	return result
}

func (p *S256Point) Hash160(compressed bool) []byte {
	return util.Hash160(p.Sec(compressed))
}

func (p *S256Point) Address(compressed bool, testnet bool) string {
	h160 := p.Hash160(compressed)
	var prefix byte
	if testnet {
		prefix = 0x6f
	} else {
		prefix = 0
	}
	withPrefix := make([]byte, len(h160)+1)
	withPrefix[0] = prefix
	copy(withPrefix[1:], h160)
	return util.EncodeBase58Checksum(withPrefix)
}
