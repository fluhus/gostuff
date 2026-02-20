// Package hashx provides simple hashing functions for various input types.
package hashx

import (
	"hash"
	"unsafe"

	"github.com/spaolacci/murmur3"
	"golang.org/x/exp/constraints"
)

// Hashx calculates hash values for various input types.
type Hashx struct {
	h   hash.Hash64
	buf []byte
}

// NewSeed returns a new Hashx with the given seed.
func NewSeed(seed uint32) *Hashx {
	return &Hashx{murmur3.New64WithSeed(seed), make([]byte, 8)}
}

// New returns a new Hashx.
func New() *Hashx {
	return &Hashx{murmur3.New64(), make([]byte, 8)}
}

// Bytes returns the hash value of the given byte sequence.
func (h *Hashx) Bytes(b []byte) uint64 {
	h.h.Reset()
	h.h.Write(b)
	return h.h.Sum64()
}

// String returns the hash value of the given string.
func (h *Hashx) String(s string) uint64 {
	return h.Bytes(unsafe.Slice(unsafe.StringData(s), len(s)))
}

// IntHashx returns the hash value of the given integer,
// using the given Hashx instance.
func IntHashx[I constraints.Integer](h *Hashx, i I) uint64 {
	h.h.Reset()
	u := uint64(i)
	for j := range 8 {
		h.buf[j] = byte(u >> (j * 8))
	}
	h.h.Write(h.buf)
	return h.h.Sum64()
}

// The default hashx.
var dflt = New()

// Bytes returns the hash value of the given byte sequence.
func Bytes(b []byte) uint64 { return dflt.Bytes(b) }

// String returns the hash value of the given string.
func String(s string) uint64 { return dflt.String(s) }

// Int returns the hash value of the given integer.
func Int[I constraints.Integer](i I) uint64 { return IntHashx(dflt, i) }
