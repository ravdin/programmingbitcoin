package tx

import (
	"bytes"
	"github.com/ravdin/programmingbitcoin/script"
	"github.com/ravdin/programmingbitcoin/util"
)

type TxIn struct {
	PrevTx    []byte
	PrevIndex int32
	ScriptSig *script.Script
	Sequence  int32
}

func ParseTxIn(s *bytes.Reader) *TxIn {
	// prev_tx is 32 bytes, little endian
	// prev_index is an integer in 4 bytes, little endian
	// use Script.parse to get the ScriptSig
	// sequence is an integer in 4 bytes, little-endian
	// return an instance of the class
	buffer := make([]byte, 32)
	s.Read(buffer)
	prev_tx := reverse(buffer)
	buffer = make([]byte, 4)
	s.Read(buffer)
	prev_index := util.LittleEndianToInt32(buffer)
	script_sig := script.Parse(s)
	s.Read(buffer)
	sequence := util.LittleEndianToInt32(buffer)
	return &TxIn{PrevTx: prev_tx, PrevIndex: prev_index, ScriptSig: script_sig, Sequence: sequence}
}

func reverse(b []byte) []byte {
	result := make([]byte, len(b))
	for i := 0; i < len(b); i++ {
		result[len(b)-i-1] = b[i]
	}
	return result
}
