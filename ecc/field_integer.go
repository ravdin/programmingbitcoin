package ecc

import (
	"math/big"
)

type FieldInteger interface {
	Add(other FieldInteger) FieldInteger
	Sub(other FieldInteger) FieldInteger
	Mul(other FieldInteger) FieldInteger
	Div(other FieldInteger) FieldInteger
	Pow(exponent *big.Int) FieldInteger
	Rmul(coeff *big.Int) FieldInteger
	Eq(other FieldInteger) bool
	Ne(other FieldInteger) bool
}

type FieldIntegerConverter func(interface{}) FieldInteger
