package ecc

import (
	"bytes"
	"fmt"
	"math/big"
)

type Signature struct {
	R *big.Int
	S *big.Int
}

func NewSignature(r *big.Int, s *big.Int) *Signature {
	return &Signature{R: r, S: s}
}

func (sig *Signature) String() string {
	return fmt.Sprintf("Signature(%v,%v)", sig.R, sig.S)
}

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
	rbin := encode(sig.R)
	sbin := encode(sig.S)
	result := []byte{0x30, byte(len(rbin) + len(sbin))}
	result = append(result, rbin...)
	result = append(result, sbin...)
	return result
}

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
	var r *big.Int = new(big.Int)
	buffer := make([]byte, rlength)
	reader.Read(buffer)
	r.SetBytes(buffer)
	marker, _ = reader.ReadByte()
	if marker != 0x02 {
		panic("Bad Signature")
	}
	slength, _ := reader.ReadByte()
	var s *big.Int = new(big.Int)
	buffer = make([]byte, slength)
	s.SetBytes(buffer)
	if len(signatureBin) != 6+int(rlength+slength) {
		panic("Signature too long")
	}
	return NewSignature(r, s)
}
