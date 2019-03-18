package code_ch02

// Helper class for testing Point with non-finite integer fields.
type IntWrapper int64

func (self IntWrapper) Add(other interface{}) FieldInteger {
	return self + other.(IntWrapper)
}

func (self IntWrapper) Sub(other interface{}) FieldInteger {
	return self - other.(IntWrapper)
}

func (self IntWrapper) Mul(other interface{}) FieldInteger {
	return self * other.(IntWrapper)
}

func (self IntWrapper) Div(other interface{}) FieldInteger {
	return self / other.(IntWrapper)
}

func (self IntWrapper) Pow(exponent int64) FieldInteger {
	var result IntWrapper = 1
	var current IntWrapper = self
	for exponent > 0 {
		if exponent&1 == 1 {
			result *= current
		}
		current *= current
		exponent /= 2
	}
	return result
}

func (self IntWrapper) Rmul(coeff int64) FieldInteger {
	panic("Not implemented")
}

func (self IntWrapper) Eq(other interface{}) bool {
	return self == other.(IntWrapper)
}

func (self IntWrapper) Ne(other interface{}) bool {
	return self != other.(IntWrapper)
}
