package ecc

import (
	"github.com/ravdin/programmingbitcoin/util"
	"math/big"
	"math/rand"
	"testing"
)

func TestPrivateKey(t *testing.T) {
	t.Run("Test Sign", func(t *testing.T) {
		r := rand.New(rand.NewSource(42))
		pk := NewPrivateKey(new(big.Int).Rand(r, N))
		z := new(big.Int).Rand(r, util.HexStringToBigInt("ffffffffffffffffffffffffffffffffffffffffffffffff"))
		sig := pk.Sign(z)
		if !pk.Point.Verify(z, sig) {
			t.Errorf("Private key signature failed!")
		}
	})

	t.Run("Test WIF", func(t *testing.T) {
		pk := NewPrivateKey(util.HexStringToBigInt("ffffffffffffff80000000000000000000000000000000000000000000000000"))
		expected := "L5oLkpV3aqBJ4BgssVAsax1iRa77G5CVYnv9adQ6Z87te7TyUdSC"
		actual := pk.Wif(true, false)
		if actual != expected {
			t.Errorf("Expected %v, got %v", expected, actual)
		}
		pk = NewPrivateKey(util.HexStringToBigInt("fffffffffffffe00000000000000000000000000000000000000000000000000"))
		expected = "93XfLeifX7Jx7n7ELGMAf1SUR6f9kgQs8Xke8WStMwUtrDucMzn"
		actual = pk.Wif(false, true)
		if actual != expected {
			t.Errorf("Expected %v, got %v", expected, actual)
		}
		pk = NewPrivateKey(util.HexStringToBigInt("0dba685b4511dbd3d368e5c4358a1277de9486447af7b3604a69b8d9d8b7889d"))
		expected = "5HvLFPDVgFZRK9cd4C5jcWki5Skz6fmKqi1GQJf5ZoMofid2Dty"
		actual = pk.Wif(false, false)
		if actual != expected {
			t.Errorf("Expected %v, got %v", expected, actual)
		}
		pk = NewPrivateKey(util.HexStringToBigInt("1cca23de92fd1862fb5b76e5f4f50eb082165e5191e116c18ed1a6b24be6a53f"))
		expected = "cNYfWuhDpbNM1JWc3c6JTrtrFVxU4AGhUKgw5f93NP2QaBqmxKkg"
		actual = pk.Wif(true, true)
		if actual != expected {
			t.Errorf("Expected %v, got %v", expected, actual)
		}
	})
}
