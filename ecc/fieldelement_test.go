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
		tests := [][]int64{
			{2, 15, 17},
			{17, 21, 7},
		}
		for _, test := range tests {
			a := NewFieldElement(test[0], 31)
			b := NewFieldElement(test[1], 31)
			actual := new(FieldElement).Add(a, b)
			expected := NewFieldElement(test[2], 31)
			assertEqual(actual, expected, t)
		}
	})

	t.Run("TestSub", func(t *testing.T) {
		tests := [][]int64{
			{29, 4, 25},
			{15, 30, 16},
		}
		for _, test := range tests {
			a := NewFieldElement(test[0], 31)
			b := NewFieldElement(test[1], 31)
			actual := new(FieldElement).Sub(a, b)
			expected := NewFieldElement(test[2], 31)
			assertEqual(actual, expected, t)
		}
	})

	t.Run("TestMul", func(t *testing.T) {
		tests := [][]int64{
			{24, 19, 22},
		}
		for _, test := range tests {
			a := NewFieldElement(test[0], 31)
			b := NewFieldElement(test[1], 31)
			actual := new(FieldElement).Mul(a, b)
			expected := NewFieldElement(test[2], 31)
			assertEqual(actual, expected, t)
		}
	})

	t.Run("TestPow", func(t *testing.T) {
		tests := [][]int64{
			{1, 17, 3, 15},
			{18, 5, 5, 16},
		}
		for _, test := range tests {
			actual := NewFieldElement(test[1], 31)
			actual.Pow(actual, big.NewInt(test[2])).Mul(actual, NewFieldElement(test[0], 31))
			expected := NewFieldElement(test[3], 31)
			assertEqual(actual, expected, t)
		}
	})

	t.Run("TestDiv", func(t *testing.T) {
		tests := [][]int64{
			{3, 24, 1, 1, 4},
			{17, 1, -3, 1, 29},
			{4, 1, -4, 11, 13},
		}
		for _, test := range tests {
			a := NewFieldElement(test[0], 31)
			b := NewFieldElement(test[1], 31)
			c := big.NewInt(test[2])
			d := NewFieldElement(test[3], 31)
			expected := NewFieldElement(test[4], 31)
			actual := new(FieldElement)
			actual.Div(a, b).Pow(actual, c).Mul(actual, d)
			assertEqual(actual, expected, t)
		}
	})
}
