package util

import "testing"

func TestMurmur3(t *testing.T) {
	tests := map[string][][]uint32{
		"hello world": [][]uint32{
			{42, 3926694905},
			{4221880255, 3771040970},
		},
		"goodbye": [][]uint32{
			{42, 1987570198},
			{4221880255, 757891893},
		},
	}
	for phrase, vals := range tests {
		for _, val := range vals {
			seed := val[0]
			expected := val[1]
			actual := Murmur3([]byte(phrase), seed)
			if expected != actual {
				t.Errorf("Expected %d, got %d", expected, actual)
			}
		}
	}
}
