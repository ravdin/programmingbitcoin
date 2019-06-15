package bloomfilter

import (
	"github.com/ravdin/programmingbitcoin/network"
	"github.com/ravdin/programmingbitcoin/util"
)

// Bip37Constant is a seed for BIP0037 bloom filters
const Bip37Constant uint32 = 0xfba4c795

// BloomFilter is a filter implementation for all possible transactions
type BloomFilter struct {
	Size          int
	BitField      []byte
	FunctionCount uint32
	Tweak         uint32
}

// New creates a new BloomFilter
func New(size int, functionCount uint32, tweak uint32) *BloomFilter {
	bitField := make([]byte, size*8)
	return &BloomFilter{
		Size:          size,
		BitField:      bitField,
		FunctionCount: functionCount,
		Tweak:         tweak,
	}
}

// Add an item to the filter
func (bf *BloomFilter) Add(item []byte) {
	for i := uint32(0); i < bf.FunctionCount; i++ {
		seed := i*Bip37Constant + bf.Tweak
		h := util.Murmur3(item, seed)
		bit := int(h % uint32(len(bf.BitField)))
		bf.BitField[bit] = 1
	}
}

func (bf *BloomFilter) filterBytes() []byte {
	return util.BitFieldToBytes(bf.BitField)
}

// FilterLoad returns the filterload message
func (bf *BloomFilter) FilterLoad(flag byte) network.Message {
	size := util.EncodeVarInt(bf.Size)
	filter := bf.filterBytes()
	payload := make([]byte, len(size)+len(filter)+9)
	copy(payload, size)
	pos := len(size)
	copy(payload[pos:], filter)
	pos += len(filter)
	copy(payload[pos:], util.Int32ToLittleEndian(bf.FunctionCount))
	copy(payload[pos+4:], util.Int32ToLittleEndian(bf.Tweak))
	payload[pos+8] = flag
	return network.NewGenericMessage([]byte("filterload"), payload)
}
