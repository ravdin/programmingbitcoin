package network

type Message interface {
	Command() []byte
	Serialize() []byte
}
