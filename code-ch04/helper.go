package code_ch04

import (
	"crypto/sha256"
  "crypto/ripemd160"
)

const BASE58_ALPHABET string = "123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz"

// sha256 followed by ripemd160
func hash160(s string) string {
  hasherSHA256 := sha256.New()
  hasherMD160 := ripemd160.New()
  hasherSHA256.Write([]byte(s))
  hasherMD160.Write(hasherSHA256.Sum(nil))
  return string(hasherMD160.Sum(nil))
}
