package reservoir

import (
	"math"
	"testing"
)

func TestSamplerFunc_sqrt(t *testing.T) {
	const n = 100
	cnt := make([]int, n)
	for range 10000 {
		s := NewFunc[int](func(n int) float64 { return math.Sqrt(float64(n)) })
		for i := range n {
			s.Add(i)
		}
		for _, i := range s.Elements() {
			cnt[i]++
		}
	}
	want1, want2 := 850, 1150
	for i, c := range cnt {
		if c < want1 || c > want2 {
			t.Errorf("count(%v)=%v, want %v-%v", i, c, want1, want2)
		}
	}
}

func TestSamplerFunc_ratio(t *testing.T) {
	const n = 100
	cnt := make([]int, n)
	for range 10000 {
		s := NewFunc[int](func(n int) float64 { return float64(n) / 10 })
		for i := range n {
			s.Add(i)
		}
		for _, i := range s.Elements() {
			cnt[i]++
		}
	}
	want1, want2 := 850, 1150
	for i, c := range cnt {
		if c < want1 || c > want2 {
			t.Errorf("count(%v)=%v, want %v-%v", i, c, want1, want2)
		}
	}
}
