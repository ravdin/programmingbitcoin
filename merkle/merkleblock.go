package merkle

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/ravdin/programmingbitcoin/util"
)

type MerkleBlock struct {
	Version    uint32
	PrevBlock  [32]byte
	MerkleRoot [32]byte
	Timestamp  uint32
	Bits       [4]byte
	Nonce      [4]byte
	Total      uint32
	Hashes     [][]byte
	Flags      []byte
}

func NewMerkleBlock(version uint32,
	prevBlock []byte,
	merkleRoot []byte,
	timestamp uint32,
	bits []byte,
	nonce []byte,
	total uint32,
	hashes [][]byte,
	flags []byte) *MerkleBlock {
	result := &MerkleBlock{
		Version:   version,
		Timestamp: timestamp,
		Total:     total,
		Hashes:    hashes,
		Flags:     flags,
	}
	copy(result.PrevBlock[:32], prevBlock)
	copy(result.MerkleRoot[:32], merkleRoot)
	copy(result.Bits[:4], bits)
	copy(result.Nonce[:4], nonce)
	return result
}

func (self *MerkleBlock) String() string {
	result := make([]string, 3)
	hashes := make([]string, len(self.Hashes))
	for i, h := range self.Hashes {
		hashes[i] = fmt.Sprintf("\t%x", h)
	}
	result[0] = fmt.Sprintf("%d", self.Version)
	result[1] = strings.Join(hashes, "")
	result[2] = fmt.Sprintf("%x", self.Flags)
	return strings.Join(result, "\n")
}

// Takes a byte stream and parses a merkle block. Returns a Merkle Block object
func (self *MerkleBlock) Parse(reader *bytes.Reader) *MerkleBlock {
	buffer4 := make([]byte, 4)
	buffer32 := make([]byte, 32)
	reader.Read(buffer4)
	self.Version = util.LittleEndianToInt32(buffer4)
	reader.Read(buffer32)
	copy(self.PrevBlock[:32], util.ReverseByteArray(buffer32))
	reader.Read(buffer32)
	copy(self.MerkleRoot[:32], util.ReverseByteArray(buffer32))
	reader.Read(buffer4)
	self.Timestamp = util.LittleEndianToInt32(buffer4)
	reader.Read(buffer4)
	copy(self.Bits[:4], buffer4)
	reader.Read(buffer4)
	copy(self.Nonce[:4], buffer4)
	reader.Read(buffer4)
	self.Total = util.LittleEndianToInt32(buffer4)
	numHashes := util.ReadVarInt(reader)
	hashes := make([][]byte, numHashes)
	for i := range hashes {
		reader.Read(buffer32)
		hashes[i] = make([]byte, 32)
		copy(hashes[i], util.ReverseByteArray(buffer32))
	}
	self.Hashes = hashes
	flagLength := util.ReadVarInt(reader)
	self.Flags = make([]byte, flagLength)
	reader.Read(self.Flags)
	return self
}

// Verifies whether the merkle tree information validates to the merkle root
func (self *MerkleBlock) IsValid() bool {
	flagBits := util.BytesToBitField(self.Flags)
	hashes := make([][]byte, len(self.Hashes))
	for i, hash := range self.Hashes {
		hashes[i] = make([]byte, len(hash))
		copy(hashes[i], hash)
		util.ReverseByteArray(hashes[i])
	}
	tree := NewMerkleTree(int(self.Total))
	tree.PopulateTree(flagBits, hashes)
	return bytes.Equal(util.ReverseByteArray(tree.Root()), self.MerkleRoot[:])
}
