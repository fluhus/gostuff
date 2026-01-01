package reservoir

import (
	"fmt"
	"math"
	"math/rand/v2"
)

const debugPrints = false

// NewFunc returns a new sampler that samples f(n) elements.
// The sampling process chooses each element with a probability of
// f(n) / n, so the number of elements may vary but on average is f(n).
//
// f should return the expected number of elements for a given a total of
// n elements.
// f should be non-decreasing, meaning for n1<n2, f(n1)<=f(n2).
//
// The memory footprint is O(f(n)) for functions that are <= O(sqrt(n)).
func NewFunc[T any](f func(n int) float64) *SamplerFunc[T] {
	return &SamplerFunc[T]{
		f: f,
		r: rand.New(rand.NewPCG(rand.Uint64(), rand.Uint64())),
	}
}

// Add maybe adds t to the selected sample.
func (s *SamplerFunc[T]) Add(t T) {
	s.n++
	n := s.f(s.n)
	if math.IsInf(n, 0) || math.IsNaN(n) {
		panic(fmt.Sprintf("expected elements is not a finite number: %v", n))
	}
	if n < s.last {
		panic(fmt.Sprintf("expected elements decreased: from %v to %v",
			s.last, n))
	}
	if n > float64(s.n) {
		panic(fmt.Sprintf("expected elements is higher than n: %v > %v",
			n, s.n))
	}
	if n == 0 {
		return
	}
	s.last = n
	p := n / float64(s.n)
	if s.r.Float64() < p {
		s.e = append(s.e, tuple[T]{t, p})
	}
}

// Elements returns a new slice containing on average f(n) elements.
// Each call may return a different slice because of a second round of
// sampling that takes place.
func (s *SamplerFunc[T]) Elements() []T {
	if debugPrints {
		fmt.Println(len(s.e))
	}
	var e []T
	for _, t := range s.e {
		p := s.f(s.n) / float64(s.n) / t.p
		if rand.Float64() < p {
			e = append(e, t.v)
		}
	}
	return e
}

// SamplerFunc is a sampler that uses a function to determine
// how many elements is should sample.
type SamplerFunc[T any] struct {
	f    func(int) float64
	e    []tuple[T]
	r    *rand.Rand
	n    int
	last float64
}

// A sampled value along with the probability with which it was sampled.
type tuple[T any] struct {
	v T
	p float64
}
