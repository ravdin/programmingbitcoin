package network

import (
	"bytes"
	"reflect"
)

// PingMessage represents a "ping" message.
type PingMessage struct {
	Nonce [8]byte
}

// PingMessageOption returns a helper function to identify this type in a network response.
func PingMessageOption() ReceiveMessageTypeOption {
	return func() reflect.Type {
		return reflect.TypeOf((*PingMessage)(nil))
	}
}

// NewPingMessage returns a new "ping" message with an 8 byte nonce.
func NewPingMessage(nonce [8]byte) *PingMessage {
	return &PingMessage{Nonce: nonce}
}

// Command sequence that identifies this type of message.
func (*PingMessage) Command() []byte {
	return []byte("ping")
}

// Serialize this message to send over the network
func (msg *PingMessage) Serialize() []byte {
	return msg.Nonce[:]
}

// Parse a message from a byte steam.
func (msg *PingMessage) Parse(reader *bytes.Reader) Message {
	reader.Read(msg.Nonce[:])
	return msg
}
