package script

import (
	"bytes"
	"github.com/ravdin/programmingbitcoin/ecc"
	"github.com/ravdin/programmingbitcoin/util"
	"math/big"
)

func op_verify(stack *OpStack, args ...[][]byte) bool {
	if stack.Length < 1 {
		return false
	}
	elem := stack.Pop()
	return decodeNum(elem) != 0
}

func op_dup(stack *OpStack, args ...[][]byte) bool {
	if stack.Length < 1 {
		return false
	}
	stack.Push(stack.Peek())
	return true
}

func op_equal(stack *OpStack, args ...[][]byte) bool {
	if stack.Length < 2 {
		return false
	}
	item1 := stack.Pop()
	item2 := stack.Pop()
	if bytes.Equal(item1, item2) {
		stack.Push(encodeNum(1))
	} else {
		stack.Push(encodeNum(0))
	}
	return true
}

func op_equalverify(stack *OpStack, args ...[][]byte) bool {
	return op_equal(stack) && op_verify(stack)
}

func op_hash160(stack *OpStack, args ...[][]byte) bool {
	if stack.Length < 1 {
		return false
	}
	element := stack.Pop()
	h160 := util.Hash160(element)
	stack.Push(h160)
	return true
}

func op_sha256(stack *OpStack, args ...[][]byte) bool {
	panic("Not implemented")
}

func op_checksig(stack *OpStack, args ...[][]byte) bool {
	if stack.Length < 2 {
		return false
	}
	z := new(big.Int)
	z.SetBytes(args[0][0])
	// the top element of the stack is the SEC pubkey
	sec_pubkey := stack.Pop()
	// the next element of the stack is the DER signature
	// take off the last byte of the signature as that's the hash_type
	der_signature := stack.Pop()
	der_signature = der_signature[:len(der_signature)-1]
	// parse the serialized pubkey and signature into objects
	point := ecc.ParseS256Point(sec_pubkey)
	sig := ecc.ParseSignature(der_signature)
	if point.Verify(z, sig) {
		stack.Push(encodeNum(1))
	} else {
		stack.Push(encodeNum(0))
	}
	return true
}

func encodeNum(num int) []byte {
	var result []byte = make([]byte, 0)
	if num == 0 {
		return result
	}
	absNum := num
	negative := num < 0
	if negative {
		absNum = -absNum
	}
	for absNum > 0 {
		result = append(result, byte(absNum&0xff))
		absNum >>= 8
	}
	// if the top bit is set,
	// for negative numbers we ensure that the top bit is set
	// for positive numbers we ensure that the top bit is not set
	length := len(result)
	if result[length-1]&0x80 == 0x80 {
		if negative {
			result = append(result, 0x80)
		} else {
			result = append(result, 0)
		}
	} else if negative {
		result[length-1] |= 0x80
	}
	return result
}

func decodeNum(element []byte) int {
	length := len(element)
	if length == 0 {
		return 0
	}
	var result int = int(element[length-1])
	var negative bool = false
	if element[length-1]&0x80 == 0x80 {
		negative = true
		result &= 0x7f
	}
	for i := length - 2; i >= 0; i-- {
		result <<= 8
		result += int(element[i])
	}
	if negative {
		result = -result
	}
	return result
}
