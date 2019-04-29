package block

import (
	"bytes"
	"encoding/hex"
	"math/big"
	"testing"

	"github.com/ravdin/programmingbitcoin/util"
)

func TestBlock(t *testing.T) {
	serialized := `020000208ec39428b17323fa0ddec8e887b4a7c53b8c0a0a220cfd0000000000000000005b0750fce0a889502d40508d39576821155e9c9e3f5c3157f961db38fd8b25be1e77a759e93c0118a4ffd71d`
	block := parseBlockFromString(serialized)

	t.Run("test parse", func(t *testing.T) {
		var expectedVersion uint32 = 0x20000002
		expectedPrevBlock := util.HexStringToBytes(`000000000000000000fd0c220a0a8c3bc5a7b487e8c8de0dfa2373b12894c38e`)
		expectedMerkleRoot := util.HexStringToBytes(`be258bfd38db61f957315c3f9e9c5e15216857398d50402d5089a8e0fc50075b`)
		var expectedTimestamp uint32 = 0x59a7771e
		expectedBits := util.HexStringToBytes(`e93c0118`)
		expectedNonce := util.HexStringToBytes(`a4ffd71d`)
		if block.Version != expectedVersion {
			t.Errorf("Expected %d, got %d", expectedVersion, block.Version)
		}
		if !bytes.Equal(expectedPrevBlock, block.PrevBlock[:]) {
			t.Errorf("Expected %v, got %v", expectedPrevBlock, block.PrevBlock)
		}
		if !bytes.Equal(expectedMerkleRoot, block.MerkleRoot[:]) {
			t.Errorf("Expected %v, got %v", expectedMerkleRoot, block.MerkleRoot)
		}
		if block.Timestamp != expectedTimestamp {
			t.Errorf("Expected %d, got %d", expectedTimestamp, block.Timestamp)
		}
		if !bytes.Equal(expectedBits, block.Bits[:]) {
			t.Errorf("Expected %v, got %v", expectedBits, block.Bits)
		}
		if !bytes.Equal(expectedNonce, block.Nonce[:]) {
			t.Errorf("Expected %v, got %v", expectedNonce, block.Nonce)
		}
	})

	t.Run("test serialize", func(t *testing.T) {
		actual := hex.EncodeToString(block.Serialize())
		if actual != serialized {
			t.Errorf("Expected %s, got %s", serialized, actual)
		}
	})

	t.Run("test hash", func(t *testing.T) {
		actual := hex.EncodeToString(block.Hash())
		expected := `0000000000000000007e9e4c586439b0cdbe13b1370bdd9435d76a644d047523`
		if actual != expected {
			t.Errorf("Expected %s, got %s", expected, actual)
		}
	})

	t.Run("test bip9", func(t *testing.T) {
		if !block.Bip9() {
			t.Errorf("Expected true")
		}
		block2 := parseBlockFromString(`0400000039fa821848781f027a2e6dfabbf6bda920d9ae61b63400030000000000000000ecae536a304042e3154be0e3e9a8220e5568c3433a9ab49ac4cbb74f8df8e8b0cc2acf569fb9061806652c27`)
		if block2.Bip9() {
			t.Errorf("Expected false")
		}
	})

	t.Run("test bip91", func(t *testing.T) {
		if block.Bip91() {
			t.Errorf("Expected false")
		}
		block2 := parseBlockFromString(`1200002028856ec5bca29cf76980d368b0a163a0bb81fc192951270100000000000000003288f32a2831833c31a25401c52093eb545d28157e200a64b21b3ae8f21c507401877b5935470118144dbfd1`)
		if !block2.Bip91() {
			t.Errorf("Expected true")
		}
	})

	t.Run("test bip141", func(t *testing.T) {
		if !block.Bip141() {
			t.Errorf("Expected true")
		}
		block2 := parseBlockFromString(`0000002066f09203c1cf5ef1531f24ed21b1915ae9abeb691f0d2e0100000000000000003de0976428ce56125351bae62c5b8b8c79d8297c702ea05d60feabb4ed188b59c36fa759e93c0118b74b2618`)
		if block2.Bip141() {
			t.Errorf("Expected false")
		}
	})

	t.Run("test target", func(t *testing.T) {
		actual := block.Target()
		expected := util.HexStringToBigInt(`13ce9000000000000000000000000000000000000000000`)
		if actual.Cmp(expected) != 0 {
			t.Errorf("Expected %v, got %v", expected, actual)
		}
	})

	t.Run("test difficulty", func(t *testing.T) {
		actual := block.Difficulty()
		expected := big.NewInt(888171856257)
		if actual.Cmp(expected) != 0 {
			t.Errorf("Expected %v, got %v", expected, actual)
		}
	})

	t.Run("test CheckPow", func(t *testing.T) {
		valid := parseBlockFromString(`04000000fbedbbf0cfdaf278c094f187f2eb987c86a199da22bbb20400000000000000007b7697b29129648fa08b4bcd13c9d5e60abb973a1efac9c8d573c71c807c56c3d6213557faa80518c3737ec1`)
		invalid := parseBlockFromString(`04000000fbedbbf0cfdaf278c094f187f2eb987c86a199da22bbb20400000000000000007b7697b29129648fa08b4bcd13c9d5e60abb973a1efac9c8d573c71c807c56c3d6213557faa80518c3737ec0`)
		if !valid.CheckPow() {
			t.Errorf("Expected true")
		}
		if invalid.CheckPow() {
			t.Errorf("Expected false")
		}
	})
}

func parseBlockFromString(str string) *Block {
	raw := util.HexStringToBytes(str)
	reader := bytes.NewReader(raw)
	return Parse(reader)
}
