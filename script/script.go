package script

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"os"
	"strings"

	"github.com/ravdin/programmingbitcoin/util"
)

// Script represents a Bitcoin script.
type Script struct {
	cmds [][]byte
}

// Add x to y and return the result.
func (scr *Script) Add(x, y *Script) *Script {
	cmds := make([][]byte, len(x.cmds)+len(y.cmds))
	copy(cmds, x.cmds)
	copy(cmds[len(x.cmds):], y.cmds)
	scr.cmds = cmds
	return scr
}

func (scr *Script) String() string {
	result := make([]string, len(scr.cmds))
	for i, cmd := range scr.cmds {
		if len(cmd) == 1 {
			opcode := int(cmd[0])
			if name, ok := opCodeNames[opcode]; ok {
				result[i] = name
			} else {
				result[i] = fmt.Sprintf(`OP_[%d]`, opcode)
			}
		} else {
			result[i] = hex.EncodeToString(cmd)
		}
	}
	return strings.Join(result, " ")
}

// NewScript initializes a new Script object.
func NewScript(cmds [][]byte) *Script {
	return &Script{cmds: cmds}
}

// P2pkhScript takes a hash160 and returns the p2pkh ScriptPubKey
func P2pkhScript(h160 []byte) *Script {
	cmds := [][]byte{
		{0x76},
		{0xa9},
		h160,
		{0x88},
		{0xac},
	}
	return NewScript(cmds)
}

// Parse a new Script from a byte reader.
func Parse(s *bytes.Reader) *Script {
	length := util.ReadVarInt(s)
	var cmds [][]byte
	var count int
	for count < length {
		currentByte, _ := s.ReadByte()
		count++
		if currentByte >= 1 && currentByte <= 75 {
			// we have an cmd set n to be the current byte
			n := int(currentByte)
			// add the next n bytes as an cmd
			buffer := make([]byte, n)
			s.Read(buffer)
			cmds = append(cmds, buffer)
			count += n
		} else if currentByte == 76 {
			// op_pushdata1
			n, _ := s.ReadByte()
			dataLength := int(n)
			buffer := make([]byte, dataLength)
			s.Read(buffer)
			cmds = append(cmds, buffer)
			count += dataLength + 1
		} else if currentByte == 77 {
			// op_pushdata2
			data := make([]byte, 2)
			s.Read(data)
			dataLength := int(util.LittleEndianToInt16(data))
			buffer := make([]byte, dataLength)
			s.Read(buffer)
			cmds = append(cmds, buffer)
			count += dataLength + 2
		} else {
			// we have an opcode. set the current byte to op_code
			opCode := currentByte
			// add the op_code to the list of cmds
			cmds = append(cmds, []byte{opCode})
		}
	}
	if count != length {
		panic("parsing script failed")
	}
	return &Script{cmds: cmds}
}

// Serialize the script as a byte array.
func (scr *Script) Serialize() []byte {
	var raw []byte
	for _, cmd := range scr.cmds {
		length := len(cmd)
		if length == 1 {
			// This is an op code
			raw = append(raw, cmd[0])
		} else {
			// Otherwise, this is an element.
			// for large lengths, we have to use a pushdata opcode
			if length < 76 {
				raw = append(raw, byte(length))
			} else if length >= 76 && length < 0x100 {
				// 76 is pushdata1
				raw = append(raw, byte(76))
				raw = append(raw, byte(length))
			} else if length >= 0x100 && length <= 520 {
				// 77 is pushdata2
				raw = append(raw, byte(77))
				raw = append(raw, util.Int16ToLittleEndian(uint16(length))...)
			}
			raw = append(raw, cmd...)
		}
	}
	total := util.EncodeVarInt(len(raw))
	result := make([]byte, len(total)+len(raw))
	copy(result, total)
	copy(result[len(total):], raw)
	return result
}

// Peek at the stack for a given index.
func (scr *Script) Peek(index int) []byte {
	return scr.cmds[index]
}

// Evaluate the script.
// Return true if the script execution succeeded and false otherwise.
func (scr *Script) Evaluate(z []byte) bool {
	// create a copy as we may need to add to this list if we have a RedeemScript
	cmds := make([][]byte, len(scr.cmds))
	copy(cmds, scr.cmds)
	stack := newOpStack(nil)
	for len(cmds) > 0 {
		cmd := cmds[0]
		cmds = cmds[1:]
		if len(cmd) == 1 {
			// This is an opcode, do what it says.
			opcode := int(cmd[0])
			operation := opCodeFunctions[opcode]
			//fmt.Fprintf(os.Stderr, "Running %s...\n", OpCodeNames[opcode])
			switch opcode {
			case 99, 100:
				// if, notif
				panic("Not implemented")
			case 107, 108:
				// stack to altstack
				panic("Not implemented")
			case 172, 173, 174, 175:
				// Signing operations.
				if !operation(stack, [][]byte{z}) {
					// TODO: Log output
					fmt.Fprintf(os.Stderr, "Op %s failed!\n", opCodeNames[opcode])
					return false
				}
			default:
				if !operation(stack) {
					// TODO: log output
					fmt.Fprintf(os.Stderr, "Op %s failed!\n", opCodeNames[opcode])
					return false
				}
			}
		} else {
			stack.push(cmd)
		}
	}
	if stack.Length == 0 {
		return false
	}
	if bytes.Equal(stack.pop(), []byte{0}) {
		return false
	}
	return true
}
