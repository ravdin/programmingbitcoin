package ecc_test

import (
	"programmingbitcoin/code-ch01"
	"testing"
)

type FieldElementTest struct {
	a        *ecc.FieldElement
	b        *ecc.FieldElement
	expected *ecc.FieldElement
}

func TestAdd(t *testing.T) {
	tests := []FieldElementTest{
		FieldElementTest{
			a:        ecc.NewFieldElement(2, 31),
			b:        ecc.NewFieldElement(15, 31),
			expected: ecc.NewFieldElement(17, 31),
		},
		FieldElementTest{
			a:        ecc.NewFieldElement(17, 31),
			b:        ecc.NewFieldElement(21, 31),
			expected: ecc.NewFieldElement(7, 31),
		},
	}
	for _, test := range tests {
		if actual := test.a.Add(test.b); !actual.Eq(test.expected) {
			t.Errorf("Wanted %+v, got %+v", test.expected, actual)
		}
	}
}

func TestPow(t *testing.T) {
	a := ecc.NewFieldElement(17, 31)
	actual := a.Pow(3)
	expected := ecc.NewFieldElement(15, 31)
	if !actual.Eq(expected) {
		t.Errorf("Wanted %+v, got %+v", expected, actual)
	}
}
