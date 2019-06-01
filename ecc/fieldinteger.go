package ecc

import (
	"math/big"
)

// Representation of an integer for finite field math.
type FieldInteger interface {
	// Set z to x+y and return z
	Add(x, y FieldInteger) FieldInteger
	// Set z to x-y and return z
	Sub(x, y FieldInteger) FieldInteger
	// Set z to x*y and return z
	Mul(x, y FieldInteger) FieldInteger
	// Set z to x/y and return z
	Div(x, y FieldInteger) FieldInteger
	// Set z to n**exponent and return z
	Pow(n FieldInteger, exponent *big.Int) FieldInteger
	// Set z to n*c where c is an integer coefficient and reutrn z
	Cmul(n FieldInteger, coefficient *big.Int) FieldInteger
	// Return true if z is equal to another FieldInteger
	Eq(other FieldInteger) bool
	// Return true if z is not equal to another FieldInteger
	Ne(other FieldInteger) bool
	// Copy the value to another FieldInteger
	Copy() FieldInteger
	// Set z to the value of another FieldInteger
	Set(n FieldInteger) FieldInteger
}
