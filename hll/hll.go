// Package hll provides an implementation of the HyperLogLog algorithm.
//
// A HyperLogLog counter can approximate the cardinality of a set with high
// accuracy and little memory.
//
// # Accuracy
//
// Average error for 1,000,000,000 elements for different values of logSize:
//
//	logSize    average error %
//	4          21
//	5          12
//	6          10
//	7          8.1
//	8          4.8
//	9          3.6
//	10         1.9
//	11         1.2
//	12         1.0
//	13         0.7
//	14         0.5
//	15         0.33
//	16         0.25
//
// # Citation
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
//
// Deprecated: use HLL2.
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
