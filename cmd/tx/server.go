package main

import (
	"encoding/hex"
	"encoding/json"
	"log"
	"net/http"

	"github.com/ravdin/programmingbitcoin/ecc"
	"github.com/ravdin/programmingbitcoin/script"
	"github.com/ravdin/programmingbitcoin/tx"
	"github.com/ravdin/programmingbitcoin/util"
)

func createTx(rw http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	var t transaction
	err := decoder.Decode(&t)
	if err != nil {
		panic(err)
	}
	txIns := make([]*tx.Input, len(t.Inputs))
	for i, input := range t.Inputs {
		prevTx := util.HexStringToBytes(input.PreviousValue)
		txIns[i] = tx.NewInput(prevTx, input.PreviousIndex, nil, 0xffffffff)
	}

	txOuts := make([]*tx.Output, len(t.Outputs))
	for i, output := range t.Outputs {
		script := script.P2pkhScript(util.DecodeBase58(output.Address))
		txOuts[i] = tx.NewOutput(output.Amount, script)
	}

	txObj := tx.NewTransaction(t.Version, txIns, txOuts, 0, t.Testnet)
	secret := util.LittleEndianToBigInt(util.Hash256([]byte(t.Passphrase)))
	pk := ecc.NewPrivateKey(secret)
	if txObj.SignInput(0, pk) {
		serialized := hex.EncodeToString(txObj.Serialize())
		rw.Write([]byte(serialized))
	}
}

func main() {
	http.HandleFunc("/tx/create", createTx)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
