// Package hashx provides simple hashing functions for various input types.
package hashx

import (
	"hash"
	"unsafe"

	"github.com/spaolacci/murmur3"
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

// The default hashx.
var dflt = New()

// Bytes returns the hash value of the given byte sequence.
func Bytes(b []byte) uint64 { return dflt.Bytes(b) }

// String returns the hash value of the given string.
func String(s string) uint64 { return dflt.String(s) }
