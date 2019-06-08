package util

// murmur3 hash
func Murmur3(data []byte, seed uint32) uint32 {
	var c1 uint32 = 0xcc9e2d51
	var c2 uint32 = 0x1b873593
	length := len(data)
	h1 := seed
	roundedEnd := (length & 0xfffffffc) // round down to 4 byte block
	for i := 0; i < roundedEnd; i += 4 {
		// little endian load order
		k1 := LittleEndianToInt32(data[i : i+4])
		k1 *= c1
		k1 = (k1 << 15) | (k1 >> 17) // ROTL32(k1,15)
		k1 *= c2
		h1 ^= k1
		h1 = (h1 << 13) | (h1 >> 19) // ROTL32(h1,13)
		h1 = h1*5 + 0xe6546b64
	}
	// tail
	var k1 uint32 = 0
	val := length & 0x03
	if val == 3 {
		k1 = uint32(data[roundedEnd+2]) << 16
	}
	// fallthrough
	if val >= 2 {
		k1 |= uint32(data[roundedEnd+1]) << 8
	}
	// fallthrough
	if val >= 1 {
		k1 |= uint32(data[roundedEnd])
		k1 *= c1
		k1 = (k1 << 15) | (k1 >> 17) // ROTL32(k1,15)
		k1 *= c2
		h1 ^= k1
	}
	// finalization
	h1 ^= uint32(length)
	// fmix(h1)
	h1 ^= (h1 >> 16)
	h1 *= 0x85ebca6b
	h1 ^= (h1 >> 13)
	h1 *= 0xc2b2ae35
	h1 ^= (h1 >> 16)
	return h1
}
