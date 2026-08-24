//go:build !go1.27

// Old int functions, deprecated in Go 1.27.

package hashx

import "golang.org/x/exp/constraints"

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

// Int returns the hash value of the given integer.
func Int[I constraints.Integer](i I) uint64 { return IntHashx(dflt, i) }
