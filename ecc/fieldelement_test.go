package ecc

import (
	"testing"
)

func TestNe(t *testing.T) {
	a := NewFieldElement(2, 31)
	b := NewFieldElement(2, 31)
	c := NewFieldElement(15, 31)
	assertEqual(a, b, t)
	assertTrue(fieldElementNeOp(), a, c, t)
	assertFalse(fieldElementNeOp(), a, b, t)
}
