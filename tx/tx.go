package tx

import (
	"bytes"
	"encoding/hex"
	"github.com/ravdin/programmingbitcoin/ecc"
	"github.com/ravdin/programmingbitcoin/util"
	"math/big"
)

type Tx struct {
	Version  uint32
	TxIns    []*TxIn
	TxOuts   []*TxOut
	Locktime uint32
	Testnet  bool
}

// Human-readable hexadecimal of the transaction hash
func (self *Tx) Id() string {
	return hex.EncodeToString(self.Hash())
}

// Binary hash of the legacy serialization
func (self *Tx) Hash() []byte {
	serialized := self.Serialize()
	// reverse the array
	length := len(serialized)
	for i := 0; i < length/2; i++ {
		serialized[i], serialized[length-i-1] = serialized[length-i-1], serialized[i]
	}
	return util.Hash256(serialized)
}

// Returns the byte serialization of the transaction
func (self *Tx) Serialize() []byte {
	panic("Not implemented")
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
	return &Tx{Version: version, TxIns: tx_ins, TxOuts: tx_outs, Locktime: locktime, Testnet: testnet}
}

func (self *Tx) Fee() uint64 {
	//Returns the fee of this transaction in satoshi'
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
func (self *Tx) SigHash(inputIndex int) *big.Int {
	panic("Not implemented")
}

// Returns whether the input has a valid signature
func (self *Tx) verifyInput(inputIndex int) bool {
	panic("Not implemented")
}

// Verify this transaction
func (self *Tx) Verify() bool {
	panic("Not implemented")
}

func (self *Tx) SignInput(inputIndex int, pk *ecc.PrivateKey) bool {
	panic("Not implemented")
}
