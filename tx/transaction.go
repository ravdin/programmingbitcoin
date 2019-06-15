package tx

import (
	"bytes"
	"encoding/hex"
	"math/big"

	"github.com/ravdin/programmingbitcoin/ecc"
	"github.com/ravdin/programmingbitcoin/script"
	"github.com/ravdin/programmingbitcoin/util"
)

// Transaction represents a bitcoin transaction.
type Transaction struct {
	Version  uint32
	Inputs   []*Input
	Outputs  []*Output
	Locktime uint32
	Testnet  bool
}

// NewTransaction initializes a Transaction object.
func NewTransaction(version uint32, txIns []*Input, txOuts []*Output, locktime uint32, testnet bool) *Transaction {
	return &Transaction{
		Version:  version,
		Inputs:   txIns,
		Outputs:  txOuts,
		Locktime: locktime,
		Testnet:  testnet,
	}
}

// ID is a human-readable hexadecimal of the transaction hash
func (tx *Transaction) ID() string {
	return hex.EncodeToString(tx.Hash())
}

// Hash of the legacy serialization
func (tx *Transaction) Hash() []byte {
	hash := util.Hash256(tx.Serialize())
	// reverse the array
	return util.ReverseByteArray(hash)
}

// Serialize the transaction.
func (tx *Transaction) Serialize() []byte {
	result := util.Int32ToLittleEndian(tx.Version)
	result = append(result, util.EncodeVarInt(len(tx.Inputs))...)
	for _, txIn := range tx.Inputs {
		result = append(result, txIn.Serialize()...)
	}
	result = append(result, util.EncodeVarInt(len(tx.Outputs))...)
	for _, txOut := range tx.Outputs {
		result = append(result, txOut.Serialize()...)
	}
	result = append(result, util.Int32ToLittleEndian(tx.Locktime)...)
	return result
}

// ParseTransaction parses a transaction from a byte reader.
func ParseTransaction(s *bytes.Reader, testnet bool) *Transaction {
	buffer := make([]byte, 4)
	s.Read(buffer)
	version := util.LittleEndianToInt32(buffer)
	numInputs := util.ReadVarInt(s)
	inputs := make([]*Input, numInputs)
	for i := 0; i < numInputs; i++ {
		inputs[i] = ParseInput(s)
	}
	numOutputs := util.ReadVarInt(s)
	outputs := make([]*Output, numOutputs)
	for i := 0; i < numOutputs; i++ {
		outputs[i] = ParseOutput(s)
	}
	s.Read(buffer)
	locktime := util.LittleEndianToInt32(buffer)
	return NewTransaction(version, inputs, outputs, locktime, testnet)
}

// Fee returns the fee of this transaction in satoshi
func (tx *Transaction) Fee() uint64 {
	// initialize input sum and output sum
	// use Input.value() to sum up the input amounts
	// use Output.amount to sum up the output amounts
	// fee is input sum - output sum
	var result uint64
	for _, txIn := range tx.Inputs {
		result += txIn.Value(false)
	}
	for _, txOut := range tx.Outputs {
		result -= txOut.Amount
	}
	return result
}

// SigHash rturns the integer representation of the hash that needs to get
// signed for index inputIndex
func (tx *Transaction) SigHash(inputIndex int) []byte {
	serialized := util.Int32ToLittleEndian(tx.Version)
	serialized = append(serialized, util.EncodeVarInt(len(tx.Inputs))...)
	for i, txIn := range tx.Inputs {
		var scriptSig *script.Script
		if i == inputIndex {
			scriptSig = txIn.ScriptPubKey(tx.Testnet)
		}
		serialized = append(serialized, NewInput(
			txIn.PrevTx,
			txIn.PrevIndex,
			scriptSig,
			txIn.Sequence,
		).Serialize()...)
	}
	serialized = append(serialized, util.EncodeVarInt(len(tx.Outputs))...)
	for _, txOut := range tx.Outputs {
		serialized = append(serialized, txOut.Serialize()...)
	}
	serialized = append(serialized, util.Int32ToLittleEndian(tx.Locktime)...)
	serialized = append(serialized, util.Int32ToLittleEndian(util.SIGHASH_ALL)...)
	return util.Hash256(serialized)
}

// Returns whether the input has a valid signature
func (tx *Transaction) verifyInput(inputIndex int) bool {
	txIn := tx.Inputs[inputIndex]
	scriptPubKey := txIn.ScriptPubKey(tx.Testnet)
	z := tx.SigHash(inputIndex)
	// combine the current ScriptSig and the previous ScriptPubKey
	script := new(script.Script)
	script.Add(txIn.ScriptSig, scriptPubKey)
	// evaluate the combined script
	return script.Evaluate(z)
}

// Verify this transaction
func (tx *Transaction) Verify() bool {
	if tx.Fee() < 0 {
		return false
	}
	for i := range tx.Inputs {
		if !tx.verifyInput(i) {
			return false
		}
	}
	return true
}

// SignInput signs a transaction input with a private key.
func (tx *Transaction) SignInput(inputIndex int, pk *ecc.PrivateKey) bool {
	z := new(big.Int)
	z.SetBytes(tx.SigHash(inputIndex))
	// get der signature of z from private key
	der := pk.Sign(z).Der()
	der = append(der, byte(util.SIGHASH_ALL))
	// calculate the sec
	sec := pk.Point.Sec(true)
	// initialize a new script with [sig, sec] as the cmds
	script := script.NewScript([][]byte{der, sec})
	// change input's scriptSig to new script
	tx.Inputs[inputIndex].ScriptSig = script
	// return whether sig is valid using tx.verifyInput
	return tx.verifyInput(inputIndex)
}

// IsCoinbase returns whether this transaction is a coinbase transaction or not
func (tx *Transaction) IsCoinbase() bool {
	if len(tx.Inputs) != 1 {
		return false
	}
	txIn := tx.Inputs[0]
	return bytes.Equal(txIn.PrevTx, make([]byte, 32)) && txIn.PrevIndex == 0xffffffff
}

// Returns the height of the block this coinbase transaction is in
func (tx *Transaction) coinbaseHeight() *big.Int {
	if !tx.IsCoinbase() {
		return nil
	}
	cmd := tx.Inputs[0].ScriptSig.Peek(0)
	return util.LittleEndianToBigInt(cmd)
}
