package common

import (
	"bytes"
	"testing"
)

func TestAddSchemaRegistryHeader(t *testing.T) {
	testSchemaIds := [5]uint32{0, 1, 10, 500, 1000000}
	expectedBytes := [5][5]byte{
		{0, 0, 0, 0, 0},
		{0, 0, 0, 0, 1},
		{0, 0, 0, 0, 10},
		{0, 0, 0, 1, 244},
		{0, 0, 15, 66, 64},
	}
	cnt := 0
	for _, testSchemaID := range testSchemaIds {
		var buf bytes.Buffer
		AddSchemaRegistryHeader(&buf, testSchemaID)
		actual := buf.Bytes()
		expected := expectedBytes[cnt]
		if len(actual) != len(expected) {
			t.Errorf("Length of arrays ard different for case %d", cnt)
		}
		for i := 0; i < 5; i++ {
			if actual[i] != expected[i] {
				t.Errorf("Different values for case %d pos %d. Expected: %d, got %d", cnt, i, expected[i], actual[i])
			}
		}
		cnt++
	}
}
