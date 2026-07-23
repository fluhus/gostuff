// Package bloom provides a simple bloom filter implementation.
package bloom

import (
	"encoding/binary"
	"fmt"
	"hash"
	"iter"
	"math"
	"unsafe"

	"github.com/fluhus/gostuff/bits"
	"github.com/fluhus/gostuff/snm"
	"github.com/spaolacci/murmur3"
	"golang.org/x/exp/constraints"
)

// Filter is a single bloom filter.
//
// A single filter is unsafe for asynchronous calls.
// Calls to Add and Has should be synchronized from the outside.
type Filter[T any] struct {
	b []byte // Filter data.
	h []hash.Hash64
	e func(T) []byte // Encodes T for hash.
}

// NHash returns the number of hash functions this filter uses.
func (f *Filter[T]) NHash() int {
	return len(f.h)
}

// NBits returns the number of bits this filter uses.
func (f *Filter[T]) NBits() int {
	return 8 * len(f.b)
}

// NElements returns an approximation of the number of elements added to the
// filter.
func (f *Filter[T]) NElements() int {
	m := float64(f.NBits())
	k := float64(f.NHash())
	x := float64(bits.Sum(f.b))
	return int(math.Round(-m / k * math.Log(1-x/m)))
}

// Iterates over the hashes of v, mod nbits.
func (f *Filter[T]) hash(v T) iter.Seq[int] {
	return func(yield func(int) bool) {
		nbits := uint64(f.NBits())
		b := f.e(v)
		for _, h := range f.h {
			h.Reset()
			h.Write(b)
			if !yield(int(h.Sum64() % nbits)) {
				return
			}
		}
	}
}

// Has checks if all k hash values of v were encountered.
// Makes at most k hash calculations.
func (f *Filter[T]) Has(v T) bool {
	for h := range f.hash(v) {
		if bits.Get(f.b, h) == 0 {
			return false
		}
	}
	return true
}

// Add adds v to the filter, and returns the value of Has(v) before adding.
// After calling Add, Has(v) will always be true. Makes k calls to hash.
func (f *Filter[T]) Add(v T) bool {
	has := true
	for h := range f.hash(v) {
		if bits.Get(f.b, h) == 0 {
			has = false
			bits.Set(f.b, h, true)
		}
	}
	return has
}

// AddFilter merges other into f. After merging, f is equivalent to have been added
// all the elements of other.
func (f *Filter[T]) AddFilter(other *Filter[T]) {
	// Make sure the two filters are compatible.
	if f.NBits() != other.NBits() {
		panic(fmt.Sprintf("mismatching number of bits: this has %v, other has %v",
			f.NBits(), other.NBits()))
	}
	if f.NHash() != other.NHash() {
		panic(fmt.Sprintf("mismatching number of hashes: this has %v, other has %v",
			f.NHash(), other.NHash()))
	}

	// Merge.
	for i := range f.b {
		f.b[i] |= other.b[i]
	}
}

// New creates a new bloom filter with m bits and k hash functions.
// Number of bits is rounded up to the nearest multiple of 8.
// e encodes an element into a slice of bytes.
//
// See [OptimalParams] for an alternative way to decide on the parameters.
func New[T any](m int, k int, e func(T) []byte) *Filter[T] {
	if m < 1 {
		panic(fmt.Sprintf("number of bits should be at least 1, got %v", m))
	}
	if k < 1 {
		panic(fmt.Sprintf("k should be at least 1, got %v", k))
	}

	result := &Filter[T]{
		b: make([]byte, ((m-1)/8)+1),
		h: snm.Slice(k, func(i int) hash.Hash64 {
			return murmur3.New64WithSeed(uint32(i))
		}),
		e: e,
	}
	return result
}

// OptimalParams returns the optimal bits (m) and hashes (k)
// for the given element count (n) and false-positive rate (p).
//
// The calculation is taken from:
// https://en.wikipedia.org/wiki/Bloom_filter#Optimal_number_of_hash_functions
func OptimalParams(n int, p float64) (m, k int) {
	m = int(math.Round(-float64(n) * math.Log(p) / math.Ln2 / math.Ln2))
	k = int(math.Round(-math.Log2(p)))
	return
}

// Int returns a function that encodes integers.
func Int[T constraints.Integer]() func(T) []byte {
	b := make([]byte, 8)
	return func(i T) []byte {
		binary.LittleEndian.PutUint64(b, uint64(i))
		return b
	}
}

// Float returns a function that encodes floats.
func Float[T constraints.Float]() func(T) []byte {
	b := make([]byte, 8)
	return func(f T) []byte {
		binary.LittleEndian.PutUint64(b, math.Float64bits(float64(f)))
		return b
	}
}

// Bytes is the identity function.
func Bytes(b []byte) []byte {
	return b
}

// String encodes strings to bytes.
func String(s string) []byte {
	return unsafe.Slice(unsafe.StringData(s), len(s))
}
