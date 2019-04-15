package tx

import (
	"bytes"
	"github.com/ravdin/programmingbitcoin/util"
)

type Tx struct {
	Version  uint32
	TxIns    []*TxIn
	TxOuts   []*TxOut
	Locktime uint32
	Testnet  bool
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
