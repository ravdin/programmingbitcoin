package ecc

import (
	"encoding/hex"
	"math/big"
)

func hexStringToBytes(str string) []byte {
	if len(str)&1 == 1 {
		str = "0" + str
	}
	src := []byte(str)
	dst := make([]byte, hex.DecodedLen(len(str)))
	_, err := hex.Decode(dst, src)
	if err != nil {
		panic(err)
	}
	return dst
}

func hexStringToBigInt(str string) *big.Int {
	result := new(big.Int)
	result.SetBytes(hexStringToBytes(str))
	return result
}
