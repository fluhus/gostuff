// Package hll provides an implementation of the HyperLogLog algorithm.
//
// A HyperLogLog counter can approximate the cardinality of a set with high
// accuracy and little memory.
//
// Performance
//
// An HLL counter uses 65kb (2^16) memory.
//
// Using murmur3 hash it was able to estimate the cardinality of the set of
// numbers 1 to 10^9 with 0.1% error in 49 seconds on a single core (49 ns/op).
//
// Citation
//
// Flajolet, Philippe; Fusy, Éric; Gandouet, Olivier; Meunier, Frédéric (2007).
// "Hyperloglog: The analysis of a near-optimal cardinality estimation
// algorithm". Discrete Mathematics and Theoretical Computer Science
// Proceedings.
package hll

import (
	"math"
)

const (
	nbits = 16
	m     = 1 << nbits
	mask  = m - 1
	alpha = 0.7213 / (1.0 + 1.079/m)
)

// An HLL is a HyperLogLog counter for arbitrary values.
type HLL struct {
	counters []byte
	hash     func(v interface{}) uint64
}

// New creates a new HyperLogLog counter that uses the given hash function.
func New(hash func(v interface{}) uint64) *HLL {
	return &HLL{
		make([]byte, m),
		hash,
	}
}

// Add adds an element to the counter. Calls hash once.
func (h *HLL) Add(v interface{}) {
	hash := h.hash(v)
	idx := hash & mask
	fp := hash >> nbits
	z := byte(nzeros(fp)) + 1
	if z > h.counters[idx] {
		h.counters[idx] = z
	}
}

// ApproxCount returns the current approximate count. Does not alter the state
// of the counter.
//
// The approximation starts being accurate around one million. Values lower than
// that should generally be regarded as highly inaccurate.
func (h *HLL) ApproxCount() int {
	z := 0.0
	for _, v := range h.counters {
		z += math.Pow(2, -float64(v))
	}
	z = 1.0 / z
	return int(alpha * m * m * z)
}

// nzeros counts the number of zeros on the right side of a binary number.
func nzeros(a uint64) int {
	if a == 0 {
		return 64 - nbits // Number of bits after using the first nbits.
	}
	n := 0
	for a&1 == 0 {
		n++
		a /= 2
	}
	return n
}
