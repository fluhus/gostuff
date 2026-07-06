//go:build go1.27

package hashx

import (
	"encoding/binary"
	"unsafe"

	"golang.org/x/exp/constraints"
)

// Int returns the hash value of the given integer.
func (h *Hashx) Int[I constraints.Integer](i I) uint64 {
	return h.intBinary(i)
}

// Several ways to hash an int - still considered.

func (h *Hashx) intShifts[I constraints.Integer](i I) uint64 {
	h.h.Reset()
	u := uint64(i)
	for j := range 8 {
		h.buf[j] = byte(u >> (j * 8))
	}
	h.h.Write(h.buf)
	return h.h.Sum64()
}

func (h *Hashx) intUnsafe[I constraints.Integer](i I) uint64 {
	h.h.Reset()
	u := uint64(i)
	a := *(*[8]byte)(unsafe.Pointer(&u))
	copy(h.buf, a[:])
	h.h.Write(h.buf)
	return h.h.Sum64()
}

func (h *Hashx) intBinary[I constraints.Integer](i I) uint64 {
	h.h.Reset()
	u := uint64(i)
	binary.LittleEndian.PutUint64(h.buf, u)
	h.h.Write(h.buf)
	return h.h.Sum64()
}
