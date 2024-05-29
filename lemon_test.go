package lemon

import (
	"testing"
)

func TestAdd(t *testing.T) {
	result := Print()
	expected := 5

	if result != expected {
		t.Errorf("Print() = %d; want %d", result, expected)
	}
}
