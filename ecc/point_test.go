package ecc

import (
	"testing"
)

func withIntWrapper() FieldIntegerConverter {
	return func(n interface{}) FieldInteger {
		return newIntWrapper(n.(int))
	}
}

func TestPoint(t *testing.T) {
	option := withIntWrapper()

	t.Run("TestNe", func(t *testing.T) {
		a, _ := NewPoint(3, -7, 5, 7, option)
		b, _ := NewPoint(18, 77, 5, 7, option)
		if a.Eq(b) {
			t.Errorf("Expected a != b")
		}
		if !a.Eq(a) {
			t.Errorf("Expected a == a")
		}
	})

	t.Run("TestOnCurve", func(t *testing.T) {
		_, err := NewPoint(2, 4, 5, 7, option)
		if err == nil {
			t.Errorf("Point is not on curve, expected error!")
		}
		// These should not raise an error
		p1, err := NewPoint(3, -7, 5, 7, option)
		if p1 == nil || err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		p2, err := NewPoint(18, 77, 5, 7, option)
		if p2 == nil || err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
	})
	/*

	   def test_add0(self):
	       a = Point(x=None, y=None, a=5, b=7)
	       b = Point(x=2, y=5, a=5, b=7)
	       c = Point(x=2, y=-5, a=5, b=7)
	       self.assertEqual(a + b, b)
	       self.assertEqual(b + a, b)
	       self.assertEqual(b + c, a)

	   def test_add1(self):
	       a = Point(x=3, y=7, a=5, b=7)
	       b = Point(x=-1, y=-1, a=5, b=7)
	       self.assertEqual(a + b, Point(x=2, y=-5, a=5, b=7))

	   def test_add2(self):
	       a = Point(x=-1, y=1, a=5, b=7)
	       self.assertEqual(a + a, Point(x=18, y=-77, a=5, b=7))
	*/
}
