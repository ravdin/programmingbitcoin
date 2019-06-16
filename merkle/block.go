package merkle

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/ravdin/programmingbitcoin/util"
)

// Block represents a merkle block.
type Block struct {
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

// NewBlock returns a new merkle block.
func NewBlock(version uint32,
	prevBlock []byte,
	merkleRoot []byte,
	timestamp uint32,
	bits []byte,
	nonce []byte,
	total uint32,
	hashes [][]byte,
	flags []byte) *Block {
	result := &Block{
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

func (block *Block) String() string {
	result := make([]string, 3)
	hashes := make([]string, len(block.Hashes))
	for i, h := range block.Hashes {
		hashes[i] = fmt.Sprintf("\t%x", h)
	}
	result[0] = fmt.Sprintf("%d", block.Version)
	result[1] = strings.Join(hashes, "")
	result[2] = fmt.Sprintf("%x", block.Flags)
	return strings.Join(result, "\n")
}

// Parse a merkle block from a byte reader. Returns a Merkle Block object
func (block *Block) Parse(reader *bytes.Reader) *Block {
	buffer4 := make([]byte, 4)
	buffer32 := make([]byte, 32)
	reader.Read(buffer4)
	block.Version = util.LittleEndianToInt32(buffer4)
	reader.Read(buffer32)
	copy(block.PrevBlock[:32], util.ReverseByteArray(buffer32))
	reader.Read(buffer32)
	copy(block.MerkleRoot[:32], util.ReverseByteArray(buffer32))
	reader.Read(buffer4)
	block.Timestamp = util.LittleEndianToInt32(buffer4)
	reader.Read(buffer4)
	copy(block.Bits[:4], buffer4)
	reader.Read(buffer4)
	copy(block.Nonce[:4], buffer4)
	reader.Read(buffer4)
	block.Total = util.LittleEndianToInt32(buffer4)
	numHashes := util.ReadVarInt(reader)
	hashes := make([][]byte, numHashes)
	for i := range hashes {
		reader.Read(buffer32)
		hashes[i] = make([]byte, 32)
		copy(hashes[i], util.ReverseByteArray(buffer32))
	}
	block.Hashes = hashes
	flagLength := util.ReadVarInt(reader)
	block.Flags = make([]byte, flagLength)
	reader.Read(block.Flags)
	return block
}

// IsValid verifies whether the merkle tree information validates to the merkle root
func (block *Block) IsValid() bool {
	flagBits := util.BytesToBitField(block.Flags)
	hashes := make([][]byte, len(block.Hashes))
	for i, hash := range block.Hashes {
		hashes[i] = make([]byte, len(hash))
		copy(hashes[i], hash)
		util.ReverseByteArray(hashes[i])
	}
	tree := NewTree(int(block.Total))
	tree.PopulateTree(flagBits, hashes)
	return bytes.Equal(util.ReverseByteArray(tree.Root()), block.MerkleRoot[:])
}
