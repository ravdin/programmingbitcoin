package network

import "bytes"

type Message interface {
	// Command sequence that identifies this type of message.
	Command() []byte
	// Serialize this message to send over the network
	Serialize() []byte
	// Read a message from a byte steam.
	Parse(reader *bytes.Reader) Message
	// Message that should be returned after reading an incoming message.
	AckMessage() Message
}
