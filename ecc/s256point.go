package ecc

import (
	"fmt"
	"github.com/ravdin/programmingbitcoin/util"
	"math/big"
)

var (
	G *S256Point
	A *S256Field
	B *S256Field
	N *big.Int
	P *big.Int
)

func init() {
	N = util.HexStringToBigInt("fffffffffffffffffffffffffffffffebaaedce6af48a03bbfd25e8cd0364141")

	// 2^256 - 2^32 - 2^9 - 2^8 - 2^7 - 2^6 - 2^4 - 1
	P = util.HexStringToBigInt("fffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f")
	A = NewS256Field(big.NewInt(0), P)
	B = NewS256Field(big.NewInt(7), P)

	x := util.HexStringToBigInt("79be667ef9dcbbac55a06295ce870b07029bfcdb2dce28d959f2815b16f81798")
	y := util.HexStringToBigInt("483ada7726a3c4655da4fbfc0e1108a8fd17b448a68554199c47d08ffb10d4b8")
	G, _ = NewS256Point(x, y)
}

type S256Point struct {
	X *S256Field
	Y *S256Field
}

func NewS256Point(x *big.Int, y *big.Int) (*S256Point, error) {
	p, err := NewPoint(NewS256Field(x, P), NewS256Field(y, P), A, B)
	if err != nil {
		return nil, err
	}
	return &S256Point{X: p.X.(*S256Field), Y: p.Y.(*S256Field)}, nil
}

// returns a Point object from a SEC binary (not hex)
func ParseS256Point(secBin []byte) *S256Point {
	if secBin[0] == 4 {
		var x, y *big.Int = new(big.Int), new(big.Int)
		x.SetBytes(secBin[1:33])
		y.SetBytes(secBin[33:65])
		result, _ := NewS256Point(x, y)
		return result
	}
	isEven := secBin[0] == 2
	var xval *big.Int = new(big.Int)
	xval.SetBytes(secBin[1:])
	x := NewS256Field(xval, P)
	// right side of the equation y^2 = x^3 + 7
	alpha := new(S256Field)
	alpha.Pow(x, big.NewInt(3))
	alpha.Add(alpha, B)
	// solve for left side
	beta := alpha.Sqrt()
	var even_beta, odd_beta *S256Field
	var betaOffset *big.Int = new(big.Int)
	betaOffset.Sub(P, beta.Num)
	if beta.Num.Bit(0) == 0 {
		even_beta = beta
		odd_beta = NewS256Field(betaOffset, P)
	} else {
		even_beta = NewS256Field(betaOffset, P)
		odd_beta = beta
	}
	if isEven {
		return &S256Point{X: x, Y: even_beta}
	} else {
		return &S256Point{X: x, Y: odd_beta}
	}
}

func (self *S256Point) Point() *Point {
	if result, err := NewPoint(self.X, self.Y, A, B); err == nil {
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

// Set z to p1 + p2 and return z.
func (z *S256Point) Add(p1, p2 *S256Point) *S256Point {
	result := new(Point)
	result.Add(p1.Point(), p2.Point())
	if result.X == nil && result.Y == nil {
		*z = S256Point{X: nil, Y: nil}
	} else {
		*z = S256Point{X: result.X.(*S256Field), Y: result.Y.(*S256Field)}
	}
	return z
}

// Set z to c * p and return z.
func (z *S256Point) Cmul(p *S256Point, coefficient *big.Int) *S256Point {
	coef := new(big.Int)
	coef.Mod(coefficient, N)
	result := new(Point)
	result.Cmul(p.Point(), coef)
	if result.X == nil && result.Y == nil {
		*z = S256Point{X: nil, Y: nil}
	} else {
		*z = S256Point{X: result.X.(*S256Field), Y: result.Y.(*S256Field)}
	}
	return z
}

func (self *S256Point) Verify(z *big.Int, sig *Signature) bool {
	// By Fermat's Little Theorem, 1/s = pow(s, N-2, N)
	s_inv := new(big.Int)
	e := new(big.Int)
	e.Sub(N, big.NewInt(2))
	s_inv.Exp(sig.S, e, N)
	// u = z / s
	u := new(big.Int)
	u.Mul(z, s_inv).Mod(u, N)
	// v = r / s
	v := new(big.Int)
	v.Mul(sig.R, s_inv).Mod(v, N)
	// u*G + v*P should have as the x coordinate, r
	total := new(S256Point)
	total.Cmul(G, u)
	total.Add(total, new(S256Point).Cmul(self, v))
	return total.X.Num.Cmp(sig.R) == 0
}

// returns the binary version of the SEC format
func (self *S256Point) Sec(compressed bool) []byte {
	x := util.IntToBytes(self.X.Num, 32)
	y := util.IntToBytes(self.Y.Num, 32)
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

func (self *S256Point) Hash160(compressed bool) []byte {
	return util.Hash160(self.Sec(compressed))
}

func (self *S256Point) Address(compressed bool, testnet bool) string {
	h160 := self.Hash160(compressed)
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
