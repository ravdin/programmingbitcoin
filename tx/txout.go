package tx

import (
	"bytes"
	"github.com/ravdin/programmingbitcoin/script"
	"github.com/ravdin/programmingbitcoin/util"
)

type TxOut struct {
	Amount       uint64
	ScriptPubKey *script.Script
}

// Takes a byte stream and parses the tx_output at the start
// return a TxOut object
func ParseTxOut(s *bytes.Reader) *TxOut {
	// amount is an integer in 8 bytes, little endian
	// use Script.parse to get the ScriptPubKey
	// return an instance of the class
	buffer := make([]byte, 8)
	s.Read(buffer)
	amount := util.LittleEndianToInt64(buffer)
	scriptPubkey := script.Parse(s)
	return &TxOut{Amount: amount, ScriptPubKey: scriptPubkey}
}

// Returns the byte serialization of the transaction output
func (self *TxOut) Serialize() []byte {
	result := util.Int64ToLittleEndian(self.Amount)
	result = append(result, self.ScriptPubKey.Serialize()...)
	return result
}
