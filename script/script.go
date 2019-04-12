package script

import (
	"bytes"
	"github.com/ravdin/programmingbitcoin/util"
)

type Script struct {
	cmds []byte
}

func Parse(s bytes.Reader) *Script {
	// get the length of the entire field
	length := int(util.ReadVarInt(s))
	// initialize the cmds array
	var cmds []byte
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
			cmds = append(cmds, buffer...)
			// increase the count by n
			count += n
		} else if current_byte == 76 {
			// op_pushdata1
			n, _ := s.ReadByte()
			data_length := int(util.LittleEndianToInt16([]byte{n}))
			buffer := make([]byte, data_length)
			s.Read(buffer)
			cmds = append(cmds, buffer...)
			count += data_length + 1
		} else if current_byte == 77 {
			// op_pushdata2
			var data []byte = make([]byte, 2)
			s.Read(data)
			data_length := int(util.LittleEndianToInt16(data))
			buffer := make([]byte, data_length)
			s.Read(buffer)
			cmds = append(cmds, buffer...)
			count += data_length + 2
		} else {
			// we have an opcode. set the current byte to op_code
			op_code := current_byte
			// add the op_code to the list of cmds
			cmds = append(cmds, op_code)
		}
	}
	if count != length {
		panic("parsing script failed")
	}
	return &Script{cmds: cmds}
}
