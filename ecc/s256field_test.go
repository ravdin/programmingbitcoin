package ecc

import (
	"math/big"
	"testing"
)

func Test256Field(t *testing.T) {
	t.Run("TestNe", func(t *testing.T) {
		a := NewS256Field(big.NewInt(2), big.NewInt(31))
		b := NewS256Field(big.NewInt(2), big.NewInt(31))
		c := NewS256Field(big.NewInt(15), big.NewInt(31))
		assertEqual(a, b, t)
		assertTrue(fieldIntegerNeOp(), a, c, t)
		assertFalse(fieldIntegerNeOp(), a, b, t)
	})
}
