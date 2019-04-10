package util

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"math/big"
)

const BASE58ALPHABET string = "123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz"

func HexStringToBytes(str string) []byte {
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

func HexStringToBigInt(str string) *big.Int {
	result := new(big.Int)
	result.SetBytes(HexStringToBytes(str))
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
		num.QuoRem(num, b58, mod)
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

func EncodeBase58Checksum(b []byte) string {
	return encodeBase58(string(append(b, Hash256(b)[:4]...)))
}

func IntToBytes(num *big.Int, size int) []byte {
	var result []byte = make([]byte, size)
	var raw = num.Bytes()
	copy(result[size-len(raw):], raw)
	return result
}

func ReadVarInt(r bytes.Reader) int64 {
	b, _ := r.ReadByte()
	var bufsize int
	switch (b) {
	case 0xfd:
		bufsize = 2
		break
	case 0xfe:
		bufsize = 4
		break
	case 0xff:
		bufsize = 8
		break
	default:
		return int64(b)
	}
	var buffer []byte = make([]byte, bufsize)
	r.Read(buffer)
	return LittleEndianToInt64(buffer)
}

func EncodeVarInt(i int) []byte {
	if i < 0xfd {
		return []byte{byte(i)}
	}
	if i < 0x10000 {
		result := make([]byte, 3)
		copy(result[1:], Int16ToLittleEndian(int16(i)))
		result[0] = 0xfd
		return result
	}
	if i < 0x100000000 {
		result := make([]byte, 5)
		copy(result[1:], Int32ToLittleEndian(int32(i)))
		result[0] = 0xfe
		return result
	}

	result := make([]byte, 9)
	copy(result[1:], Int64ToLittleEndian(int64(i)))
	result[0] = 0xff
	return result
}

func LittleEndianToInt16(b []byte) int16 {
	if len(b) > 2 {
		panic ("Value is too large!")
	}
	var result int16
	buf := bytes.NewReader(b)
	err := binary.Read(buf, binary.LittleEndian, &result)
	if err != nil {
		panic(err)
	}
	return result
}

func LittleEndianToInt32(b []byte) int32 {
	if len(b) > 4 {
		panic ("Value is too large!")
	}
	var result int32
	buf := bytes.NewReader(b)
	err := binary.Read(buf, binary.LittleEndian, &result)
	if err != nil {
		panic(err)
	}
	return result
}

func LittleEndianToInt64(b []byte) int64 {
	if len(b) > 8 {
		panic ("Value is too large!")
	}
	var result int64
	buf := bytes.NewReader(b)
	err := binary.Read(buf, binary.LittleEndian, &result)
	if err != nil {
		panic(err)
	}
	return result
}

func Int16ToLittleEndian(num int16) []byte {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.LittleEndian, &num)
	if err != nil {
		panic(err)
	}
	return buf.Bytes()
}

func Int32ToLittleEndian(num int32) []byte {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.LittleEndian, &num)
	if err != nil {
		panic(err)
	}
	return buf.Bytes()
}

func Int64ToLittleEndian(num int64) []byte {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.LittleEndian, &num)
	if err != nil {
		panic(err)
	}
	return buf.Bytes()
}
