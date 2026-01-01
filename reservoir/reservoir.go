// Package reservoir implements reservoir sampling.
//
// Reservoir sampling allows sampling m uniformly random elements
// from a stream,
// using O(m) memory regardless of the stream length.
package reservoir

import "math/rand/v2"

// Sampler samples a fixed number of elements with uniform distribution
// from a stream.
type Sampler[T any] struct {
	Elements []T // Elements selected so far.
	r        *rand.Rand
	n        int
}

// New returns a new sampler that samples n elements.
func New[T any](n int) *Sampler[T] {
	return &Sampler[T]{
		Elements: make([]T, 0, n),
		r:        rand.New(rand.NewPCG(rand.Uint64(), rand.Uint64())),
	}
}

// Add maybe adds t to the selected sample.
func (s *Sampler[T]) Add(t T) {
	s.n++
	if len(s.Elements) < cap(s.Elements) {
		s.Elements = append(s.Elements, t)
		return
	}
	i := s.r.IntN(s.n)
	if i >= len(s.Elements) {
		return
	}
	s.Elements[i] = t
}
