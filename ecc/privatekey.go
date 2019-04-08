package ecc

import (
	//"fmt"
	"crypto/hmac"
	"crypto/sha256"
	"math/big"
)

type PrivateKey struct {
	secret *big.Int
	point  *S256Point
}

func NewPrivateKey(secret *big.Int) *PrivateKey {
	return &PrivateKey{secret: secret, point: G.Rmul(secret)}
}

func (self *PrivateKey) Sign(z *big.Int) *Signature {
	k := self.deterministicK(z)
	// r is the x coordinate of the resulting point k*G
	r := G.Rmul(k).X.Num
	// remember 1/k = pow(k, N-2, N)
	e := new(big.Int)
	e.Sub(N, big.NewInt(2))
	k_inv := new(big.Int)
	k_inv.Exp(k, e, N)
	// s = (z+r*secret) / k
	s := new(big.Int)
	s.Mul(r, self.secret)
	s.Add(s, z)
	s.Mul(s, k_inv)
	s.Mod(s, N)
	tmp := new(big.Int)
	tmp.Mul(s, big.NewInt(2))
	if tmp.Cmp(N) > 0 {
		s.Sub(N, s)
	}
	// return an instance of Signature:
	// Signature(r, s)
	return NewSignature(r, s)
}

func (self *PrivateKey) Wif(compressed bool, testnet bool) string {
	var secretBytes = make([]byte, 33)
	copy(secretBytes[1:], intToBytes(self.secret, 32))
	if testnet {
		secretBytes[0] = 0xef
	} else {
		secretBytes[0] = 0x80
	}
	if compressed {
		secretBytes = append(secretBytes, 1)
	}
	return encodeBase58Checksum(secretBytes)
}

func (self *PrivateKey) deterministicK(z *big.Int) *big.Int {
	var ztmp *big.Int = new(big.Int)
	var k []byte = make([]byte, 32)
	var v []byte = make([]byte, 32)
	for i := 0; i < 32; i++ {
		k[i] = 0
		v[i] = 1
	}
	ztmp.Set(z)
	if ztmp.Cmp(N) > 0 {
		ztmp.Sub(ztmp, N)
	}
	zBytes := intToBytes(ztmp, 32)
	secretBytes := intToBytes(self.secret, 32)
	mac := hmac.New(sha256.New, k)
	mac.Write(v)
	mac.Write([]byte{0})
	mac.Write(secretBytes)
	mac.Write(zBytes)
	k = mac.Sum(nil)
	mac = hmac.New(sha256.New, k)
	mac.Write(v)
	v = mac.Sum(nil)
	mac.Reset()
	mac.Write(v)
	mac.Write([]byte{1})
	mac.Write(secretBytes)
	mac.Write(zBytes)
	k = mac.Sum(nil)
	mac = hmac.New(sha256.New, k)
	mac.Write(v)
	v = mac.Sum(nil)
	mac.Reset()
	var candidate *big.Int = new(big.Int)
	for true {
		mac.Write(v)
		v = mac.Sum(nil)
		candidate.SetBytes(v)
		if candidate.Cmp(big.NewInt(1)) >= 0 && candidate.Cmp(N) < 0 {
			break
		}
		mac.Reset()
		mac.Write(v)
		mac.Write([]byte{0})
		k = mac.Sum(nil)
		mac = hmac.New(sha256.New, k)
		mac.Write(v)
		v = mac.Sum(nil)
		mac.Reset()
	}
	return candidate
}
