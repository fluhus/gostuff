package heaps

import (
	"testing"

	"golang.org/x/exp/slices"
)

func TestHeap(t *testing.T) {
	input := []string{"bb", "a", "ffff", "ddddd"}
	want := []string{"a", "bb", "ddddd", "ffff"}
	h := Min[string]()
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

func TestHeap_big(t *testing.T) {
	input := []int{
		5, 8, 25, 21, 22, 15, 13, 20, 1, 14,
		24, 12, 7, 18, 27, 3, 30, 28, 23, 29,
		19, 2, 6, 4, 26, 9, 17, 10, 11, 16,
	}
	want := []int{
		1, 2, 3, 4, 5, 6, 7, 8, 9, 10,
		11, 12, 13, 14, 15, 16, 17, 18, 19, 20,
		21, 22, 23, 24, 25, 26, 27, 28, 29, 30,
	}
	h := Min[int]()
	for _, v := range input {
		h.Push(v)
	}
	if ln := h.Len(); ln != len(input) {
		t.Fatalf("Len()=%d, want %d", ln, len(input))
	}
	var got []int
	for h.Len() > 0 {
		got = append(got, h.Pop())
	}
	if !slices.Equal(got, want) {
		t.Fatalf("Pop=%v, want %v", got, want)
	}
}
