package util

import (
	"crypto/sha256"
	"hash"

	"golang.org/x/crypto/ripemd160"
)

// Calculate the hash of hasher over buf.
func calcHash(buf []byte, hasher hash.Hash) []byte {
	hasher.Write(buf)
	return hasher.Sum(nil)
}

// Hash160 is sha256 followed by ripemd160
func Hash160(buf []byte) []byte {
	return calcHash(calcHash(buf, sha256.New()), ripemd160.New())
}

// Hash256 is two rounds of sha256
func Hash256(buf []byte) []byte {
	return calcHash(calcHash(buf, sha256.New()), sha256.New())
}
