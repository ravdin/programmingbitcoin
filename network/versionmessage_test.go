package network

import (
	"encoding/hex"
	"testing"
)

func TestVersionMessage(t *testing.T) {
	args := map[int]interface{}{
		TIMESTAMP: uint64(0),
		NONCE:     [8]byte{0, 0, 0, 0, 0, 0, 0, 0},
	}
	message := NewVersionMessage(args)
	t.Run("test serialize", func(t *testing.T) {
		actual := hex.EncodeToString(message.Serialize())
		expected := `7f11010000000000000000000000000000000000000000000000000000000000000000000000ffff000000008d20000000000000000000000000000000000000ffff000000008d200000000000000000182f70726f6772616d6d696e67626974636f696e3a302e312f0000000000`
		if actual != expected {
			t.Errorf("Expected %s, got %s", expected, actual)
		}
	})
}
