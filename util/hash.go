package util

import (
	"crypto/sha256"
	"golang.org/x/crypto/ripemd160"
	"hash"
)

// Calculate the hash of hasher over buf.
func calcHash(buf []byte, hasher hash.Hash) []byte {
	hasher.Write(buf)
	return hasher.Sum(nil)
}

// sha256 followed by ripemd160
func Hash160(buf []byte) []byte {
	return calcHash(calcHash(buf, sha256.New()), ripemd160.New())
}

// Two rounds of sha256
func Hash256(buf []byte) []byte {
	return calcHash(calcHash(buf, sha256.New()), sha256.New())
}
