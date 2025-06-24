package reservoir

import (
	"fmt"
	"slices"
	"testing"
)

func TestSmall(t *testing.T) {
	var want []int
	r := New[int](10)
	if len(r.Elements) != 0 {
		t.Fatalf("len(E)=%v, want 0", len(r.Elements))
	}
	for i := range 10 {
		r.Add(i)
		want = append(want, i)
		if !slices.Equal(r.Elements, want) {
			t.Fatalf("E=%v, want %v", r.Elements, want)
		}
	}
	r.Add(10)
	want = append(want, 10)
	if slices.Equal(r.Elements, want) {
		t.Fatalf("E=%v, want smaller", r.Elements)
	}
	if len(r.Elements) != 10 {
		t.Fatalf("len(E)=%v, want 10", len(r.Elements))
	}
}

func TestBig(t *testing.T) {
	counts := make([]int, 10)
	for range 1000 {
		r := New[int](5)
		for i := range 10 {
			r.Add(i)
		}
		for _, i := range r.Elements {
			counts[i]++
		}
	}
	wantMin, wantMax := 450, 550
	for i, c := range counts {
		if c < wantMin || c > wantMax {
			t.Errorf("count[%v]=%v, want %v-%v", i, c, wantMin, wantMax)
		}
	}
}

func Example() {
	// Select 10 random elements uniformly out of a stream.
	sampler := New[int](10)
	for i := range 1000 {
		sampler.Add(i)
	}
	fmt.Println("Selected subsample:", sampler.Elements)
}
