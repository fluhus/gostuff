// Package hll provides an implementation of the HyperLogLog algorithm.
//
// A HyperLogLog counter can approximate the cardinality of a set with high
// accuracy and little memory.
//
// Accuracy
//
// The counter is built to be accurate up to +-1% for any cardinality starting
// from 0, with a high probability. This is verified in the tests.
//
// Performance
//
// An HLL counter uses 65kb memory. Adding an element of size up to 100
// bytes takes an order of 100ns. Calculating the approximate count takes an
// order of 4ms.
//
// Citation
//
// Flajolet, Philippe; Fusy, Éric; Gandouet, Olivier; Meunier, Frédéric (2007).
// "Hyperloglog: The analysis of a near-optimal cardinality estimation
// algorithm". Discrete Mathematics and Theoretical Computer Science
// Proceedings.
package hll

import (
	"encoding/json"
	"fmt"
	"hash"
	"math"
	_ "unsafe"

	"github.com/spaolacci/murmur3"
)

const (
	nbits = 16
	m     = 1 << nbits
	mask  = m - 1
	alpha = 0.7213 / (1.0 + 1.079/m)
)

//go:linkname fastrand runtime.fastrand
func fastrand() uint32

// An HLL is a HyperLogLog counter for arbitrary values.
type HLL struct {
	counters []byte
	h        hash.Hash64
	seed     uint32
}

// New creates a new HyperLogLog counter with a random hash seed.
func New() *HLL {
	return NewSeed(fastrand())
}

// NewSeed creates a new HyperLogLog counter with the given hash seed.
func NewSeed(seed uint32) *HLL {
	return &HLL{
		counters: make([]byte, m),
		h:        murmur3.New64WithSeed(seed),
		seed:     seed,
	}
}

// Add adds v to the counter. Calls hash once.
func (h *HLL) Add(v []byte) {
	h.h.Reset()
	h.h.Write(v)
	hash := h.h.Sum64()

	idx := hash & mask
	fp := hash >> nbits
	z := byte(nzeros(fp)) + 1
	if z > h.counters[idx] {
		h.counters[idx] = z
	}
}

// ApproxCount returns the current approximate count.
// Does not alter the state of the counter.
func (h *HLL) ApproxCount() int {
	z := 0.0
	for _, v := range h.counters {
		z += math.Pow(2, -float64(v))
	}
	z = 1.0 / z
	result := int(alpha * m * m * z)

	if result < m*5/2 {
		zeros := 0
		for _, v := range h.counters {
			if v == 0 {
				zeros++
			}
		}
		// If some registers are zero, use linear counting.
		if zeros > 0 {
			result = int(m * math.Log(m/float64(zeros)))
		}
	}

	return result
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

// AddHLL adds the state of another counter to h.
// The result is equivalent to adding all the values of other to h.
func (h *HLL) AddHLL(other *HLL) {
	if h.seed != other.seed {
		panic(fmt.Sprintf("seeds don't match: %v, %v", h.seed, other.seed))
	}
	for i, b := range other.counters {
		if h.counters[i] < b {
			h.counters[i] = b
		}
	}
}

// Used for JSON marshaling/unmarshaling.
type jsonHLL struct {
	Counters []byte
	Seed     uint32
}

// MarshalJSON implements the json.Marshaler interface.
func (h *HLL) MarshalJSON() ([]byte, error) {
	return json.Marshal(&jsonHLL{Counters: h.counters, Seed: h.seed})
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (h *HLL) UnmarshalJSON(b []byte) error {
	jh := &jsonHLL{}
	if err := json.Unmarshal(b, jh); err != nil {
		return err
	}
	h.counters = jh.Counters
	h.h = murmur3.New64WithSeed(jh.Seed)
	h.seed = jh.Seed
	return nil
}
