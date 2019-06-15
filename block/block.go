package block

import (
	"bytes"
	"math/big"

	"github.com/ravdin/programmingbitcoin/util"
)

// Useful constants
var (
	GenesisBlock     []byte = util.HexStringToBytes(`0100000000000000000000000000000000000000000000000000000000000000000000003ba3edfd7a7b12b27ac72c3e67768f617fc81bc3888a51323a9fb8aa4b1e5e4a29ab5f49ffff001d1dac2b7c`)
	TestGenesisBlock []byte = util.HexStringToBytes(`0100000000000000000000000000000000000000000000000000000000000000000000003ba3edfd7a7b12b27ac72c3e67768f617fc81bc3888a51323a9fb8aa4b1e5e4adae5494dffff001d1aa4ae18`)
	LowestBits       []byte = util.HexStringToBytes(`ffff001d`)
)

// Block is a batch of transations.
type Block struct {
	Version    uint32
	PrevBlock  [32]byte
	MerkleRoot [32]byte
	Timestamp  uint32
	Bits       [4]byte
	Nonce      [4]byte
	TxHashes   [][]byte
}

// NewBlock creates a new Block instance.
func NewBlock(version uint32, prevBlock []byte, merkleRoot []byte, timestamp uint32, bits []byte, nonce []byte, txHashes [][]byte) *Block {
	block := &Block{Version: version, Timestamp: timestamp}
	copy(block.PrevBlock[:32], prevBlock)
	copy(block.MerkleRoot[:32], merkleRoot)
	copy(block.Bits[:4], bits)
	copy(block.Nonce[:4], nonce)
	block.TxHashes = txHashes
	return block
}

// Parse akes a byte stream and parses a block. Returns a Block object
func Parse(s *bytes.Reader) *Block {
	block := new(Block)
	buffer4 := make([]byte, 4)
	buffer32 := make([]byte, 32)
	s.Read(buffer4)
	block.Version = util.LittleEndianToInt32(buffer4)
	s.Read(buffer32)
	copy(block.PrevBlock[:], util.ReverseByteArray(buffer32))
	s.Read(buffer32)
	copy(block.MerkleRoot[:], util.ReverseByteArray(buffer32))
	s.Read(buffer4)
	block.Timestamp = util.LittleEndianToInt32(buffer4)
	s.Read(buffer4)
	copy(block.Bits[:], buffer4)
	s.Read(buffer4)
	copy(block.Nonce[:], buffer4)
	return block
}

// Serialize eturns the 80 byte block header
func (b *Block) Serialize() []byte {
	result := make([]byte, 80)
	position := 0
	copy(result[position:position+4], util.Int32ToLittleEndian(b.Version))
	position += 4
	copy(result[position:position+32], b.PrevBlock[:])
	util.ReverseByteArray(result[position : position+32])
	position += 32
	copy(result[position:position+32], b.MerkleRoot[:])
	util.ReverseByteArray(result[position : position+32])
	position += 32
	copy(result[position:position+4], util.Int32ToLittleEndian(b.Timestamp))
	position += 4
	copy(result[position:position+4], b.Bits[:])
	position += 4
	copy(result[position:position+4], b.Nonce[:])
	return result
}

// Hash returns the hash256 interpreted little endian of the block
func (b *Block) Hash() []byte {
	serialized := b.Serialize()
	return util.ReverseByteArray(util.Hash256(serialized))
}

// Bip9 eturns whether this block is signaling readiness for BIP9
func (b *Block) Bip9() bool {
	// BIP9 is signalled if the top 3 bits are 001
	return b.Version>>29 == 1
}

// Bip91 eturns whether this block is signaling readiness for BIP91
func (b *Block) Bip91() bool {
	// BIP91 is signalled if the 5th bit from the right is 1
	return b.Version&0x10 == 0x10
}

// Bip141 eturns whether this block is signaling readiness for BIP141
func (b *Block) Bip141() bool {
	// BIP141 is signalled if the 2nd bit from the right is 1
	return b.Version&2 == 2
}

// Target eturns the proof-of-work target based on the bits
func (b *Block) Target() *big.Int {
	return util.BitsToTarget(b.Bits[:])
}

// Difficulty returns the block difficulty based on the bits
func (b *Block) Difficulty() *big.Int {
	// note difficulty is (target of lowest difficulty) / (b's target)
	// lowest difficulty has bits that equal 0xffff001d
	lowest := util.BitsToTarget(util.HexStringToBytes("ffff001d"))
	result := new(big.Int)
	result.Div(lowest, b.Target())
	return result
}

// CheckPow returns whether this block satisfies proof of work
func (b *Block) CheckPow() bool {
	sha := b.Hash()
	proof := new(big.Int)
	proof.SetBytes(sha)
	return proof.Cmp(b.Target()) < 0
}

// ValidateMerkleRoot gets the merkle root of the tx hashes and checks that it's the same as the merkle root of this block.
func (b *Block) ValidateMerkleRoot() bool {
	if b.TxHashes == nil {
		return false
	}
	numHashes := len(b.TxHashes)
	hashes := make([][]byte, numHashes)
	// Reverse each item in b.TxHashses
	for i, hash := range b.TxHashes {
		hashes[i] = make([]byte, len(hash))
		copy(hashes[i], hash)
		util.ReverseByteArray(hashes[i])
	}
	root := util.MerkleRoot(hashes)
	util.ReverseByteArray(root)
	return bytes.Equal(root, b.MerkleRoot[:])
}
