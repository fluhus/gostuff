//go:build go1.27

package hashx

import (
	"testing"

	"github.com/fluhus/gostuff/sets"
)

func TestInt_method(t *testing.T) {
	const n = 100000

	hashes := sets.Set[uint64]{}
	hx := New()
	for i := range n {
		hashes.Add(hx.Int(i))
	}
	if len(hashes) != n {
		t.Errorf("len(hashes)=%v, want %v", len(hashes), n)
	}

	hashes = sets.Set[uint64]{}
	hx = New()
	for i := range n {
		hashes.Add(hx.Int(i))
		hashes.Add(hx.Int(-i))
	}
	want := n*2 - 1
	if len(hashes) != want {
		t.Errorf("len(hashes)=%v, want %v", len(hashes), want)
	}
}

func TestInt_equal(t *testing.T) {
	const x int = 1234567890
	h := New()
	a, b, c, d := h.intShifts(x), h.intUnsafe(x), h.intBinary(x), IntHashx(h, x)
	if a != b || a != c || a != d {
		t.Fatalf("mismatching hashes: %v %v %v %v", a, b, c, d)
	}
}

// Benchmarking several ways to hash an int.

func BenchmarkInt(b *testing.B) {
	const x int = 1234567890
	b.Run("shifts", func(b *testing.B) {
		h := New()
		for b.Loop() {
			h.intShifts(x)
		}
	})
	b.Run("unsafe", func(b *testing.B) {
		h := New()
		for b.Loop() {
			h.intUnsafe(x)
		}
	})
	b.Run("binary", func(b *testing.B) {
		h := New()
		for b.Loop() {
			h.intBinary(x)
		}
	})
}
