package tx

import (
  "bytes"
)

type Tx struct {
	Version    int
	NumInputs  int
	TxIns      []TxIn
	NumOutputs int
	TxOuts     []TxOut
	Locktime   int
}

func ParseTx(s string) *Tx {
  // s.read(n) will return n bytes
  // version is an integer in 4 bytes, little-endian
  // num_inputs is a varint, use read_varint(s)
  // parse num_inputs number of TxIns
  // num_outputs is a varint, use read_varint(s)
  // parse num_outputs number of TxOuts
  // locktime is an integer in 4 bytes, little-endian
  // return an instance of the class (see __init__ for args)
  version = little_endian_to_int(s.read(4))
  num_inputs = read_varint(s)
  tx_ins = [TxIn.parse(s) for _ in range(num_inputs)]
  num_outputs = read_varint(s)
  tx_outs = [TxOut.parse(s) for _ in range(num_outputs)]
  locktime = little_endian_to_int(s.read())
  return cls(version, tx_ins, tx_outs, locktime, testnet)
}
