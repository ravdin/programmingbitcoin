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

func ParseTxOut(s *bytes.Reader) *TxOut {
	// amount is an integer in 8 bytes, little endian
	// use Script.parse to get the ScriptPubKey
	// return an instance of the class
	buffer := make([]byte, 8)
	s.Read(buffer)
	amount := util.LittleEndianToInt64(buffer)
	script_pubkey := script.Parse(s)
	return &TxOut{Amount: amount, ScriptPubKey: script_pubkey}
}
