package rhash

import (
	"crypto/rand"
	"testing"
)

func BenchmarkWrite1K(b *testing.B) {
	const n = 20
	buf := make([]byte, 1024)
	rand.Read(buf)
	b.Run("buz64", func(b *testing.B) {
		h := NewBuz(n)
		for i := 0; i < b.N; i++ {
			h.Write(buf)
		}
	})
	b.Run("rabin32", func(b *testing.B) {
		h := NewRabinFingerprint32(n)
		for i := 0; i < b.N; i++ {
			h.Write(buf)
		}
	})
	b.Run("rabin64", func(b *testing.B) {
		h := NewRabinFingerprint64(n)
		for i := 0; i < b.N; i++ {
			h.Write(buf)
		}
	})
}
