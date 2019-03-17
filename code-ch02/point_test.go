package ecc_test

import (
	"programmingbitcoin/code-ch02"
	"testing"
)

func TestAdd0(t *testing.T) {
    var x ecc.IntWrapper = 2
    var y ecc.IntWrapper = 5
    var a ecc.IntWrapper = 5
    var b ecc.IntWrapper = 7
    pa, _ := ecc.NewPoint(nil, nil, a, b)
    pb, err := ecc.NewPoint(x, y, a, b)
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
