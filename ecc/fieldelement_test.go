package ecc

import (
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
}
