package script

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"os"
	"strings"

	"github.com/ravdin/programmingbitcoin/util"
)

type Script struct {
	cmds [][]byte
}

func (z *Script) Add(x, y *Script) *Script {
	cmds := make([][]byte, len(x.cmds)+len(y.cmds))
	copy(cmds, x.cmds)
	copy(cmds[len(x.cmds):], y.cmds)
	z.cmds = cmds
	return z
}

func (self *Script) String() string {
	result := make([]string, len(self.cmds))
	for i, cmd := range self.cmds {
		if len(cmd) == 1 {
			opcode := int(cmd[0])
			if name, ok := OpCodeNames[opcode]; ok {
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

func NewScript(cmds [][]byte) *Script {
	return &Script{cmds: cmds}
}

// Takes a hash160 and returns the p2pkh ScriptPubKey
func P2pkhScript(h160 []byte) *Script {
	cmds := [][]byte{
		[]byte{0x76},
		[]byte{0xa9},
		h160,
		[]byte{0x88},
		[]byte{0xac},
	}
	return NewScript(cmds)
}

func Parse(s *bytes.Reader) *Script {
	// get the length of the entire field
	length := util.ReadVarInt(s)
	// initialize the cmds array
	var cmds [][]byte
	// initialize the number of bytes we've read to 0
	count := 0
	// loop until we've read length bytes
	for count < length {
		// get the current byte
		current_byte, _ := s.ReadByte()
		// increment the bytes we've read
		count++
		// if the current byte is between 1 and 75 inclusive
		if current_byte >= 1 && current_byte <= 75 {
			// we have an cmd set n to be the current byte
			n := int(current_byte)
			// add the next n bytes as an cmd
			buffer := make([]byte, n)
			s.Read(buffer)
			cmds = append(cmds, buffer)
			// increase the count by n
			count += n
		} else if current_byte == 76 {
			// op_pushdata1
			n, _ := s.ReadByte()
			data_length := int(util.LittleEndianToInt16([]byte{n}))
			buffer := make([]byte, data_length)
			s.Read(buffer)
			cmds = append(cmds, buffer)
			count += data_length + 1
		} else if current_byte == 77 {
			// op_pushdata2
			var data []byte = make([]byte, 2)
			s.Read(data)
			data_length := int(util.LittleEndianToInt16(data))
			buffer := make([]byte, data_length)
			s.Read(buffer)
			cmds = append(cmds, buffer)
			count += data_length + 2
		} else {
			// we have an opcode. set the current byte to op_code
			op_code := current_byte
			// add the op_code to the list of cmds
			cmds = append(cmds, []byte{op_code})
		}
	}
	if count != length {
		panic("parsing script failed")
	}
	return &Script{cmds: cmds}
}

func (self *Script) Serialize() []byte {
	var raw []byte
	for _, cmd := range self.cmds {
		length := len(cmd)
		if length == 1 {
			// This is an op code
			raw = append(raw, cmd[0])
		} else {
			// Otherwise, this is an element.
			// for large lengths, we have to use a pushdata opcode
			if length < 76 {
				raw = append(raw, util.ByteToLittleEndian(byte(length)))
			} else if length >= 76 && length < 0x100 {
				// 76 is pushdata1
				raw = append(raw, util.ByteToLittleEndian(76))
				raw = append(raw, util.ByteToLittleEndian(byte(length)))
			} else if length >= 0x100 && length <= 520 {
				// 77 is pushdata2
				raw = append(raw, util.ByteToLittleEndian(77))
				raw = append(raw, util.Int16ToLittleEndian(uint16(length))...)
			}
			raw = append(raw, cmd...)
		}
	}
	total := util.EncodeVarInt(len(raw))
	var result []byte = make([]byte, len(total)+len(raw))
	copy(result, total)
	copy(result[len(total):], raw)
	return result
}

func (self *Script) Peek(index int) []byte {
	return self.cmds[index]
}

func (self *Script) Evaluate(z []byte) bool {
	// create a copy as we may need to add to this list if we have a RedeemScript
	cmds := make([][]byte, len(self.cmds))
	copy(cmds, self.cmds)
	stack := NewOpStack(nil)
	for len(cmds) > 0 {
		cmd := cmds[0]
		cmds = cmds[1:]
		if len(cmd) == 1 {
			// This is an opcode, do what it says.
			opcode := int(cmd[0])
			operation := OpCodeFunctions[opcode]
			fmt.Fprintf(os.Stderr, "Running %s...\n", OpCodeNames[opcode])
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
					fmt.Fprintf(os.Stderr, "Op %s failed!\n", OpCodeNames[opcode])
					return false
				}
			default:
				if !operation(stack) {
					// TODO: log output
					fmt.Fprintf(os.Stderr, "Op %s failed!\n", OpCodeNames[opcode])
					return false
				}
			}
		} else {
			stack.Push(cmd)
		}
	}
	if stack.Length == 0 {
		return false
	}
	if bytes.Equal(stack.Pop(), []byte{0}) {
		return false
	}
	return true
}
