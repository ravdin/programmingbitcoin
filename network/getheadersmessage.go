package network

import (
	"bytes"

	"github.com/ravdin/programmingbitcoin/util"
)

type GetHeadersMessage struct {
	Version    uint32
	NumHashes  int
	StartBlock [32]byte
	EndBlock   [32]byte
}

const (
	defaultVersion   uint32 = 70015
	defaultNumHashes int    = 1
)

func NewGetHeadersMessage(startBlock []byte) *GetHeadersMessage {
	var startBlockData [32]byte
	var endBlockData [32]byte
	copy(startBlockData[:], startBlock[:32])
	return &GetHeadersMessage{
		Version:    defaultVersion,
		NumHashes:  defaultNumHashes,
		StartBlock: startBlockData,
		EndBlock:   endBlockData,
	}
}

func (*GetHeadersMessage) Command() []byte {
	return []byte("getheaders")
}

func (self *GetHeadersMessage) Serialize() []byte {
	version := util.Int32ToLittleEndian(self.Version)
	numHashes := util.EncodeVarInt(self.NumHashes)
	startBlock := make([]byte, 32)
	endBlock := make([]byte, 32)
	copy(startBlock, self.StartBlock[:])
	copy(endBlock, self.EndBlock[:])
	util.ReverseByteArray(startBlock)
	util.ReverseByteArray(endBlock)
	result := make([]byte, 68+len(numHashes))
	copy(result[:4], version)
	copy(result[4:4+len(numHashes)], numHashes)
	copy(result[4+len(numHashes):36+len(numHashes)], startBlock)
	copy(result[36+len(numHashes):], endBlock)
	return result
}

func (self *GetHeadersMessage) Parse(reader *bytes.Reader) Message {
	version := make([]byte, 4)
	reader.Read(version)
	self.Version = util.LittleEndianToInt32(version)
	self.NumHashes = util.ReadVarInt(reader)
	blockData := make([]byte, 32)
	reader.Read(blockData)
	copy(self.StartBlock[:], util.ReverseByteArray(blockData))
	reader.Read(blockData)
	copy(self.EndBlock[:], util.ReverseByteArray(blockData))
	return self
}
