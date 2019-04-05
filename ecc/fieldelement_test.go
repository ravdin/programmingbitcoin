package ecc

import (
	"math/big"
	"testing"
)

func TestFieldElement(t *testing.T) {
	t.Run("TestNe", func(t *testing.T) {
		a := NewFieldElement(2, 31)
		b := NewFieldElement(2, 31)
		c := NewFieldElement(15, 31)
		assertEqual(a, b, t)
		assertTrue(fieldIntegerNeOp(), a, c, t)
		assertFalse(fieldIntegerNeOp(), a, b, t)
	})

	t.Run("TestAdd", func(t *testing.T) {
		a := NewFieldElement(2, 31)
		b := NewFieldElement(15, 31)
		assertEqual(a.Add(b), NewFieldElement(17, 31), t)
		a = NewFieldElement(17, 31)
		b = NewFieldElement(21, 31)
		assertEqual(a.Add(b), NewFieldElement(7, 31), t)
	})

	t.Run("TestSub", func(t *testing.T) {
		a := NewFieldElement(29, 31)
		b := NewFieldElement(4, 31)
		assertEqual(a.Sub(b), NewFieldElement(25, 31), t)
		a = NewFieldElement(15, 31)
		b = NewFieldElement(30, 31)
		assertEqual(a.Sub(b), NewFieldElement(16, 31), t)
	})

	t.Run("TestMul", func(t *testing.T) {
		a := NewFieldElement(24, 31)
		b := NewFieldElement(19, 31)
		assertEqual(a.Mul(b), NewFieldElement(22, 31), t)
	})

	t.Run("TestPow", func(t *testing.T) {
		a := NewFieldElement(17, 31)
		assertEqual(a.Pow(big.NewInt(3)), NewFieldElement(15, 31), t)
		a = NewFieldElement(5, 31)
		b := NewFieldElement(18, 31)
		assertEqual(a.Pow(big.NewInt(5)).Mul(b), NewFieldElement(16, 31), t)
	})

	t.Run("TestDiv", func(t *testing.T) {
		a := NewFieldElement(3, 31)
		b := NewFieldElement(24, 31)
		assertEqual(a.Div(b), NewFieldElement(4, 31), t)
		a = NewFieldElement(17, 31)
		assertEqual(a.Pow(big.NewInt(-3)), NewFieldElement(29, 31), t)
		a = NewFieldElement(4, 31)
		b = NewFieldElement(11, 31)
		assertEqual(a.Pow(big.NewInt(-4)).Mul(b), NewFieldElement(13, 31), t)
	})
}
