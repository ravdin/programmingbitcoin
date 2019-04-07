package ecc

import (
	"encoding/hex"
	"math/big"
)

func hexStringToBigInt(str string) *big.Int {
	if len(str)&1 == 1 {
		str = "0" + str
	}
	src := []byte(str)
	dst := make([]byte, hex.DecodedLen(len(str)))
	_, err := hex.Decode(dst, src)
	if err != nil {
		panic(err)
	}
	result := new(big.Int)
	result.SetBytes(dst)
	return result
}
