package tx

import (
	"github.com/ravdin/programmingbitcoin/script"
)

type TxIn struct {
	PrevTx    [32]byte
	PrevIndex int
	ScriptSig *script.Script
	Sequence  int
}
