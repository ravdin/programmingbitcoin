package code_ch01

import (
	"testing"
)

type FieldElementTest struct {
	a        *FieldElement
	b        *FieldElement
	expected *FieldElement
}

func TestAdd(t *testing.T) {
	tests := []FieldElementTest{
		FieldElementTest{
			a:        NewFieldElement(2, 31),
			b:        NewFieldElement(15, 31),
			expected: NewFieldElement(17, 31),
		},
		FieldElementTest{
			a:        NewFieldElement(17, 31),
			b:        NewFieldElement(21, 31),
			expected: NewFieldElement(7, 31),
		},
	}
	for _, test := range tests {
		if actual := test.a.Add(test.b); !actual.Eq(test.expected) {
			t.Errorf("Wanted %+v, got %+v", test.expected, actual)
		}
	}
}

func TestPow(t *testing.T) {
	a := NewFieldElement(17, 31)
	actual := a.Pow(3)
	expected := NewFieldElement(15, 31)
	if !actual.Eq(expected) {
		t.Errorf("Wanted %+v, got %+v", expected, actual)
	}
}
