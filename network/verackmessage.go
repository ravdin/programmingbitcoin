package network

import (
	"bytes"
	"reflect"
)

type VerackMessage struct {
}

func VerackMessageOption() ReceiveMessageTypeOption {
	return func() reflect.Type {
		return reflect.TypeOf((*VerackMessage)(nil))
	}
}

func NewVerackMessage() *VerackMessage {
	return new(VerackMessage)
}

func (*VerackMessage) Command() []byte {
	return []byte("verack")
}

func (self *VerackMessage) Serialize() []byte {
	return make([]byte, 0)
}

func (self *VerackMessage) Parse(reader *bytes.Reader) Message {
	return self
}
