package network

import "bytes"

// Message defines a common interface for all messages.
type Message interface {
	Command() []byte
	Serialize() []byte
	Parse(reader *bytes.Reader) Message
}
