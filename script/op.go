package script

import (
//"math/big"
)

type OpCodeFunction func(stack [][]byte, args ...interface{}) bool

func OpHash160(stack [][]byte, args ...interface{}) bool {
	panic("Not implemented")
}

func OpSHA256(stack [][]byte, args ...interface{}) bool {
	panic("Not implemented")
}

func OpCheckSig(stack [][]byte, args ...interface{}) bool {
	panic("Not implemented")
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
