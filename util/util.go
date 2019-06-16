package util

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"math/big"
)

// Useful constants
const (
	SigHashAll     uint32 = 1
	SigHashNone    uint32 = 2
	SigHashSingle  uint32 = 3
	base58Alphabet string = "123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz"
	twoWeeks       int    = 60 * 60 * 24 * 14
	maxTarget      string = `ffff0000000000000000000000000000000000000000000000000000`
)

// HexStringToBytes converts a hex string to a byte array.
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

// HexStringToBigInt converts a hex string to a big.Int.
func HexStringToBigInt(str string) *big.Int {
	result := new(big.Int)
	result.SetBytes(HexStringToBytes(str))
	return result
}

// ReverseByteArray reverses a byte array in place and returns the result.
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
			count++
		} else {
			break
		}
	}
	var encoded []int
	prefix := make([]int, count)
	num := new(big.Int)
	mod := new(big.Int)
	b58 := big.NewInt(58)
	num.SetBytes(chars)
	for num.Sign() > 0 {
		num.QuoRem(num, b58, mod)
		encoded = append(encoded, int(mod.Int64()))
	}
	encoded = append(encoded, prefix...)
	alphabet := []byte(base58Alphabet)
	result := make([]byte, len(encoded))
	for i := len(encoded) - 1; i >= 0; i-- {
		result[len(encoded)-i-1] = alphabet[encoded[i]]
	}
	return string(result)
}

// EncodeBase58Checksum returns a base58 encoded string with an appended checksum.
func EncodeBase58Checksum(b []byte) string {
	return encodeBase58(string(append(b, Hash256(b)[:4]...)))
}

// DecodeBase58 decodes a base58 string and verifies the checksum.
func DecodeBase58(encoded string) []byte {
	num := big.NewInt(0)
	b58 := big.NewInt(58)
	alphabet := []byte(base58Alphabet)
	chars := []byte(encoded)
	for _, c := range chars {
		num.Mul(num, b58)
		num.Add(num, big.NewInt(int64(bytes.IndexByte(alphabet, c))))
	}
	combined := num.Bytes()
	length := len(combined)
	checksum := combined[length-4:]
	if !bytes.Equal(Hash256(combined[:length-4])[:4], checksum) {
		panic(fmt.Sprintf("Bad address: %x %x", checksum, Hash256(combined[:length-4])[:4]))
	}
	return combined[1 : length-4]
}

// IntToBytes returns a byte array of a given size from a big.Int.
func IntToBytes(num *big.Int, size int) []byte {
	result := make([]byte, size)
	var raw = num.Bytes()
	copy(result[size-len(raw):], raw)
	return result
}

// ReadVarInt reads an integer of a variable size (up to 8 bytes) from a byte reader.
func ReadVarInt(r *bytes.Reader) int {
	b, err := r.ReadByte()
	if err != nil {
		panic(err)
	}
	var bufsize int
	switch b {
	case 0xfd:
		bufsize = 2
	case 0xfe:
		bufsize = 4
	case 0xff:
		bufsize = 8
	default:
		return int(b)
	}
	buffer := make([]byte, bufsize)
	r.Read(buffer)
	return int(LittleEndianToInt64(buffer))
}

// EncodeVarInt returns a variable size byte array from an integer.
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

// LittleEndianToInt16 returns a uint16 from a byte array.
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

// LittleEndianToInt32 returns a uint32 from a byte array.
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

// LittleEndianToInt64 returns a uint64 from a byte array.
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

// LittleEndianToBigInt returns a big.Int from a byte array.
func LittleEndianToBigInt(b []byte) *big.Int {
	result := new(big.Int)
	ReverseByteArray(b)
	result.SetBytes(b)
	return result
}

// Int16ToLittleEndian returns a byte array from a uint16.
func Int16ToLittleEndian(num uint16) []byte {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.LittleEndian, &num)
	if err != nil {
		panic(err)
	}
	return buf.Bytes()
}

// Int32ToLittleEndian returns a byte array from a uint32.
func Int32ToLittleEndian(num uint32) []byte {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.LittleEndian, &num)
	if err != nil {
		panic(err)
	}
	return buf.Bytes()
}

// Int64ToLittleEndian returns a byte array from a uint64.
func Int64ToLittleEndian(num uint64) []byte {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.LittleEndian, &num)
	if err != nil {
		panic(err)
	}
	return buf.Bytes()
}

// H160ToP2pkhAddress takes a byte sequence hash160 and returns a p2pkh address string
func H160ToP2pkhAddress(h160 []byte, testnet bool) string {
	var prefix byte
	if testnet {
		prefix = 0x6f
	}
	b := make([]byte, len(h160)+1)
	b[0] = prefix
	copy(b[1:], h160)
	return EncodeBase58Checksum(b)
}

// H160ToP2shAddress takes a byte sequence hash160 and returns a p2sh address string
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

// BitsToTarget turns bits into a target (large 256-bit integer)
func BitsToTarget(bits []byte) *big.Int {
	length := len(bits)
	exponent := bits[length-1]
	coefficient := new(big.Int)
	cbits := make([]byte, length-1)
	copy(cbits[:], bits[:length-1])
	ReverseByteArray(cbits)
	coefficient.SetBytes(cbits)
	// The formula is coefficient * 256**(exponent-3)
	result := big.NewInt(256)
	result.Exp(result, big.NewInt(int64(exponent-3)), nil)
	result.Mul(result, coefficient)
	return result
}

// TargetToBits turns a target integer back into bits
func TargetToBits(target *big.Int) []byte {
	rawBytes := target.Bytes()
	coefficient := make([]byte, 3)
	var exponent int
	if rawBytes[0] > 0x7f {
		copy(coefficient[1:], rawBytes[:2])
		exponent = len(rawBytes) + 1
	} else {
		copy(coefficient, rawBytes[:3])
		exponent = len(rawBytes)
	}
	result := make([]byte, 4)
	copy(result, ReverseByteArray(coefficient))
	result[3] = byte(exponent)
	return result
}

// CalculateNewBits calculates the new bits given a 2016-block time differential and the previous bits
func CalculateNewBits(previousBits []byte, timeDifferential int) []byte {
	// if the time differential is greater than 8 weeks, set to 8 weeks
	if timeDifferential > twoWeeks*4 {
		timeDifferential = twoWeeks * 4
	}
	// if the time differential is less than half a week, set to half a week
	if timeDifferential < twoWeeks/4 {
		timeDifferential = twoWeeks / 4
	}
	// the new target is the previous target * time differential / two weeks
	target := BitsToTarget(previousBits)
	target.Mul(target, big.NewInt(int64(timeDifferential)))
	target.Div(target, big.NewInt(int64(twoWeeks)))

	biMaxTarget := HexStringToBigInt(maxTarget)
	if target.Cmp(biMaxTarget) > 0 {
		target = biMaxTarget
	}

	// convert the new target to bits
	return TargetToBits(target)
}

// MerkleParent takes the binary hashes and calculates the hash256
func MerkleParent(hash1, hash2 []byte) []byte {
	return Hash256(append(hash1, hash2...))
}

// MerkleParentLevel takes a list of binary hashes and returns a list that's half the length
func MerkleParentLevel(hashes [][]byte) [][]byte {
	if len(hashes) == 1 {
		panic("Cannot take a parent level with only 1 item")
	}
	length := len(hashes)
	result := make([][]byte, (length+1)/2)
	for i := range result {
		hash1 := hashes[2*i]
		hash2 := hashes[2*i]
		if i*2 < length-1 {
			hash2 = hashes[2*i+1]
		}
		result[i] = MerkleParent(hash1, hash2)
	}
	return result
}

// MerkleRoot takes a list of binary hashes and returns the merkle root
func MerkleRoot(hashes [][]byte) []byte {
	current := hashes
	for len(current) > 1 {
		current = MerkleParentLevel(current)
	}
	return current[0]
}

// BitFieldToBytes converts a bit field to a byte array.
func BitFieldToBytes(bits []byte) []byte {
	if len(bits)%8 != 0 {
		panic("bits does not have a length that is divisible by 8")
	}
	result := make([]byte, len(bits)/8)
	for i, bit := range bits {
		byteIndex := i / 8
		bitIndex := uint(i % 8)
		if bit != 0 {
			result[byteIndex] |= 1 << bitIndex
		}
	}
	return result
}

// BytesToBitField converts a byte array to a bit field.
func BytesToBitField(bytes []byte) []byte {
	result := make([]byte, len(bytes)*8)
	for byteIndex, b := range bytes {
		var mask byte = 1
		bits := make([]byte, 8)
		for i := 0; i < 8; i++ {
			if b&mask != 0 {
				bits[i] = 1
			}
			mask <<= 1
		}
		copy(result[byteIndex*8:(byteIndex+1)*8], bits)
	}
	return result
}
