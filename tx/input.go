package tx

import (
	"bytes"
	"encoding/hex"

	"github.com/ravdin/programmingbitcoin/script"
	"github.com/ravdin/programmingbitcoin/util"
)

// Input represents a transaction input.
type Input struct {
	PrevTx    []byte
	PrevIndex int
	ScriptSig *script.Script
	Sequence  uint32
}

// NewInput initializes a new transaction input.
func NewInput(prevTx []byte, prevIndex int, scriptSig *script.Script, sequence uint32) *Input {
	result := &Input{PrevTx: prevTx, PrevIndex: prevIndex, Sequence: sequence}
	if scriptSig == nil {
		result.ScriptSig = new(script.Script)
	} else {
		result.ScriptSig = scriptSig
	}
	return result
}

// ParseInput takes a byte stream and parses the tx_input at the start
// return a Input object
func ParseInput(s *bytes.Reader) *Input {
	prevTx := make([]byte, 32)
	s.Read(prevTx)
	util.ReverseByteArray(prevTx)
	buffer := make([]byte, 4)
	s.Read(buffer)
	prevIndex := int(util.LittleEndianToInt32(buffer))
	scriptSig := script.Parse(s)
	s.Read(buffer)
	sequence := util.LittleEndianToInt32(buffer)
	return NewInput(prevTx, prevIndex, scriptSig, sequence)
}

// Serialize the transaction input
func (in *Input) Serialize() []byte {
	prevTx := make([]byte, len(in.PrevTx))
	copy(prevTx, in.PrevTx)
	util.ReverseByteArray(prevTx)
	prevIndex := util.Int32ToLittleEndian(uint32(in.PrevIndex))
	scriptSig := in.ScriptSig.Serialize()
	sequence := util.Int32ToLittleEndian(in.Sequence)
	result := make([]byte, len(prevTx)+len(prevIndex)+len(scriptSig)+len(sequence))
	index := 0
	for _, item := range [][]byte{prevTx, prevIndex, scriptSig, sequence} {
		copy(result[index:], item)
		index += len(item)
	}
	return result
}

// Value is the output value by looking up the tx hash
// Returns the amount in satoshi
func (in *Input) Value(testnet bool) uint64 {
	tx := in.fetchTx(testnet)
	return tx.Outputs[in.PrevIndex].Amount
}

// ScriptPubKey looks up the tx hash
// Returns a Script object
func (in *Input) ScriptPubKey(testnet bool) *script.Script {
	tx := in.fetchTx(testnet)
	return tx.Outputs[in.PrevIndex].ScriptPubKey
}

func (in *Input) fetchTx(testnet bool) *Transaction {
	fetcher := newTxFetcher()
	return fetcher.fetch(hex.EncodeToString(in.PrevTx), testnet, false)
}
