package heaps

import (
	"fmt"
	"math/rand/v2"
	"slices"
	"testing"

	"github.com/fluhus/gostuff/snm"
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
	h := New(func(i1, i2 int) bool {
		return i1 < i2
	})
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

func TestHeap_pushSlice(t *testing.T) {
	input := []int{
		5, 8, 25, 21, 22, 15, 13, 20, 1, 14,
		24, 12, 7, 18, 27, 3, 30, 28, 23, 29,
		19, 2, 6, 4, 26, 9, 17, 10, 11, 16,
	}
	h := Min[int]()
	h.PushSlice(input)
	for i := range h.a {
		if i == 0 {
			continue
		}
		ia := (i - 1) / 2
		if h.a[i] < h.a[ia] {
			t.Errorf("h[%d] < h[%d]: %d < %d", i, ia, h.a[i], h.a[ia])
		}
	}
}

func Benchmark(b *testing.B) {
	for _, n := range []int{1000, 10000, 100000, 1000000} {
		nums := snm.Slice(n, func(i int) int {
			return rand.Int()
		})
		b.Run(fmt.Sprint("Heap.Push.", n), func(b *testing.B) {
			for range b.N {
				h := Min[int]()
				for _, i := range nums {
					h.Push(i)
				}
			}
		})
		b.Run(fmt.Sprint("Heap.PushSlice.", n), func(b *testing.B) {
			for range b.N {
				h := Min[int]()
				h.PushSlice(nums)
			}
		})
	}
}

func FuzzHeap(f *testing.F) {
	f.Add(0, 0, 0, 0, 0, 0, 0)
	f.Fuzz(func(t *testing.T, a, b, c, d, e, f, g int) {
		h := Min[int]()
		h.Push(a)
		h.Push(b)
		h.Push(c)
		h.Push(d)
		h.Push(e)
		h.Push(f)
		h.Push(g)
		got := make([]int, 0, 7)
		for h.Len() > 0 {
			got = append(got, h.Pop())
		}
		if !slices.IsSorted(got) {
			t.Fatalf("Min().Pop()=%v, want sorted", got)
		}
	})
}
