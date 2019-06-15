package tx

import (
	"bytes"

	"github.com/ravdin/programmingbitcoin/script"
	"github.com/ravdin/programmingbitcoin/util"
)

// Output represents a transaction output.
type Output struct {
	Amount       uint64
	ScriptPubKey *script.Script
}

// NewOutput returns a new transaction output.
// amount: Amount in the output.
// scriptPubKey: Script to verify the output.
func NewOutput(amount uint64, scriptPubkey *script.Script) *Output {
	return &Output{Amount: amount, ScriptPubKey: scriptPubkey}
}

// ParseOutput parses an Output object from a byte reader.
func ParseOutput(s *bytes.Reader) *Output {
	// amount is an integer in 8 bytes, little endian
	// use Script.parse to get the ScriptPubKey
	// return an instance of the class
	buffer := make([]byte, 8)
	s.Read(buffer)
	amount := util.LittleEndianToInt64(buffer)
	scriptPubkey := script.Parse(s)
	return NewOutput(amount, scriptPubkey)
}

// Serialize the transaction output
func (out *Output) Serialize() []byte {
	result := util.Int64ToLittleEndian(out.Amount)
	result = append(result, out.ScriptPubKey.Serialize()...)
	return result
}
