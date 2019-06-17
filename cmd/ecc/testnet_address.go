package main

import (
	"fmt"
	"os"

	"github.com/ravdin/programmingbitcoin/ecc"
	"github.com/ravdin/programmingbitcoin/util"
)

// From Chapter 4, excercise 9.
// Create a testnet address from a passphrase.
func main() {
	passphrase := os.Args[1]
	secret := util.LittleEndianToBigInt(util.Hash256([]byte(passphrase)))
	pk := ecc.NewPrivateKey(secret)
	fmt.Fprintf(os.Stdout, "%s\n", pk.Point.Address(true, true))
}
