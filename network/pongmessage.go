package network

import (
	"bytes"
)

type PongMessage struct {
	Nonce [8]byte
}

func NewPongMessage(nonce [8]byte) *PongMessage {
	return &PongMessage{Nonce: nonce}
}

func (*PongMessage) Command() []byte {
	return []byte("pong")
}

func (self *PongMessage) Serialize() []byte {
	return self.Nonce[:]
}

func (self *PongMessage) Parse(reader *bytes.Reader) Message {
	reader.Read(self.Nonce[:])
	return self
}

func (self *PongMessage) AckMessage() Message {
	return nil
}
