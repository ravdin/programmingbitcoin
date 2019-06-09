package network

import (
	"bytes"
)

type GenericMessage struct {
	CommandVal []byte
	Payload    []byte
}

func NewGenericMessage(command []byte, payload []byte) *GenericMessage {
	return &GenericMessage{
		CommandVal: command,
		Payload:    payload,
	}
}

func (self *GenericMessage) Command() []byte {
	return self.CommandVal
}

func (self *GenericMessage) Serialize() []byte {
	return self.Payload
}

func (self *GenericMessage) Parse(reader *bytes.Reader) Message {
	panic("Not implemented!")
}
