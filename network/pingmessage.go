package network

import (
	"bytes"
	"reflect"
)

type PingMessage struct {
	Nonce [8]byte
}

func PingMessageOption() ReceiveMessageTypeOption {
	return func() reflect.Type {
		return reflect.TypeOf((*PingMessage)(nil))
	}
}

func NewPingMessage(nonce [8]byte) *PingMessage {
	return &PingMessage{Nonce: nonce}
}

func (*PingMessage) Command() []byte {
	return []byte("ping")
}

func (self *PingMessage) Serialize() []byte {
	return self.Nonce[:]
}

func (self *PingMessage) Parse(reader *bytes.Reader) Message {
	reader.Read(self.Nonce[:])
	return self
}
