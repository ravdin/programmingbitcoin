package ecc

import (
	"crypto/hmac"
	"crypto/sha256"
	"fmt"
	"math/big"

	"github.com/ravdin/programmingbitcoin/util"
)

// PrivateKey represents a private key with a secret.
type PrivateKey struct {
	secret *big.Int
	Point  *S256Point
}

// NewPrivateKey returns a PrivateKey instance.
func NewPrivateKey(secret *big.Int) *PrivateKey {
	return &PrivateKey{secret: secret, Point: new(S256Point).Cmul(G, secret)}
}

// Hex returns the private key in hex format.
func Hex(pk *PrivateKey) string {
	return fmt.Sprintf("%x", pk.secret.Bytes())
}

// Sign returns a Signature instance.
func (pk *PrivateKey) Sign(z *big.Int) *Signature {
	k := pk.deterministicK(z)
	// r is the x coordinate of the resulting point k*G
	r := new(S256Point).Cmul(G, k).X.Num
	// remember 1/k = pow(k, N-2, N)
	e := new(big.Int)
	e.Sub(N, big.NewInt(2))
	kInv := new(big.Int)
	kInv.Exp(k, e, N)
	// s = (z+r*secret) / k
	s := new(big.Int)
	s.Mul(r, pk.secret)
	s.Add(s, z)
	s.Mul(s, kInv)
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

// Wif converts the secret from integer to a 32-bytes in big endian
func (pk *PrivateKey) Wif(compressed bool, testnet bool) string {
	var secretBytes = make([]byte, 33)
	copy(secretBytes[1:], util.IntToBytes(pk.secret, 32))
	if testnet {
		secretBytes[0] = 0xef
	} else {
		secretBytes[0] = 0x80
	}
	if compressed {
		secretBytes = append(secretBytes, 1)
	}
	return util.EncodeBase58Checksum(secretBytes)
}

func (pk *PrivateKey) deterministicK(z *big.Int) *big.Int {
	k := make([]byte, 32)
	v := make([]byte, 32)
	for i := 0; i < 32; i++ {
		k[i] = 0
		v[i] = 1
	}
	ztmp := new(big.Int)
	ztmp.Mod(z, N)
	zBytes := util.IntToBytes(ztmp, 32)
	secretBytes := util.IntToBytes(pk.secret, 32)
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
	candidate := new(big.Int)
	for true {
		mac.Write(v)
		v = mac.Sum(nil)
		candidate.SetBytes(v)
		if candidate.Sign() > 0 && candidate.Cmp(N) < 0 {
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
