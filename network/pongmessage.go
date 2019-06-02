package network

import (
	"bytes"
)

type PongMessage struct {
	Nonce [8]byte
}

func NewPongMessage(nonce []byte) *PongMessage {
	result := new(PongMessage)
	copy(result.Nonce[:], nonce)
	return result
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
