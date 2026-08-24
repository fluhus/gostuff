package hashx

import (
	"testing"

	"github.com/fluhus/gostuff/sets"
)

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
