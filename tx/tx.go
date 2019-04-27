package tx

import (
	"bytes"
	"encoding/hex"
	"math/big"

	"github.com/ravdin/programmingbitcoin/ecc"
	"github.com/ravdin/programmingbitcoin/script"
	"github.com/ravdin/programmingbitcoin/util"
)

type Tx struct {
	Version  uint32
	TxIns    []*TxIn
	TxOuts   []*TxOut
	Locktime uint32
	Testnet  bool
}

func NewTx(version uint32, txIns []*TxIn, txOuts []*TxOut, locktime uint32, testnet bool) *Tx {
	return &Tx{Version: version, TxIns: txIns, TxOuts: txOuts, Locktime: locktime, Testnet: testnet}
}

// Human-readable hexadecimal of the transaction hash
func (self *Tx) Id() string {
	return hex.EncodeToString(self.Hash())
}

// Binary hash of the legacy serialization
func (self *Tx) Hash() []byte {
	hash := util.Hash256(self.Serialize())
	// reverse the array
	return util.ReverseByteArray(hash)
}

// Returns the byte serialization of the transaction
func (self *Tx) Serialize() []byte {
	result := util.Int32ToLittleEndian(self.Version)
	result = append(result, util.EncodeVarInt(len(self.TxIns))...)
	for _, txIn := range self.TxIns {
		result = append(result, txIn.Serialize()...)
	}
	result = append(result, util.EncodeVarInt(len(self.TxOuts))...)
	for _, txOut := range self.TxOuts {
		result = append(result, txOut.Serialize()...)
	}
	result = append(result, util.Int32ToLittleEndian(self.Locktime)...)
	return result
}

func ParseTx(s *bytes.Reader, testnet bool) *Tx {
	// s.read(n) will return n bytes
	// version is an integer in 4 bytes, little-endian
	// num_inputs is a varint, use read_varint(s)
	// parse num_inputs number of TxIns
	// num_outputs is a varint, use read_varint(s)
	// parse num_outputs number of TxOuts
	// locktime is an integer in 4 bytes, little-endian
	// return an instance of the class
	buffer := make([]byte, 4)
	s.Read(buffer)
	version := util.LittleEndianToInt32(buffer)
	num_inputs := util.ReadVarInt(s)
	tx_ins := make([]*TxIn, num_inputs)
	for i := 0; i < num_inputs; i++ {
		tx_ins[i] = ParseTxIn(s)
	}
	num_outputs := util.ReadVarInt(s)
	tx_outs := make([]*TxOut, num_outputs)
	for i := 0; i < num_outputs; i++ {
		tx_outs[i] = ParseTxOut(s)
	}
	s.Read(buffer)
	locktime := util.LittleEndianToInt32(buffer)
	return NewTx(version, tx_ins, tx_outs, locktime, testnet)
}

// Returns the fee of this transaction in satoshi
func (self *Tx) Fee() uint64 {
	// initialize input sum and output sum
	// use TxIn.value() to sum up the input amounts
	// use TxOut.amount to sum up the output amounts
	// fee is input sum - output sum
	var result uint64 = 0
	for _, txIn := range self.TxIns {
		result += txIn.Value(false)
	}
	for _, txOut := range self.TxOuts {
		result -= txOut.Amount
	}
	return result
}

// Returns the integer representation of the hash that needs to get
// signed for index inputIndex
func (self *Tx) SigHash(inputIndex int) []byte {
	serialized := util.Int32ToLittleEndian(self.Version)
	serialized = append(serialized, util.EncodeVarInt(len(self.TxIns))...)
	for i, txIn := range self.TxIns {
		var scriptSig *script.Script = nil
		if i == inputIndex {
			scriptSig = txIn.ScriptPubKey(self.Testnet)
		}
		serialized = append(serialized, NewTxIn(
			txIn.PrevTx,
			txIn.PrevIndex,
			scriptSig,
			txIn.Sequence,
		).Serialize()...)
	}
	serialized = append(serialized, util.EncodeVarInt(len(self.TxOuts))...)
	for _, txOut := range self.TxOuts {
		serialized = append(serialized, txOut.Serialize()...)
	}
	serialized = append(serialized, util.Int32ToLittleEndian(self.Locktime)...)
	serialized = append(serialized, util.Int32ToLittleEndian(util.SIGHASH_ALL)...)
	return util.Hash256(serialized)
}

// Returns whether the input has a valid signature
func (self *Tx) verifyInput(inputIndex int) bool {
	txIn := self.TxIns[inputIndex]
	scriptPubKey := txIn.ScriptPubKey(self.Testnet)
	z := self.SigHash(inputIndex)
	// combine the current ScriptSig and the previous ScriptPubKey
	script := new(script.Script)
	script.Add(txIn.ScriptSig, scriptPubKey)
	// evaluate the combined script
	return script.Evaluate(z)
}

// Verify this transaction
func (self *Tx) Verify() bool {
	if self.Fee() < 0 {
		return false
	}
	for i := range self.TxIns {
		if !self.verifyInput(i) {
			return false
		}
	}
	return true
}

func (self *Tx) SignInput(inputIndex int, pk *ecc.PrivateKey) bool {
	z := new(big.Int)
	z.SetBytes(self.SigHash(inputIndex))
	// get der signature of z from private key
	der := pk.Sign(z).Der()
	der = append(der, byte(util.SIGHASH_ALL))
	// calculate the sec
	sec := pk.Point.Sec(true)
	// initialize a new script with [sig, sec] as the cmds
	script := script.NewScript([][]byte{der, sec})
	// change input's script_sig to new script
	self.TxIns[inputIndex].ScriptSig = script
	// return whether sig is valid using self.verify_input
	return self.verifyInput(inputIndex)
}
