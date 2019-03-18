package code_ch01

import (
	"testing"
)

type FieldElementOp func(*FieldElement, *FieldElement) bool

func fieldElementNeOp() FieldElementOp {
	return func(a *FieldElement, b *FieldElement) bool {
		return a.Ne(b)
	}
}

func assertEqual(a *FieldElement, b *FieldElement, t *testing.T) {
	if !a.Eq(b) {
		t.Errorf("%v is not equal to %v\n", a, b)
	}
}

func assertTrue(op FieldElementOp, a *FieldElement, b *FieldElement, t *testing.T) {
	if !op(a, b) {
		t.Errorf("Assertion failed!\n")
	}
}

func assertFalse(op FieldElementOp, a *FieldElement, b *FieldElement, t *testing.T) {
	if op(a, b) {
		t.Errorf("Assertion failed!\n")
	}
}
