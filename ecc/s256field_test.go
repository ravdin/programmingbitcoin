package ecc

import (
	"math/big"
	"testing"
)

func Test256Field(t *testing.T) {
	t.Run("TestNe", func(t *testing.T) {
		a := NewS256FieldFromInt64(2, 31)
		b := NewS256FieldFromInt64(2, 31)
		c := NewS256FieldFromInt64(15, 31)
		assertEqual(a, b, t)
		assertTrue(fieldIntegerNeOp(), a, c, t)
		assertFalse(fieldIntegerNeOp(), a, b, t)
	})

	t.Run("TestAdd", func(t *testing.T) {
		a := NewS256FieldFromInt64(2, 31)
		b := NewS256FieldFromInt64(15, 31)
		assertEqual(a.Add(b), NewS256FieldFromInt64(17, 31), t)
		a = NewS256FieldFromInt64(17, 31)
		b = NewS256FieldFromInt64(21, 31)
		assertEqual(a.Add(b), NewS256FieldFromInt64(7, 31), t)
	})

	t.Run("TestSub", func(t *testing.T) {
		a := NewS256FieldFromInt64(29, 31)
		b := NewS256FieldFromInt64(4, 31)
		assertEqual(a.Sub(b), NewS256FieldFromInt64(25, 31), t)
		a = NewS256FieldFromInt64(15, 31)
		b = NewS256FieldFromInt64(30, 31)
		assertEqual(a.Sub(b), NewS256FieldFromInt64(16, 31), t)
	})

	t.Run("TestMul", func(t *testing.T) {
		a := NewS256FieldFromInt64(24, 31)
		b := NewS256FieldFromInt64(19, 31)
		assertEqual(a.Mul(b), NewS256FieldFromInt64(22, 31), t)
	})

	t.Run("TestPow", func(t *testing.T) {
		a := NewS256FieldFromInt64(17, 31)
		assertEqual(a.Pow(big.NewInt(3)), NewS256FieldFromInt64(15, 31), t)
		a = NewS256FieldFromInt64(5, 31)
		b := NewS256FieldFromInt64(18, 31)
		assertEqual(a.Pow(big.NewInt(5)).Mul(b), NewS256FieldFromInt64(16, 31), t)
	})

	t.Run("TestDiv", func(t *testing.T) {
		a := NewS256FieldFromInt64(3, 31)
		b := NewS256FieldFromInt64(24, 31)
		assertEqual(a.Div(b), NewS256FieldFromInt64(4, 31), t)
		a = NewS256FieldFromInt64(17, 31)
		assertEqual(a.Pow(big.NewInt(-3)), NewS256FieldFromInt64(29, 31), t)
		a = NewS256FieldFromInt64(4, 31)
		b = NewS256FieldFromInt64(11, 31)
		assertEqual(a.Pow(big.NewInt(-4)).Mul(b), NewS256FieldFromInt64(13, 31), t)
	})
}
