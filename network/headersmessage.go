package network

import (
	"bytes"
	"reflect"

	"github.com/ravdin/programmingbitcoin/block"
	"github.com/ravdin/programmingbitcoin/util"
)

type HeadersMessage struct {
	Blocks []*block.Block
}

func HeadersMessageOption() ReceiveMessageTypeOption {
	return func() reflect.Type {
		return reflect.TypeOf((*HeadersMessage)(nil))
	}
}

func NewHeadersMessage(blocks []*block.Block) *HeadersMessage {
	return &HeadersMessage{Blocks: blocks}
}

func (*HeadersMessage) Command() []byte {
	return []byte("headers")
}

func (self *HeadersMessage) Serialize() []byte {
	panic("Not implemented")
}

func (self *HeadersMessage) Parse(reader *bytes.Reader) Message {
	numHeaders := util.ReadVarInt(reader)
	blocks := make([]*block.Block, numHeaders)
	for i := 0; i < numHeaders; i++ {
		blocks[i] = block.Parse(reader)
		numTx := util.ReadVarInt(reader)
		if numTx != 0 {
			panic("number of txs not 0")
		}
	}
	self.Blocks = blocks
	return self
}
