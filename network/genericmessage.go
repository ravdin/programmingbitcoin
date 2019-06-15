package network

import (
	"bytes"
)

// GenericMessage represents a "generic" message with a payload.
type GenericMessage struct {
	CommandVal []byte
	Payload    []byte
}

// NewGenericMessage initializes a GenericMessage object.
func NewGenericMessage(command []byte, payload []byte) *GenericMessage {
	return &GenericMessage{
		CommandVal: command,
		Payload:    payload,
	}
}

// Command sequence that identifies this type of message.
func (msg *GenericMessage) Command() []byte {
	return msg.CommandVal
}

// Serialize this message to send over the network.
func (msg *GenericMessage) Serialize() []byte {
	return msg.Payload
}

// Parse a message from a byte steam.
func (msg *GenericMessage) Parse(reader *bytes.Reader) Message {
	panic("Not implemented!")
}
