package rhash

import (
	"fmt"
	"hash"
)

const rabinPrime = 16777619

var _ hash.Hash64 = &RabinFingerprint64{}

// RabinFingerprint64 implements a Rabin fingerprint rolling-hash.
// Implements [hash.Hash64].
type RabinFingerprint64 struct {
	h, pow uint64
	i      int
	hist   []byte
}

// NewRabinFingerprint64 returns a new rolling hash with a window size of n.
func NewRabinFingerprint64(n int) *RabinFingerprint64 {
	if n < 1 {
		panic(fmt.Sprintf("bad n: %d", n))
	}
	return &RabinFingerprint64{0, 1, 0, make([]byte, n)}
}

// Write updates the hash with the given bytes. Always returns len(data), nil.
func (h *RabinFingerprint64) Write(data []byte) (int, error) {
	for _, b := range data {
		h.WriteByte(b)
	}
	return len(data), nil
}

// WriteByte updates the hash with the given byte. Always returns nil.
func (h *RabinFingerprint64) WriteByte(b byte) error {
	h.h = h.h*rabinPrime + uint64(b)
	i := h.i % len(h.hist)
	h.h -= h.pow * uint64(h.hist[i])
	h.hist[i] = b
	if h.i < len(h.hist) {
		h.pow *= rabinPrime
	}
	h.i++
	return nil
}

// Sum appends the current hash to b and returns the resulting slice.
func (h *RabinFingerprint64) Sum(b []byte) []byte {
	s := h.Sum64()
	n := h.Size()
	for i := 0; i < n; i++ {
		b = append(b, byte(s))
		s >>= 8
	}
	return b
}

// Sum64 returns the current hash.
func (h *RabinFingerprint64) Sum64() uint64 {
	return h.h
}

// Size returns the number of bytes Sum will return, which is eight.
func (h *RabinFingerprint64) Size() int {
	return 8
}

// BlockSize returns the hash's block size, which is one.
func (h *RabinFingerprint64) BlockSize() int {
	return 1
}

// Reset resets the hash to its initial state.
func (h *RabinFingerprint64) Reset() {
	h.h = 0
	h.i = 0
	h.pow = 1
	for i := range h.hist {
		h.hist[i] = 0
	}
}

// RabinFingerprintSum64 returns the Rabin fingerprint of data.
func RabinFingerprintSum64(data []byte) uint64 {
	if len(data) == 0 {
		return 0
	}
	h := NewRabinFingerprint64(len(data))
	h.Write(data)
	return h.Sum64()
}
