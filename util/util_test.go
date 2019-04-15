package util

import (
	"bytes"
	"testing"
)

func TestLittleEndianToInt(t *testing.T) {
	tests := map[uint64]string{
		10011545: "99c3980000000000",
		32454049: "a135ef0100000000",
		254:      "fe00",
	}
	for expected, h := range tests {
		actual := LittleEndianToInt64(HexStringToBytes(h))
		if actual != expected {
			t.Errorf("Expected %d, got %d", expected, actual)
		}
	}
}

func TestIntToLittleEndian(t *testing.T) {
	expected := []byte{1, 0, 0, 0}
	actual := Int32ToLittleEndian(1)
	if !bytes.Equal(expected, actual) {
		t.Errorf("Expected %v, got %v", expected, actual)
	}
}
