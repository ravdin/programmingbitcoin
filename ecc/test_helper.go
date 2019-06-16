package ecc

import (
	"testing"
)

type fieldIntegerOp func(FieldInteger, FieldInteger) bool

func fieldIntegerNeOp() fieldIntegerOp {
	return func(a FieldInteger, b FieldInteger) bool {
		return a.Ne(b)
	}
}

func assertEqual(a FieldInteger, b FieldInteger, t *testing.T) {
	if !a.Eq(b) {
		t.Errorf("%v is not equal to %v\n", a, b)
	}
}

func assertTrue(op fieldIntegerOp, a FieldInteger, b FieldInteger, t *testing.T) {
	if !op(a, b) {
		t.Errorf("Assertion failed!\n")
	}
}

func assertFalse(op fieldIntegerOp, a FieldInteger, b FieldInteger, t *testing.T) {
	if op(a, b) {
		t.Errorf("Assertion failed!\n")
	}
}
