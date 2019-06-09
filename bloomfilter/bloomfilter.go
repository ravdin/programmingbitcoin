package bloomfilter

import (
	"github.com/ravdin/programmingbitcoin/network"
	"github.com/ravdin/programmingbitcoin/util"
)

const BIP37_CONSTANT uint32 = 0xfba4c795

type BloomFilter struct {
	Size          int
	BitField      []byte
	FunctionCount uint32
	Tweak         uint32
}

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
func (self *BloomFilter) Add(item []byte) {
	for i := uint32(0); i < self.FunctionCount; i++ {
		seed := i*BIP37_CONSTANT + self.Tweak
		h := util.Murmur3(item, seed)
		bit := int(h % uint32(len(self.BitField)))
		self.BitField[bit] = 1
	}
}

func (self *BloomFilter) FilterBytes() []byte {
	return util.BitFieldToBytes(self.BitField)
}

func (self *BloomFilter) FilterLoad(flag byte) network.Message {
	size := util.EncodeVarInt(self.Size)
	filter := self.FilterBytes()
	payload := make([]byte, len(size)+len(filter)+9)
	copy(payload, size)
	pos := len(size)
	copy(payload[pos:], filter)
	pos += len(filter)
	copy(payload[pos:], util.Int32ToLittleEndian(self.FunctionCount))
	copy(payload[pos+4:], util.Int32ToLittleEndian(self.Tweak))
	payload[pos+8] = flag
	return network.NewGenericMessage([]byte("filterload"), payload)
}
