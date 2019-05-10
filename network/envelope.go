package network

import (
	"bytes"
	"encoding/hex"
	"fmt"

	"github.com/ravdin/programmingbitcoin/util"
)

type Envelope struct {
	Command []byte
	Payload []byte
	Magic   [4]byte
}

var (
	NETWORK_MAGIC      = [4]byte{0xf9, 0xbe, 0xb4, 0xd9}
	TEST_NETWORK_MAGIC = [4]byte{0x0b, 0x11, 0x09, 0x07}
)

func NewEnvelope(command []byte, payload []byte, testnet bool) *Envelope {
	var magic [4]byte
	if testnet {
		magic = TEST_NETWORK_MAGIC
	} else {
		magic = NETWORK_MAGIC
	}
	return &Envelope{Command: command, Payload: payload, Magic: magic}
}

func ParseEnvelope(reader *bytes.Reader, testnet bool) *Envelope {
	magic := make([]byte, 4)
	reader.Read(magic)
	var expectedMagic [4]byte
	if testnet {
		expectedMagic = TEST_NETWORK_MAGIC
	} else {
		expectedMagic = NETWORK_MAGIC
	}
	if !bytes.Equal(magic[:], expectedMagic[:]) {
		panic(fmt.Sprintf("magic is not right %v vs %v", hex.EncodeToString(magic[:]), hex.EncodeToString(expectedMagic[:])))
	}
	command := make([]byte, 12)
	reader.Read(command)
	command = bytes.TrimRightFunc(command, func(r rune) bool { return r == 0 })
	buffer := make([]byte, 4)
	reader.Read(buffer)
	payloadLength := util.LittleEndianToInt32(buffer)
	reader.Read(buffer)
	payload := make([]byte, payloadLength)
	reader.Read(payload)
	// Verify checksum
	if !bytes.Equal(util.Hash256(payload)[:4], buffer) {
		panic("Invalid checksum!")
	}
	return NewEnvelope(command, payload, testnet)
}

func (self *Envelope) Serialize() []byte {
	command := make([]byte, 12)
	copy(command, self.Command)
	payloadLength := len(self.Payload)
	checksum := util.Hash256(self.Payload)[:4]
	result := make([]byte, payloadLength+24)
	copy(result[:4], self.Magic[:])
	copy(result[4:16], command)
	copy(result[16:20], util.Int32ToLittleEndian(uint32(payloadLength)))
	copy(result[20:24], checksum)
	copy(result[24:], self.Payload)
	return result
}

func (self *Envelope) String() string {
	return fmt.Sprintf("%s: %s", string(self.Command), hex.EncodeToString(self.Payload))
}

func (self *Envelope) Stream() *bytes.Reader {
	return bytes.NewReader(self.Payload)
}
