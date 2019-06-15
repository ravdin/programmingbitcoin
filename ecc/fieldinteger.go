package ecc

import (
	"math/big"
)

// FieldInteger is a representation of an integer for finite field math.
type FieldInteger interface {
	Add(x, y FieldInteger) FieldInteger
	Sub(x, y FieldInteger) FieldInteger
	Mul(x, y FieldInteger) FieldInteger
	Div(x, y FieldInteger) FieldInteger
	Pow(n FieldInteger, exponent *big.Int) FieldInteger
	Cmul(n FieldInteger, coefficient *big.Int) FieldInteger
	Eq(other FieldInteger) bool
	Ne(other FieldInteger) bool
	Copy() FieldInteger
	Set(n FieldInteger) FieldInteger
}
