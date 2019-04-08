package ecc

import (
	"encoding/hex"
	"math/big"
)

const BASE58ALPHABET string = "123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz"

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

func encodeBase58(s string) string {
	count := 0
	chars := []byte(s)
	for _, c := range chars {
		if c == 0 {
			count += 1
		} else {
			break
		}
	}
	var prefix []int = make([]int, count)
	var encoded []int
	var num *big.Int = new(big.Int)
	var mod *big.Int = new(big.Int)
	var b58 *big.Int = big.NewInt(58)
	num.SetBytes(chars)
	for num.Cmp(big.NewInt(0)) > 0 {
		mod.Mod(num, b58)
		num.Div(num, b58)
		encoded = append(encoded, int(mod.Int64()))
	}
	encoded = append(encoded, prefix...)
	alphabet := []byte(BASE58ALPHABET)
	var result []byte = make([]byte, len(encoded))
	for i := len(encoded) - 1; i >= 0; i-- {
		result[len(encoded)-i-1] = alphabet[encoded[i]]
	}
	return string(result)
}

func encodeBase58Checksum(b []byte) string {
	return encodeBase58(string(append(b, hash256(b)[:4]...)))
}
