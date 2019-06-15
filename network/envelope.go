package network

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"os"

	"github.com/ravdin/programmingbitcoin/util"
)

// Envelope wraps a network message.
type Envelope struct {
	Command []byte
	Payload []byte
	Magic   [4]byte
}

var (
	networkMagic     = [4]byte{0xf9, 0xbe, 0xb4, 0xd9}
	testNetworkMagic = [4]byte{0x0b, 0x11, 0x09, 0x07}
)

// NewEnvelope initializes a new Envelope.
func NewEnvelope(command []byte, payload []byte, testnet bool) *Envelope {
	var magic [4]byte
	if testnet {
		magic = testNetworkMagic
	} else {
		magic = networkMagic
	}
	return &Envelope{Command: command, Payload: payload, Magic: magic}
}

// ParseEnvelope takes a stream and creates a network.Envelope
func ParseEnvelope(reader *bytes.Reader, testnet bool) *Envelope {
	if reader.Len() == 0 {
		panic("Connection reset!")
	}
	magic := make([]byte, 4)
	reader.Read(magic)
	var expectedMagic [4]byte
	if testnet {
		expectedMagic = testNetworkMagic
	} else {
		expectedMagic = networkMagic
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
	checksum := make([]byte, 4)
	reader.Read(checksum)
	payload := make([]byte, payloadLength)
	reader.Read(payload)
	// Verify checksum
	if !bytes.Equal(util.Hash256(payload)[:4], checksum) {
		fmt.Fprintf(os.Stderr, "%x %x\n", util.Hash256(payload)[:4], checksum)
		panic("Invalid checksum!")
	}
	return NewEnvelope(command, payload, testnet)
}

// Serialize returns the byte serialization of the entire network message
func (env *Envelope) Serialize() []byte {
	command := make([]byte, 12)
	copy(command, env.Command)
	payloadLength := len(env.Payload)
	checksum := util.Hash256(env.Payload)[:4]
	result := make([]byte, payloadLength+24)
	copy(result[:4], env.Magic[:])
	copy(result[4:16], command)
	copy(result[16:20], util.Int32ToLittleEndian(uint32(payloadLength)))
	copy(result[20:24], checksum)
	copy(result[24:], env.Payload)
	return result
}

func (env *Envelope) String() string {
	return fmt.Sprintf("%s: %s", string(env.Command), hex.EncodeToString(env.Payload))
}

// Stream opens a new Reader for the payload's bytes.
func (env *Envelope) Stream() *bytes.Reader {
	return bytes.NewReader(env.Payload)
}
