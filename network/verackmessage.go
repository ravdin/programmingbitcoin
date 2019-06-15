package network

import (
	"bytes"
	"reflect"
)

// VerackMessage represents a "verack" message we expect to receive in response to a handshake.
type VerackMessage struct {
}

// VerackMessageOption returns a helper function to identify this type in a network response.
func VerackMessageOption() ReceiveMessageTypeOption {
	return func() reflect.Type {
		return reflect.TypeOf((*VerackMessage)(nil))
	}
}

// NewVerackMessage returns a new VerackMessage.
func NewVerackMessage() *VerackMessage {
	return new(VerackMessage)
}

// Command sequence that identifies this type of message.
func (*VerackMessage) Command() []byte {
	return []byte("verack")
}

// Serialize this message to send over the network.
func (msg *VerackMessage) Serialize() []byte {
	return make([]byte, 0)
}

// Parse a message from a byte steam.
func (msg *VerackMessage) Parse(reader *bytes.Reader) Message {
	return msg
}
