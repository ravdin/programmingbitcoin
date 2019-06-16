package util

import (
	"bytes"
	"encoding/hex"
	"testing"
)

func TestLittleEndianToInt(t *testing.T) {
	tests := map[uint64]string{
		10011545: "99c3980000000000",
		32454049: "a135ef0100000000",
		254:      "fe00",
	}
	for expected, h := range tests {
		actual := LittleEndianToInt64(HexStringToBytes(h))
		if actual != expected {
			t.Errorf("Expected %d, got %d", expected, actual)
		}
	}
}

func TestIntToLittleEndian(t *testing.T) {
	expected := []byte{1, 0, 0, 0}
	actual := Int32ToLittleEndian(1)
	if !bytes.Equal(expected, actual) {
		t.Errorf("Expected %v, got %v", expected, actual)
	}
	expected = []byte{0x99, 0xc3, 0x98, 0x00, 0x00, 0x00, 0x00, 0x00}
	actual = Int64ToLittleEndian(10011545)
	if !bytes.Equal(expected, actual) {
		t.Errorf("Expected %v, got %v", expected, actual)
	}
}

func TestBase58(t *testing.T) {
	address := "mnrVtF8DWjMu839VW3rBfgYaAfKk8983Xf"
	h160 := DecodeBase58(address)
	actual := hex.EncodeToString(h160)
	expected := "507b27411ccf7f16f10297de6cef3f291623eddf"
	if actual != expected {
		t.Errorf("Expected %s, got %s", expected, actual)
	}
	b := make([]byte, len(h160)+1)
	b[0] = 0x6f
	copy(b[1:], h160)
	encoded := EncodeBase58Checksum(b)
	if encoded != address {
		t.Errorf("Expected %s, got %s", address, encoded)
	}
}

func Testp2pkhAddress(t *testing.T) {
	h160, _ := hex.DecodeString("74d691da1574e6b3c192ecfb52cc8984ee7b6c56")
	mainnet := "1BenRpVUFK65JFWcQSuHnJKzc4M8ZP8Eqa"
	testnet := "mrAjisaT4LXL5MzE81sfcDYKU3wqWSvf9q"
	actual := H160ToP2pkhAddress(h160, false)
	if actual != mainnet {
		t.Errorf("Expected %s, got %s", mainnet, actual)
	}
	actual = H160ToP2pkhAddress(h160, true)
	if actual != testnet {
		t.Errorf("Expected %s, got %s", testnet, actual)
	}
}

func Testp2shAddress(t *testing.T) {
	h160 := HexStringToBytes("74d691da1574e6b3c192ecfb52cc8984ee7b6c56")
	mainnet := "3CLoMMyuoDQTPRD3XYZtCvgvkadrAdvdXh"
	testnet := "2N3u1R6uwQfuobCqbCgBkpsgBxvr1tZpe7B"
	actual := H160ToP2shAddress(h160, false)
	if actual != mainnet {
		t.Errorf("Expected %s, got %s", mainnet, actual)
	}
	actual = H160ToP2shAddress(h160, true)
	if actual != testnet {
		t.Errorf("Expected %s, got %s", testnet, actual)
	}
}

func TestCalculateNewBits(t *testing.T) {
	prevBits := HexStringToBytes("54d80118")
	timeDifferential := 302400
	actual := CalculateNewBits(prevBits, timeDifferential)
	expected := HexStringToBytes("00157617")
	if !bytes.Equal(actual, expected) {
		t.Errorf("Expected %x, got %x", expected, actual)
	}
}

func TestMerkleParent(t *testing.T) {
	hash1 := HexStringToBytes(`c117ea8ec828342f4dfb0ad6bd140e03a50720ece40169ee38bdc15d9eb64cf5`)
	hash2 := HexStringToBytes(`c131474164b412e3406696da1ee20ab0fc9bf41c8f05fa8ceea7a08d672d7cc5`)
	actual := MerkleParent(hash1, hash2)
	expected := HexStringToBytes(`8b30c5ba100f6f2e5ad1e2a742e5020491240f8eb514fe97c713c31718ad7ecd`)
	if !bytes.Equal(actual, expected) {
		t.Errorf("Expected %x, got %x", expected, actual)
	}
}

func TestMerkle(t *testing.T) {
	testHashes := [][]byte{
		HexStringToBytes(`c117ea8ec828342f4dfb0ad6bd140e03a50720ece40169ee38bdc15d9eb64cf5`),
		HexStringToBytes(`c131474164b412e3406696da1ee20ab0fc9bf41c8f05fa8ceea7a08d672d7cc5`),
		HexStringToBytes(`f391da6ecfeed1814efae39e7fcb3838ae0b02c02ae7d0a5848a66947c0727b0`),
		HexStringToBytes(`3d238a92a94532b946c90e19c49351c763696cff3db400485b813aecb8a13181`),
		HexStringToBytes(`10092f2633be5f3ce349bf9ddbde36caa3dd10dfa0ec8106bce23acbff637dae`),
		HexStringToBytes(`7d37b3d54fa6a64869084bfd2e831309118b9e833610e6228adacdbd1b4ba161`),
		HexStringToBytes(`8118a77e542892fe15ae3fc771a4abfd2f5d5d5997544c3487ac36b5c85170fc`),
		HexStringToBytes(`dff6879848c2c9b62fe652720b8df5272093acfaa45a43cdb3696fe2466a3877`),
		HexStringToBytes(`b825c0745f46ac58f7d3759e6dc535a1fec7820377f24d4c2c6ad2cc55c0cb59`),
		HexStringToBytes(`95513952a04bd8992721e9b7e2937f1c04ba31e0469fbe615a78197f68f52b7c`),
		HexStringToBytes(`2e6d722e5e4dbdf2447ddecc9f7dabb8e299bae921c99ad5b0184cd9eb8e5908`),
	}

	t.Run("test merkle parent level", func(t *testing.T) {
		actual := MerkleParentLevel(testHashes)
		expected := [][]byte{
			HexStringToBytes(`8b30c5ba100f6f2e5ad1e2a742e5020491240f8eb514fe97c713c31718ad7ecd`),
			HexStringToBytes(`7f4e6f9e224e20fda0ae4c44114237f97cd35aca38d83081c9bfd41feb907800`),
			HexStringToBytes(`ade48f2bbb57318cc79f3a8678febaa827599c509dce5940602e54c7733332e7`),
			HexStringToBytes(`68b3e2ab8182dfd646f13fdf01c335cf32476482d963f5cd94e934e6b3401069`),
			HexStringToBytes(`43e7274e77fbe8e5a42a8fb58f7decdb04d521f319f332d88e6b06f8e6c09e27`),
			HexStringToBytes(`1796cd3ca4fef00236e07b723d3ed88e1ac433acaaa21da64c4b33c946cf3d10`),
		}
		if len(actual) != len(expected) {
			t.Errorf("Expected length %d, got %d", len(expected), len(actual))
		}
		for i, item := range actual {
			if !bytes.Equal(item, expected[i]) {
				t.Errorf("Expected %x, got %x", expected[i], item)
			}
		}
	})

	t.Run("test merkle root", func(t *testing.T) {
		expected := HexStringToBytes(`acbcab8bcc1af95d8d563b77d24c3d19b18f1486383d75a5085c4e86c86beed6`)
		actual := MerkleRoot(append(testHashes, HexStringToBytes(`b13a750047bc0bdceb2473e5fe488c2596d7a7124b4e716fdd29b046ef99bbf0`)))
		if !bytes.Equal(actual, expected) {
			t.Errorf("Expected %x, got %x", expected, actual)
		}
	})

	t.Run("test bit field to bytes", func(t *testing.T) {
		bitfield := []byte{0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 0, 0, 1, 0, 1, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0}
		expected := HexStringToBytes(`4000600a080000010940`)
		actual := BitFieldToBytes(bitfield)
		if !bytes.Equal(actual, expected) {
			t.Errorf("Expected %x, got %x", expected, actual)
		}
		expected = BytesToBitField(expected)
		if !bytes.Equal(bitfield, expected) {
			t.Errorf("Expected %v, got %v", expected, bitfield)
		}
	})
}
