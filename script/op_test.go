package script

import (
	"encoding/hex"
	"github.com/ravdin/programmingbitcoin/util"
	"testing"
)

func TestOp(t *testing.T) {
	t.Run("Test OpHash160", func(t *testing.T) {
		stack := [][]byte{
			[]byte(`hello world`),
		}
		if !OpHash160(&stack) {
			t.Errorf("OpHash160 failed!")
		}
		expected := `d7d5ee7824ff93f94c3055af9382c86c68b5ca92`
		actual := hex.EncodeToString(stack[0])
		if actual != expected {
			t.Errorf("Expected %v, got %v", expected, actual)
		}
	})

	t.Run("Test OpCheckSig", func(t *testing.T) {
		z := util.HexStringToBigInt(`7c076ff316692a3d7eb3c3bb0f8b1488cf72e1afcd929e29307032997a838a3d`)
		sec := util.HexStringToBytes(`04887387e452b8eacc4acfde10d9aaf7f6d9a0f975aabb10d006e4da568744d06c61de6d95231cd89026e286df3b6ae4a894a3378e393e93a0f45b666329a0ae34`)
		sig := util.HexStringToBytes(`3045022000eff69ef2b1bd93a66ed5219add4fb51e11a840f404876325a1e8ffe0529a2c022100c7207fee197d27c618aea621406f6bf5ef6fca38681d82b2f06fddbdce6feab601`)
		stack := [][]byte{sec, sig}
		if !OpCheckSig(&stack, z) {
			t.Errorf("OpCheckSig failed!")
		}
		actual := decodeNum(stack[0])
		expected := 1
		if actual != expected {
			t.Errorf("Expected %v, got %v", expected, actual)
		}
	})
}
