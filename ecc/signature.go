package ecc

import (
	"bytes"
	"fmt"
	"math/big"
)

// Signature encapsulates a digital signature.
type Signature struct {
	r *big.Int
	s *big.Int
}

// NewSignature initializes a Signature object.
func NewSignature(r *big.Int, s *big.Int) *Signature {
	return &Signature{r: r, s: s}
}

func (sig *Signature) String() string {
	return fmt.Sprintf("Signature(%x,%x)", sig.r, sig.s)
}

// Der serializes the signature.
func (sig *Signature) Der() []byte {
	encode := func(num *big.Int) []byte {
		bin := num.Bytes()
		// if result has a high bit, add a \x00
		if bin[0]&0x80 == 0x80 {
			bin = append(bin, 0)
			copy(bin[1:], bin)
			bin[0] = 0
		}
		result := []byte{2, byte(len(bin))}
		result = append(result, bin...)
		return result
	}
	rbin := encode(sig.r)
	sbin := encode(sig.s)
	result := []byte{0x30, byte(len(rbin) + len(sbin))}
	result = append(result, rbin...)
	result = append(result, sbin...)
	return result
}

// ParseSignature parses a signature in DER format.
func ParseSignature(signatureBin []byte) *Signature {
	reader := bytes.NewReader(signatureBin)
	compound, _ := reader.ReadByte()
	if compound != 0x30 {
		panic("Bad Signature")
	}
	length, _ := reader.ReadByte()
	if int(length+2) != len(signatureBin) {
		panic("Bad Signature Length")
	}
	marker, _ := reader.ReadByte()
	if marker != 0x02 {
		panic("Bad Signature")
	}
	rlength, _ := reader.ReadByte()
	r := new(big.Int)
	buffer := make([]byte, rlength)
	reader.Read(buffer)
	r.SetBytes(buffer)
	marker, _ = reader.ReadByte()
	if marker != 0x02 {
		panic("Bad Signature")
	}
	slength, _ := reader.ReadByte()
	s := new(big.Int)
	buffer = make([]byte, slength)
	reader.Read(buffer)
	s.SetBytes(buffer)
	if len(signatureBin) != 6+int(rlength+slength) {
		panic("Signature too long")
	}
	return NewSignature(r, s)
}
