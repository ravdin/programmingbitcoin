package network

import (
	"bytes"
	"reflect"

	"github.com/ravdin/programmingbitcoin/block"
	"github.com/ravdin/programmingbitcoin/util"
)

// HeadersMessage represents a "headers" message with a sequence of blocks.
type HeadersMessage struct {
	Blocks []*block.Block
}

// HeadersMessageOption returns a helper function to identify this type in a network response.
func HeadersMessageOption() ReceiveMessageTypeOption {
	return func() reflect.Type {
		return reflect.TypeOf((*HeadersMessage)(nil))
	}
}

// NewHeadersMessage returns a new HeadersMessage
// blocks: an array of Block objects.
func NewHeadersMessage(blocks []*block.Block) *HeadersMessage {
	return &HeadersMessage{Blocks: blocks}
}

// Command sequence that identifies this type of message.
func (*HeadersMessage) Command() []byte {
	return []byte("headers")
}

// Serialize this message to send over the network
func (msg *HeadersMessage) Serialize() []byte {
	panic("Not implemented")
}

// Parse a message from a byte steam.
func (msg *HeadersMessage) Parse(reader *bytes.Reader) Message {
	numHeaders := util.ReadVarInt(reader)
	blocks := make([]*block.Block, numHeaders)
	for i := 0; i < numHeaders; i++ {
		blocks[i] = block.Parse(reader)
		numTx := util.ReadVarInt(reader)
		if numTx != 0 {
			panic("number of txs not 0")
		}
	}
	msg.Blocks = blocks
	return msg
}
