package code_ch02

import (
	"testing"
)

func TestAdd0(t *testing.T) {
	var x IntWrapper = 2
	var y IntWrapper = 5
	var a IntWrapper = 5
	var b IntWrapper = 7
	pa, _ := NewPoint(nil, nil, a, b)
	pb, err := NewPoint(x, y, a, b)
	if err != nil {
		t.Errorf("%v\n", err)
	}
	//pc, err := ecc.NewPoint(x, y, a, b)
	actual := pa.Add(pb)
	if !actual.Eq(pb) {
		t.Errorf("Wanted %+v, got %+v\n", pb, actual)
	}
	//self.assertEqual(a + b, b)
	//self.assertEqual(b + a, b)
	//self.assertEqual(b + c, a)
}
