// A generic test suite for rolling hashes.

package rhash

import (
	"crypto/rand"
	"hash"
	"testing"

	"github.com/fluhus/gostuff/gnum"
	"github.com/fluhus/gostuff/sets"
)

// Runs the test suite for a hash64.
func test64(t *testing.T, f func(n int) hash.Hash64) {
	t.Run("basic", func(t *testing.T) { test64basic(t, f) })
	t.Run("cyclic", func(t *testing.T) { test64cyclic(t, f) })
	t.Run("big-n", func(t *testing.T) { test64bigN(t, f) })
}

func test64basic(t *testing.T, f func(n int) hash.Hash64) {
	data := []byte("amitamit")
	tests := []struct {
		n        int
		wantSize []int
		wantEq   []int
	}{
		{2, []int{1, 2, 3, 4, 5, 5, 5, 5}, []int{-1, -1, -1, -1, -1, 1, 2, 3}},
		{3, []int{1, 2, 3, 4, 5, 6, 6, 6}, []int{-1, -1, -1, -1, -1, -1, 2, 3}},
	}
	for _, test := range tests {
		slice := []uint64{}
		set := sets.Set[uint64]{}
		h := f(test.n)
		for i := range data {
			h.Write(data[i : i+1])
			slice = append(slice, h.Sum64())
			set.Add(h.Sum64())
			if len(set) != test.wantSize[i] {
				t.Fatalf("n=%d #%d: set size=%v, want %v",
					test.n, i, len(set), test.wantSize[i])
			}
			if test.wantEq[i] != -1 && h.Sum64() != slice[test.wantEq[i]] {
				t.Fatalf("n=%d #%d: Sum64()=%d, want %d",
					test.n, i, h.Sum64(), slice[test.wantEq[i]])
			}
		}
	}
}

func test64cyclic(t *testing.T, f func(n int) hash.Hash64) {
	inputs := []string{
		"asdjadasdk",
		"uioewrmnoc",
		"wiewuwikxa",
		"mfhddl/lcc",
		"28n9789dkd",
	}
	h := f(10)
	for _, input := range inputs {
		h.Write([]byte(input))
		h2 := f(10)
		h2.Write([]byte(input))
		got, want := h.Sum64(), h2.Sum64()
		if got != want {
			t.Fatalf("Sum64(%q)=%v, want %v", input, got, want)
		}
	}
}

func test64bigN(t *testing.T, f func(n int) hash.Hash64) {
	const n = 100

	// Create random input.
	buf := make([]byte, n)
	_, err := rand.Read(buf)
	if err != nil {
		t.Fatalf("rand.Read() failed: %v", err)
	}

	// Repeat 3 times.
	buf = append(buf, buf...)
	buf = append(buf, buf...)

	h := f(n)
	hashes := sets.Set[uint64]{}
	for i := range buf {
		h.Write(buf[i : i+1])
		hashes.Add(h.Sum64())
		want := gnum.Min2(i+1, n*2-1)
		if len(hashes) != want {
			t.Fatalf("got %d unique hashes, want %d", len(hashes), want)
		}
	}
}

// Runs the test suite for a hash32.
func test32(t *testing.T, f func(n int) hash.Hash32) {
	t.Run("basic", func(t *testing.T) { test32basic(t, f) })
	t.Run("cyclic", func(t *testing.T) { test32cyclic(t, f) })
	t.Run("big-n", func(t *testing.T) { test32bigN(t, f) })
}

func test32basic(t *testing.T, f func(n int) hash.Hash32) {
	data := []byte("amitamit")
	tests := []struct {
		n        int
		wantSize []int
		wantEq   []int
	}{
		{2, []int{1, 2, 3, 4, 5, 5, 5, 5}, []int{-1, -1, -1, -1, -1, 1, 2, 3}},
		{3, []int{1, 2, 3, 4, 5, 6, 6, 6}, []int{-1, -1, -1, -1, -1, -1, 2, 3}},
	}
	for _, test := range tests {
		slice := []uint32{}
		set := sets.Set[uint32]{}
		h := f(test.n)
		for i := range data {
			h.Write(data[i : i+1])
			slice = append(slice, h.Sum32())
			set.Add(h.Sum32())
			if len(set) != test.wantSize[i] {
				t.Fatalf("n=%d #%d: set size=%v, want %v",
					test.n, i, len(set), test.wantSize[i])
			}
			if test.wantEq[i] != -1 && h.Sum32() != slice[test.wantEq[i]] {
				t.Fatalf("n=%d #%d: Sum32()=%d, want %d",
					test.n, i, h.Sum32(), slice[test.wantEq[i]])
			}
		}
	}
}

func test32cyclic(t *testing.T, f func(n int) hash.Hash32) {
	inputs := []string{
		"asdjadasdk",
		"uioewrmnoc",
		"wiewuwikxa",
		"mfhddl/lcc",
		"28n9789dkd",
	}
	h := f(10)
	for _, input := range inputs {
		h.Write([]byte(input))
		h2 := f(10)
		h2.Write([]byte(input))
		got, want := h.Sum32(), h2.Sum32()
		if got != want {
			t.Fatalf("Sum32(%q)=%v, want %v", input, got, want)
		}
	}
}

func test32bigN(t *testing.T, f func(n int) hash.Hash32) {
	const n = 100

	// Create random input.
	buf := make([]byte, n)
	_, err := rand.Read(buf)
	if err != nil {
		t.Fatalf("rand.Read() failed: %v", err)
	}

	// Repeat 3 times.
	buf = append(buf, buf...)
	buf = append(buf, buf...)

	h := f(n)
	hashes := sets.Set[uint32]{}
	for i := range buf {
		h.Write(buf[i : i+1])
		hashes.Add(h.Sum32())
		want := gnum.Min2(i+1, n*2-1)
		if len(hashes) != want {
			t.Fatalf("got %d unique hashes, want %d", len(hashes), want)
		}
	}
}
