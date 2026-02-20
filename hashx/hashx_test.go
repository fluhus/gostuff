package hashx

import (
	"testing"

	"github.com/fluhus/gostuff/sets"
)

func TestInt(t *testing.T) {
	const n = 100000

	hashes := sets.Set[uint64]{}
	hx := New()
	for i := range n {
		hashes.Add(IntHashx(hx, i))
	}
	if len(hashes) != n {
		t.Errorf("len(hashes)=%v, want %v", len(hashes), n)
	}

	hashes = sets.Set[uint64]{}
	hx = New()
	for i := range n {
		hashes.Add(IntHashx(hx, i))
		hashes.Add(IntHashx(hx, -i))
	}
	want := n*2 - 1
	if len(hashes) != want {
		t.Errorf("len(hashes)=%v, want %v", len(hashes), want)
	}
}

func TestBytesString(t *testing.T) {
	s := "abcdefghijklmnopqrstuvwxyz"
	b := []byte(s)
	hashes := sets.Set[uint64]{}
	hx := New()
	n := 0
	for i := range s {
		for j := range i {
			n++
			ss := s[j : i+1]
			h1 := hx.String(ss)
			h2 := hx.Bytes(b[j : i+1])
			if h1 != h2 {
				t.Errorf("String(%s) != Bytes(%s): %v, %v", ss, ss, h1, h2)
			}
			hashes.Add(h1)
		}
	}
	if len(hashes) != n {
		t.Errorf("len(hashes)=%v, want %v", len(hashes), n)
	}
}
