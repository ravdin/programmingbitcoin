package tx

import (
  "bytes"
	"github.com/ravdin/programmingbitcoin/util"
)

type Tx struct {
	Version    int32
	TxIns      []*TxIn
	TxOuts     []*TxOut
	Locktime   byte
  Testnet    bool
}

func ParseTx(s bytes.Reader, testnet bool) *Tx {
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
  b, _ := s.ReadByte()
  locktime := util.LittleEndianToByte(b)
  return &Tx{Version: version, TxIns: tx_ins, TxOuts: tx_outs, Locktime: locktime, Testnet: testnet}
}
