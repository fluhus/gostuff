// Package bloom provides a simple bloom filter implementation.
package bloom

// TODO(amit): Add merge function.

import (
	"fmt"
	"hash"
	"io"
	"math"
	_ "unsafe"

	"github.com/fluhus/gostuff/binio"
	"github.com/spaolacci/murmur3"
)

//go:linkname fastrand runtime.fastrand
func fastrand() uint32

// Filter is a single bloom filter.
type Filter struct {
	b    []byte        // Filter data.
	h    []hash.Hash64 // Hash functions.
	seed uint32
}

// NHash returns the number of hash functions this filter uses.
func (f *Filter) NHash() int {
	return len(f.h)
}

// NBits returns the number of bits this filter uses.
func (f *Filter) NBits() int {
	return 8 * len(f.b)
}

// Has checks if all k hash values of v were encountered.
// Makes at most k hash calculations.
func (f *Filter) Has(v []byte) bool {
	for i := range f.h {
		f.h[i].Reset()
		f.h[i].Write(v)
		hash := int(f.h[i].Sum64() % uint64(len(f.b)*8))
		if binio.GetBit(f.b, hash) == 0 {
			return false
		}
	}
	return true
}

// Add adds v to the filter, and returns the value of Has(v) before adding.
// After calling Add, Has(v) will always be true. Makes k calls to hash.
func (f *Filter) Add(v []byte) bool {
	has := true
	for i := range f.h {
		f.h[i].Reset()
		f.h[i].Write(v)
		hash := int(f.h[i].Sum64() % uint64(len(f.b)*8))
		if binio.GetBit(f.b, hash) == 0 {
			has = false
			binio.SetBit(f.b, hash, 1)
		}
	}
	return has
}

// Seed returns the hash seed of this filter.
// A new filter starts with a random seed.
func (f *Filter) Seed() uint32 {
	return f.seed
}

// SetSeed sets the hash seed of this filter.
// The filter must be empty.
func (f *Filter) SetSeed(seed uint32) {
	for _, b := range f.b {
		if b != 0 {
			panic("cannot change seed after elements were added")
		}
	}
	f.seed = seed
	h := murmur3.New32WithSeed(seed)
	for i := range f.h {
		h.Write([]byte{1})
		f.h[i] = murmur3.New64WithSeed(h.Sum32())
	}
}

// Encode writes this filter to the stream. Can be reproduced later with Decode.
func (f *Filter) Encode(w io.Writer) error {
	// Order is k, seed, bytes.
	if err := binio.WriteUint64(w, uint64(len(f.h))); err != nil {
		return err
	}
	if err := binio.WriteUint64(w, uint64(f.seed)); err != nil {
		return err
	}
	if err := binio.WriteBytes(w, f.b); err != nil {
		return err
	}
	return nil
}

// Decode reads an encoded filter from the stream and sets this filter's state
// to match it. Destroys the previously existing state of this filter.
func (f *Filter) Decode(r io.Reader) error {
	k, err := binio.ReadUint64(r)
	if err != nil {
		return err
	}
	f.h = make([]hash.Hash64, k)
	seed, err := binio.ReadUint64(r)
	if err != nil {
		return err
	}
	if seed >= 1<<32 {
		return fmt.Errorf("seed is too large: %v, want at most %v",
			seed, 1<<32-1)
	}
	f.b = nil
	f.SetSeed(uint32(seed))
	b, err := binio.ReadBytes(r)
	if err != nil {
		return err
	}
	f.b = b

	return nil
}

// New creates a new bloom filter with the given parameters. Number of
// bits is rounded up to the nearest multiple of 8.
//
// See NewOptimal for an alternative way to decide on the parameters.
func New(bits int, k int) *Filter {
	if bits < 1 {
		panic(fmt.Sprintf("number of bits should be at least 1, got %v", bits))
	}
	if k < 1 {
		panic(fmt.Sprintf("k should be at least 1, got %v", k))
	}

	result := &Filter{
		b: make([]byte, ((bits-1)/8)+1),
		h: make([]hash.Hash64, k),
	}
	result.SetSeed(fastrand())
	return result
}

// NewOptimal creates a new bloom filter, with parameters optimal for the
// expected number of elements (n) and the required false-positive rate (p).
//
// The calculation is taken from:
// https://en.wikipedia.org/wiki/Bloom_filter#Optimal_number_of_hash_functions
func NewOptimal(n int, p float64) *Filter {
	m := math.Round(-float64(n) * math.Log(p) / math.Ln2 / math.Ln2)
	k := math.Round(-math.Log2(p))
	return New(int(m), int(k))
}
