package bloomfilter

import (
	"fmt"
	"testing"
)

func TestBloomFilter(t *testing.T) {
	t.Run("Test add", func(t *testing.T) {
		bf := New(10, 5, 99)
		tests := [][]string{
			{"Hello World", "0000000a080000000140"},
			{"Goodbye!", "4000600a080000010940"},
		}
		for _, test := range tests {
			item := []byte(test[0])
			expected := test[1]
			bf.Add(item)
			actual := fmt.Sprintf("%x", bf.FilterBytes())
			if expected != actual {
				t.Errorf("Expected %s, got %s", expected, actual)
			}
		}
	})

	t.Run("Test FilterLoad", func(t *testing.T) {
		bf := New(10, 5, 99)
		bf.Add([]byte("Hello World"))
		bf.Add([]byte("Goodbye!"))
		msg := bf.FilterLoad(1)
		expected := "0a4000600a080000010940050000006300000001"
		actual := fmt.Sprintf("%x", msg.Serialize())
		if expected != actual {
			t.Errorf("Expected %s, got %s", expected, actual)
		}
	})
}
