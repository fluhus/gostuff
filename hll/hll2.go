package hll

import (
	"fmt"
	"math"
)

// An HLL2 is a HyperLogLog counter for arbitrary values.
type HLL2[T any] struct {
	counters []byte
	h        func(T) uint64
	nbits    int
	m        int
	mask     uint64
}

// New2 creates a new HyperLogLog counter.
// The counter will use 2^logSize bytes.
// h is the hash function to use for added values.
func New2[T any](logSize int, h func(T) uint64) *HLL2[T] {
	if logSize < 4 {
		panic(fmt.Sprintf("logSize=%v, should be at least 4", logSize))
	}
	m := 1 << logSize
	return &HLL2[T]{
		counters: make([]byte, m),
		h:        h,
		nbits:    logSize,
		m:        m,
		mask:     uint64(m - 1),
	}
}

// Add adds v to the counter. Calls hash once.
func (h *HLL2[T]) Add(t T) {
	hash := h.h(t)
	idx := hash & h.mask
	fp := hash >> h.nbits
	z := byte(h.nzeros(fp)) + 1
	if z > h.counters[idx] {
		h.counters[idx] = z
	}
}

// ApproxCount returns the current approximate count.
// Does not alter the state of the counter.
func (h *HLL2[T]) ApproxCount() int {
	z := 0.0
	for _, v := range h.counters {
		z += math.Pow(2, -float64(v))
	}
	z = 1.0 / z
	fm := float64(h.m)
	result := int(h.alpha() * fm * fm * z)

	if result < h.m*5/2 {
		zeros := 0
		for _, v := range h.counters {
			if v == 0 {
				zeros++
			}
		}
		// If some registers are zero, use linear counting.
		if zeros > 0 {
			result = int(fm * math.Log(fm/float64(zeros)))
		}
	}

	return result
}

// Returns the alpha value to use depending on m.
func (h *HLL2[T]) alpha() float64 {
	switch h.m {
	case 16:
		return 0.673
	case 32:
		return 0.697
	case 64:
		return 0.709
	}
	return 0.7213 / (1 + 1.079/float64(h.m))
}

// nzeros counts the number of zeros on the right side of a binary number.
func (h *HLL2[T]) nzeros(a uint64) int {
	if a == 0 {
		return 64 - h.nbits // Number of bits after using the first nbits.
	}
	n := 0
	for a&1 == 0 {
		n++
		a /= 2
	}
	return n
}

// AddHLL adds the state of another counter to h,
// assuming they use the same hash function.
// The result is equivalent to adding all the values of other to h.
func (h *HLL2[T]) AddHLL(other *HLL2[T]) {
	if len(h.counters) != len(other.counters) {
		panic("merging HLLs with different sizes")
	}
	for i, b := range other.counters {
		if h.counters[i] < b {
			h.counters[i] = b
		}
	}
}
