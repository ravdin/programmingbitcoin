package script

import (
  "math/big"
  "github.com/ravdin/programmingbitcoin/ecc"
  "github.com/ravdin/programmingbitcoin/util"
)

type OpCodeFunction func(stack *[][]byte, args ...interface{}) bool

func OpHash160(stack *[][]byte, args ...interface{}) bool {
	if len(*stack) < 1 {
    return false
  }
  element := pop(stack)
  h160 := util.Hash160(element)
  push(stack, h160)
  return true
}

func OpSHA256(stack *[][]byte, args ...interface{}) bool {
	panic("Not implemented")
}

func OpCheckSig(stack *[][]byte, args ...interface{}) bool {
  if len(*stack) < 2 {
    return false
  }
	z, ok := args[0].(*big.Int)
  if !ok {
    panic ("Failed to cast arg as big.Int!")
  }
  // the top element of the stack is the SEC pubkey
  sec_pubkey := pop(stack)
  // the next element of the stack is the DER signature
  // take off the last byte of the signature as that's the hash_type
  der_signature := pop(stack)
  der_signature = der_signature[:len(der_signature)-1]
  // parse the serialized pubkey and signature into objects
  point := ecc.ParseS256Point(sec_pubkey)
  sig := ecc.ParseSignature(der_signature)
  if point.Verify(z, sig) {
    push(stack, encodeNum(1))
  } else {
    push(stack, encodeNum(0))
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

var OpCodeFunctions = map[int]OpCodeFunction{
	168: OpSHA256,
	172: OpCheckSig,
}

func pop(stack *[][]byte) []byte {
  tmp := *stack
  result := tmp[0]
  *stack = tmp[1:]
  return result
}

func push(stack *[][]byte, item []byte) {
  result := make([][]byte, len(*stack) + 1)
  copy(result[1:], *stack)
  result[0] = item
  *stack = result
}
