package script

import (
	"bytes"
	"fmt"
	"math/big"
	"os"

	"github.com/ravdin/programmingbitcoin/ecc"
	"github.com/ravdin/programmingbitcoin/util"
)

func op0(stack *opStack, args ...[][]byte) bool {
	stack.push(encodeNum(0))
	return true
}

func opVerify(stack *opStack, args ...[][]byte) bool {
	if stack.Length < 1 {
		return false
	}
	elem := stack.pop()
	return decodeNum(elem) != 0
}

func opDup(stack *opStack, args ...[][]byte) bool {
	if stack.Length < 1 {
		return false
	}
	stack.push(stack.peek())
	return true
}

func opEqual(stack *opStack, args ...[][]byte) bool {
	if stack.Length < 2 {
		return false
	}
	item1 := stack.pop()
	item2 := stack.pop()
	if bytes.Equal(item1, item2) {
		stack.push(encodeNum(1))
	} else {
		stack.push(encodeNum(0))
	}
	return true
}

func opEqualverify(stack *opStack, args ...[][]byte) bool {
	return opEqual(stack) && opVerify(stack)
}

func opHash160(stack *opStack, args ...[][]byte) bool {
	if stack.Length < 1 {
		return false
	}
	element := stack.pop()
	h160 := util.Hash160(element)
	stack.push(h160)
	return true
}

func opSha256(stack *opStack, args ...[][]byte) bool {
	panic("Not implemented")
}

func opChecksig(stack *opStack, args ...[][]byte) bool {
	if stack.Length < 2 {
		return false
	}
	z := new(big.Int)
	z.SetBytes(args[0][0])
	// the top element of the stack is the SEC pubkey
	secPubkey := stack.pop()
	// the next element of the stack is the DER signature
	// take off the last byte of the signature as that's the hash_type
	derSignature := stack.pop()
	derSignature = derSignature[:len(derSignature)-1]
	// parse the serialized pubkey and signature into objects
	point := ecc.ParseS256Point(secPubkey)
	sig := ecc.ParseSignature(derSignature)
	if point.Verify(z, sig) {
		stack.push(encodeNum(1))
	} else {
		stack.push(encodeNum(0))
	}
	return true
}

func opCheckmultisig(stack *opStack, args ...[][]byte) bool {
	z := new(big.Int)
	z.SetBytes(args[0][0])
	if stack.Length < 1 {
		return false
	}
	n := decodeNum(stack.pop())
	if stack.Length < n+1 {
		return false
	}
	secPubkeys := make([][]byte, n)
	for i := 0; i < n; i++ {
		secPubkeys[i] = stack.pop()
	}
	m := decodeNum(stack.pop())
	if stack.Length < m+1 {
		return false
	}
	derSignatures := make([][]byte, m)
	for i := 0; i < m; i++ {
		sig := stack.pop()
		sigLength := len(sig)
		derSignatures[i] = sig[:sigLength-1]
	}
	// OP_CHECKMULTISIG bug
	stack.pop()
	secIndex := 0
	for derIndex := 0; derIndex < m; derIndex++ {
		if secIndex >= n {
			fmt.Fprintf(os.Stderr, "signatures no good or not in right order\n")
			return false
		}
		sig := ecc.ParseSignature(derSignatures[derIndex])
		for secIndex < n {
			point := ecc.ParseS256Point(secPubkeys[secIndex])
			secIndex++
			if point.Verify(z, sig) {
				break
			}
		}
	}
	// The signatures are valid, push a 1 to the stack.
	stack.push(encodeNum(1))
	return true
}

func encodeNum(num int) []byte {
	result := make([]byte, 0)
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
	result := int(element[length-1])
	negative := false
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
