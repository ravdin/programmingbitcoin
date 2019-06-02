package block

import (
	"bytes"
	"math/big"

	"github.com/ravdin/programmingbitcoin/util"
)

var (
	GENESIS_BLOCK      []byte = util.HexStringToBytes(`0100000000000000000000000000000000000000000000000000000000000000000000003ba3edfd7a7b12b27ac72c3e67768f617fc81bc3888a51323a9fb8aa4b1e5e4a29ab5f49ffff001d1dac2b7c`)
	TEST_GENESIS_BLOCK []byte = util.HexStringToBytes(`0100000000000000000000000000000000000000000000000000000000000000000000003ba3edfd7a7b12b27ac72c3e67768f617fc81bc3888a51323a9fb8aa4b1e5e4adae5494dffff001d1aa4ae18`)
	LOWEST_BITS        []byte = util.HexStringToBytes(`ffff001d`)
)

type Block struct {
	Version    uint32
	PrevBlock  [32]byte
	MerkleRoot [32]byte
	Timestamp  uint32
	Bits       [4]byte
	Nonce      [4]byte
	TxHashes   [][]byte
}

func NewBlock(version uint32, prevBlock []byte, merkleRoot []byte, timestamp uint32, bits []byte, nonce []byte, txHashes [][]byte) *Block {
	block := &Block{Version: version, Timestamp: timestamp}
	copy(block.PrevBlock[:32], prevBlock)
	copy(block.MerkleRoot[:32], merkleRoot)
	copy(block.Bits[:4], bits)
	copy(block.Nonce[:4], nonce)
	block.TxHashes = txHashes
	return block
}

// Takes a byte stream and parses a block. Returns a Block object
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

// Returns the 80 byte block header
func (self *Block) Serialize() []byte {
	result := make([]byte, 80)
	position := 0
	copy(result[position:position+4], util.Int32ToLittleEndian(self.Version))
	position += 4
	copy(result[position:position+32], self.PrevBlock[:])
	util.ReverseByteArray(result[position : position+32])
	position += 32
	copy(result[position:position+32], self.MerkleRoot[:])
	util.ReverseByteArray(result[position : position+32])
	position += 32
	copy(result[position:position+4], util.Int32ToLittleEndian(self.Timestamp))
	position += 4
	copy(result[position:position+4], self.Bits[:])
	position += 4
	copy(result[position:position+4], self.Nonce[:])
	return result
}

// Returns the hash256 interpreted little endian of the block
func (self *Block) Hash() []byte {
	serialized := self.Serialize()
	return util.ReverseByteArray(util.Hash256(serialized))
}

// Returns whether this block is signaling readiness for BIP9
func (self *Block) Bip9() bool {
	// BIP9 is signalled if the top 3 bits are 001
	return self.Version>>29 == 1
}

// Returns whether this block is signaling readiness for BIP91
func (self *Block) Bip91() bool {
	// BIP91 is signalled if the 5th bit from the right is 1
	return self.Version&0x10 == 0x10
}

// Returns whether this block is signaling readiness for BIP141
func (self *Block) Bip141() bool {
	// BIP141 is signalled if the 2nd bit from the right is 1
	return self.Version&2 == 2
}

// Returns the proof-of-work target based on the bits
func (self *Block) Target() *big.Int {
	return util.BitsToTarget(self.Bits[:])
}

// Returns the block difficulty based on the bits
func (self *Block) Difficulty() *big.Int {
	// note difficulty is (target of lowest difficulty) / (self's target)
	// lowest difficulty has bits that equal 0xffff001d
	lowest := util.BitsToTarget(util.HexStringToBytes("ffff001d"))
	result := new(big.Int)
	result.Div(lowest, self.Target())
	return result
}

// Returns whether this block satisfies proof of work
func (self *Block) CheckPow() bool {
	sha := self.Hash()
	proof := new(big.Int)
	proof.SetBytes(sha)
	return proof.Cmp(self.Target()) < 0
}

// Gets the merkle root of the tx hashes and checks that it's the same as the merkle root of this block.
func (self *Block) ValidateMerkleRoot() bool {
	if self.TxHashes == nil {
		return false
	}
	numHashes := len(self.TxHashes)
	hashes := make([][]byte, numHashes)
	// Reverse each item in self.TxHashses
	for i, hash := range self.TxHashes {
		hashes[i] = make([]byte, len(hash))
		copy(hashes[i], hash)
		util.ReverseByteArray(hashes[i])
	}
	root := util.MerkleRoot(hashes)
	util.ReverseByteArray(root)
	return bytes.Equal(root, self.MerkleRoot[:])
}
