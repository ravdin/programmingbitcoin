package tx

import (
	"bytes"
	"encoding/hex"
	"github.com/ravdin/programmingbitcoin/script"
	"github.com/ravdin/programmingbitcoin/util"
)

type TxIn struct {
	PrevTx    []byte
	PrevIndex int
	ScriptSig *script.Script
	Sequence  uint32
}

func NewTxIn(prevTx []byte, prevIndex int, scriptSig *script.Script, sequence uint32) *TxIn {
	result := &TxIn{PrevTx: prevTx, PrevIndex: prevIndex, Sequence: sequence}
	if scriptSig == nil {
		result.ScriptSig = new(script.Script)
	} else {
		result.ScriptSig = scriptSig
	}
	return result
}

// Takes a byte stream and parses the tx_input at the start
// return a TxIn object
func ParseTxIn(s *bytes.Reader) *TxIn {
	// prev_tx is 32 bytes, little endian
	// prev_index is an integer in 4 bytes, little endian
	// use Script.parse to get the ScriptSig
	// sequence is an integer in 4 bytes, little-endian
	// return an instance of the class
	prevTx := make([]byte, 32)
	s.Read(prevTx)
	util.ReverseByteArray(prevTx)
	buffer := make([]byte, 4)
	s.Read(buffer)
	prevIndex := int(util.LittleEndianToInt32(buffer))
	scriptSig := script.Parse(s)
	s.Read(buffer)
	sequence := util.LittleEndianToInt32(buffer)
	return &TxIn{PrevTx: prevTx, PrevIndex: prevIndex, ScriptSig: scriptSig, Sequence: sequence}
}

// Returns the byte serialization of the transaction input
func (self *TxIn) Serialize() []byte {
	prevTx := make([]byte, len(self.PrevTx))
	copy(prevTx, self.PrevTx)
	util.ReverseByteArray(prevTx)
	prevIndex := util.Int32ToLittleEndian(uint32(self.PrevIndex))
	scriptSig := self.ScriptSig.Serialize()
	sequence := util.Int32ToLittleEndian(self.Sequence)
	result := make([]byte, len(prevTx)+len(prevIndex)+len(scriptSig)+len(sequence))
	index := 0
	for _, item := range [][]byte{prevTx, prevIndex, scriptSig, sequence} {
		copy(result[index:], item)
		index += len(item)
	}
	return result
}

// Get the output value by looking up the tx hash
// Returns the amount in satoshi
func (self *TxIn) Value(testnet bool) uint64 {
	tx := self.fetchTx(testnet)
	return tx.TxOuts[self.PrevIndex].Amount
}

// Get the ScriptPubKey by looking up the tx hash
// Returns a Script object
func (self *TxIn) ScriptPubKey(testnet bool) *script.Script {
	tx := self.fetchTx(testnet)
	return tx.TxOuts[self.PrevIndex].ScriptPubKey
}

func (self *TxIn) fetchTx(testnet bool) *Tx {
	fetcher := NewTxFetcher()
	return fetcher.fetch(hex.EncodeToString(self.PrevTx), testnet, false)
}
