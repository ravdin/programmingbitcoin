package network

import (
	"bytes"
	"encoding/hex"
	"testing"

	"github.com/ravdin/programmingbitcoin/util"
)

func TestEnvelope(t *testing.T) {
	messages := []string{
		`f9beb4d976657261636b000000000000000000005df6e0e2`,
		`f9beb4d976657273696f6e0000000000650000005f1a69d2721101000100000000000000bc8f5e5400000000010000000000000000000000000000000000ffffc61b6409208d010000000000000000000000000000000000ffffcb0071c0208d128035cbc97953f80f2f5361746f7368693a302e392e332fcf05050001`,
	}
	commands := []string{`verack`, `version`}
	envelopes := make([]*Envelope, 2)
	for i, msg := range messages {
		data := util.HexStringToBytes(msg)
		reader := bytes.NewReader(data)
		envelopes[i] = ParseEnvelope(reader, false)
	}
	t.Run("Test parse", func(t *testing.T) {
		for i, env := range envelopes {
			expected := commands[i]
			actual := string(env.Command)
			if actual != expected {
				t.Errorf("Expected %s, got %s", expected, actual)
			}
		}
	})
	t.Run("Test serialize", func(t *testing.T) {
		for i, env := range envelopes {
			actual := hex.EncodeToString(env.Serialize())
			expected := messages[i]
			if actual != expected {
				t.Errorf("Expected %s, got %s", expected, actual)
			}
		}
	})
}
