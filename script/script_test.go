package script

import (
	"bytes"
	"encoding/hex"
	"github.com/ravdin/programmingbitcoin/util"
	"testing"
)

func TestScript(t *testing.T) {
	scriptPubKey := `6a47304402207899531a52d59a6de200179928ca900254a36b8dff8bb75f5f5d71b1cdc26125022008b422690b8461cb52c3cc30330b23d574351872b7c361e9aae3649071c1a7160121035d5c93d9ac96881f19ba1f686f15f009ded7c62efe85a872e6a19b43c15a2937`
	reader := bytes.NewReader(util.HexStringToBytes(scriptPubKey))
	s := Parse(reader)

	t.Run("Test cmds stack", func(t *testing.T) {
		cmds := []string{
			`304402207899531a52d59a6de200179928ca900254a36b8dff8bb75f5f5d71b1cdc26125022008b422690b8461cb52c3cc30330b23d574351872b7c361e9aae3649071c1a71601`,
			`035d5c93d9ac96881f19ba1f686f15f009ded7c62efe85a872e6a19b43c15a2937`,
		}
		for i, expected := range cmds {
			actual := hex.EncodeToString(s.cmds[i])
			if actual != expected {
				t.Errorf("Expected %v, got %v", expected, actual)
			}
		}
	})

	t.Run("Test serialize", func(t *testing.T) {
		serialized := s.Serialize()
		if serialized != scriptPubKey {
			t.Errorf("Expected %v, got %v", scriptPubKey, serialized)
		}
	})
}
