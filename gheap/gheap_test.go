package gheap

import (
	"testing"

	"golang.org/x/exp/slices"
)

func TestHeap(t *testing.T) {
	input := []string{"bb", "a", "ffff", "ddddd"}
	want := []string{"a", "bb", "ddddd", "ffff"}
	h := NewOrdered[string]()
	for _, v := range input {
		h.Push(v)
	}
	if ln := h.Len(); ln != len(input) {
		t.Fatalf("Len()=%d, want %d", ln, len(input))
	}
	var got []string
	for h.Len() > 0 {
		got = append(got, h.Pop())
	}
	if !slices.Equal(got, want) {
		t.Fatalf("Pop=%v, want %v", got, want)
	}
}
