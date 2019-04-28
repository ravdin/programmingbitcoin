package util

import (
	"bytes"
	"encoding/hex"
	"testing"
)

func TestUtil(t *testing.T) {
	t.Run("little endian to int", func(t *testing.T) {
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
	})

	t.Run("int to little endian", func(t *testing.T) {
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
	})

	t.Run("base 58", func(t *testing.T) {
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
	})

	t.Run("p2pkh address", func(t *testing.T) {
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
	})

	t.Run("p2sh address", func(t *testing.T) {
		h160, _ := hex.DecodeString("74d691da1574e6b3c192ecfb52cc8984ee7b6c56")
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
	})
}
