package util

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"math/big"
)

const (
	SIGHASH_ALL    uint32 = 1
	SIGHASH_NONE   uint32 = 2
	SIGHASH_SINGLE uint32 = 3
	BASE58ALPHABET string = "123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz"
)

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

// Reverse a byte array in place and return the result.
func ReverseByteArray(arr []byte) []byte {
	length := len(arr)
	for i := 0; i < length/2; i++ {
		arr[i], arr[length-i-1] = arr[length-i-1], arr[i]
	}
	return arr
}

func encodeBase58(s string) string {
	// Determine how many 0 bytes the input starts with
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
	for num.Sign() > 0 {
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

func DecodeBase58(encoded string) []byte {
	num := big.NewInt(0)
	b58 := big.NewInt(58)
	alphabet := []byte(BASE58ALPHABET)
	chars := []byte(encoded)
	for _, c := range chars {
		num.Mul(num, b58)
		num.Add(num, big.NewInt(int64(bytes.IndexByte(alphabet, c))))
	}
	combined := num.Bytes()
	length := len(combined)
	checksum := combined[length-4:]
	if !bytes.Equal(Hash256(combined[:length-4])[:4], checksum) {
		panic(fmt.Sprintf("Bad address: %v %v", checksum, Hash256(combined[:length-4])[:4]))
	}
	return combined[1 : length-4]
}

func IntToBytes(num *big.Int, size int) []byte {
	var result []byte = make([]byte, size)
	var raw = num.Bytes()
	copy(result[size-len(raw):], raw)
	return result
}

func ReadVarInt(r *bytes.Reader) int {
	b, err := r.ReadByte()
	if err != nil {
		panic(err)
	}
	var bufsize int
	switch b {
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
		return int(b)
	}
	var buffer []byte = make([]byte, bufsize)
	r.Read(buffer)
	return int(LittleEndianToInt64(buffer))
}

func EncodeVarInt(i int) []byte {
	if i < 0xfd {
		return []byte{byte(i)}
	}
	if i < 0x10000 {
		result := make([]byte, 3)
		copy(result[1:], Int16ToLittleEndian(uint16(i)))
		result[0] = 0xfd
		return result
	}
	if i < 0x100000000 {
		result := make([]byte, 5)
		copy(result[1:], Int32ToLittleEndian(uint32(i)))
		result[0] = 0xfe
		return result
	}

	result := make([]byte, 9)
	copy(result[1:], Int64ToLittleEndian(uint64(i)))
	result[0] = 0xff
	return result
}

func LittleEndianToByte(b byte) byte {
	var result byte
	buf := bytes.NewReader([]byte{b})
	err := binary.Read(buf, binary.LittleEndian, &result)
	if err != nil {
		panic(nil)
	}
	return result
}

func LittleEndianToInt16(b []byte) uint16 {
	if len(b) > 2 {
		panic("Value is too large!")
	}
	if len(b) < 2 {
		b = append(b, make([]byte, 2-len(b))...)
	}
	var result uint16
	buf := bytes.NewReader(b)
	err := binary.Read(buf, binary.LittleEndian, &result)
	if err != nil {
		panic(err)
	}
	return result
}

func LittleEndianToInt32(b []byte) uint32 {
	if len(b) > 4 {
		panic("Value is too large!")
	}
	if len(b) < 4 {
		b = append(b, make([]byte, 4-len(b))...)
	}
	var result uint32
	buf := bytes.NewReader(b)
	err := binary.Read(buf, binary.LittleEndian, &result)
	if err != nil {
		panic(err)
	}
	return result
}

func LittleEndianToInt64(b []byte) uint64 {
	if len(b) > 8 {
		panic("Value is too large!")
	}
	if len(b) < 8 {
		b = append(b, make([]byte, 8-len(b))...)
	}
	var result uint64
	buf := bytes.NewReader(b)
	err := binary.Read(buf, binary.LittleEndian, &result)
	if err != nil {
		panic(err)
	}
	return result
}

func LittleEndianToBigInt(b []byte) *big.Int {
	result := new(big.Int)
	ReverseByteArray(b)
	result.SetBytes(b)
	return result
}

func ByteToLittleEndian(num byte) byte {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.LittleEndian, &num)
	if err != nil {
		panic(err)
	}
	return buf.Bytes()[0]
}

func Int16ToLittleEndian(num uint16) []byte {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.LittleEndian, &num)
	if err != nil {
		panic(err)
	}
	return buf.Bytes()
}

func Int32ToLittleEndian(num uint32) []byte {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.LittleEndian, &num)
	if err != nil {
		panic(err)
	}
	return buf.Bytes()
}

func Int64ToLittleEndian(num uint64) []byte {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.LittleEndian, &num)
	if err != nil {
		panic(err)
	}
	return buf.Bytes()
}

// Takes a byte sequence hash160 and returns a p2pkh address string
func H160ToP2pkhAddress(h160 []byte, testnet bool) string {
	var prefix byte = 0
	if testnet {
		prefix = 0x6f
	}
	b := make([]byte, len(h160)+1)
	b[0] = prefix
	copy(b[1:], h160)
	return EncodeBase58Checksum(b)
}

// Takes a byte sequence hash160 and returns a p2sh address string
func H160ToP2shAddress(h160 []byte, testnet bool) string {
	var prefix byte = 5
	if testnet {
		prefix = 0xc4
	}
	b := make([]byte, len(h160)+1)
	b[0] = prefix
	copy(b[1:], h160)
	return EncodeBase58Checksum(b)
}
