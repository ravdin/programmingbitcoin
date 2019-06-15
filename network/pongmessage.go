package network

import (
	"bytes"
)

// PongMessage represents a "pong" message in response to a "ping".
type PongMessage struct {
	Nonce [8]byte
}

// NewPongMessage returns a new "pong" message with an 8 byte nonce.
func NewPongMessage(nonce []byte) *PongMessage {
	result := new(PongMessage)
	copy(result.Nonce[:], nonce)
	return result
}

// Command sequence that identifies this type of message.
func (*PongMessage) Command() []byte {
	return []byte("pong")
}

// Serialize this message to send over the network
func (msg *PongMessage) Serialize() []byte {
	return msg.Nonce[:]
}

// Parse a message from a byte steam.
func (msg *PongMessage) Parse(reader *bytes.Reader) Message {
	reader.Read(msg.Nonce[:])
	return msg
}
