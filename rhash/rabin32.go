package rhash

import (
	"fmt"
	"hash"
)

var _ hash.Hash32 = &RabinFingerprint32{}

// RabinFingerprint32 implements a Rabin fingerprint rolling-hash.
// Implements [hash.Hash32].
type RabinFingerprint32 struct {
	h, pow uint32
	i      int
	hist   []byte
}

// NewRabinFingerprint32 returns a new rolling hash with a window size of n.
func NewRabinFingerprint32(n int) *RabinFingerprint32 {
	if n < 1 {
		panic(fmt.Sprintf("bad n: %d", n))
	}
	return &RabinFingerprint32{0, 1, 0, make([]byte, n)}
}

// Write updates the hash with the given bytes. Always returns len(data), nil.
func (h *RabinFingerprint32) Write(data []byte) (int, error) {
	for _, b := range data {
		h.WriteByte(b)
	}
	return len(data), nil
}

// WriteByte updates the hash with the given byte. Always returns nil.
func (h *RabinFingerprint32) WriteByte(b byte) error {
	h.h = h.h*rabinPrime + uint32(b)
	i := h.i % len(h.hist)
	h.h -= h.pow * uint32(h.hist[i])
	h.hist[i] = b
	if h.i < len(h.hist) {
		h.pow *= rabinPrime
	}
	h.i++
	return nil
}

// Sum appends the current hash to b and returns the resulting slice.
func (h *RabinFingerprint32) Sum(b []byte) []byte {
	s := h.Sum32()
	n := h.Size()
	for i := 0; i < n; i++ {
		b = append(b, byte(s))
		s >>= 8
	}
	return b
}

// Sum32 returns the current hash.
func (h *RabinFingerprint32) Sum32() uint32 {
	return h.h
}

// Size returns the number of bytes Sum will return, which is four.
func (h *RabinFingerprint32) Size() int {
	return 4
}

// BlockSize returns the hash's block size, which is one.
func (h *RabinFingerprint32) BlockSize() int {
	return 1
}

// Reset resets the hash to its initial state.
func (h *RabinFingerprint32) Reset() {
	h.h = 0
	h.i = 0
	h.pow = 1
	for i := range h.hist {
		h.hist[i] = 0
	}
}

// RabinFingerprintSum32 returns the Rabin fingerprint of data.
func RabinFingerprintSum32(data []byte) uint32 {
	if len(data) == 0 {
		return 0
	}
	h := NewRabinFingerprint32(len(data))
	h.Write(data)
	return h.Sum32()
}
