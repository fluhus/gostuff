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
