package ecc

import (
	"math/big"
	"testing"
)

func TestFieldElement(t *testing.T) {
	t.Run("TestNe", func(t *testing.T) {
		a := newFieldElement(2, 31)
		b := newFieldElement(2, 31)
		c := newFieldElement(15, 31)
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
			a := newFieldElement(test[0], 31)
			b := newFieldElement(test[1], 31)
			actual := new(fieldElement).Add(a, b)
			expected := newFieldElement(test[2], 31)
			assertEqual(actual, expected, t)
		}
	})

	t.Run("TestSub", func(t *testing.T) {
		tests := [][]int64{
			{29, 4, 25},
			{15, 30, 16},
		}
		for _, test := range tests {
			a := newFieldElement(test[0], 31)
			b := newFieldElement(test[1], 31)
			actual := new(fieldElement).Sub(a, b)
			expected := newFieldElement(test[2], 31)
			assertEqual(actual, expected, t)
		}
	})

	t.Run("TestMul", func(t *testing.T) {
		tests := [][]int64{
			{24, 19, 22},
		}
		for _, test := range tests {
			a := newFieldElement(test[0], 31)
			b := newFieldElement(test[1], 31)
			actual := new(fieldElement).Mul(a, b)
			expected := newFieldElement(test[2], 31)
			assertEqual(actual, expected, t)
		}
	})

	t.Run("TestPow", func(t *testing.T) {
		tests := [][]int64{
			{1, 17, 3, 15},
			{18, 5, 5, 16},
		}
		for _, test := range tests {
			actual := newFieldElement(test[1], 31)
			actual.Pow(actual, big.NewInt(test[2])).Mul(actual, newFieldElement(test[0], 31))
			expected := newFieldElement(test[3], 31)
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
			a := newFieldElement(test[0], 31)
			b := newFieldElement(test[1], 31)
			c := big.NewInt(test[2])
			d := newFieldElement(test[3], 31)
			expected := newFieldElement(test[4], 31)
			actual := new(fieldElement)
			actual.Div(a, b).Pow(actual, c).Mul(actual, d)
			assertEqual(actual, expected, t)
		}
	})
}
