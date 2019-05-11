package network

import (
	"encoding/hex"
	"testing"

	"github.com/ravdin/programmingbitcoin/util"
)

func TestGetHeadersMessage(t *testing.T) {
	blockHex := `0000000000000000001237f46acddf58578a37e213d2a6edc4884a2fcad05ba3`
	ghm := NewGetHeadersMessage(util.HexStringToBytes(blockHex))
	actual := hex.EncodeToString(ghm.Serialize())
	expected := `7f11010001a35bd0ca2f4a88c4eda6d213e2378a5758dfcd6af437120000000000000000000000000000000000000000000000000000000000000000000000000000000000`
	if actual != expected {
		t.Errorf("Expected %s, got %s", expected, actual)
	}
}
