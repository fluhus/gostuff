package rhash

import (
	"crypto/rand"
	"fmt"
	"testing"

	"github.com/fluhus/gostuff/hashx"
)

func BenchmarkWrite1K(b *testing.B) {
	const n = 20
	buf := make([]byte, 1024)
	rand.Read(buf)
	b.Run("buz64", func(b *testing.B) {
		h := NewBuz(n)
		for b.Loop() {
			h.Write(buf)
		}
	})
	b.Run("rabin32", func(b *testing.B) {
		h := NewRabinFingerprint32(n)
		for b.Loop() {
			h.Write(buf)
		}
	})
	b.Run("rabin64", func(b *testing.B) {
		h := NewRabinFingerprint64(n)
		for b.Loop() {
			h.Write(buf)
		}
	})
}

func BenchmarkRolling(b *testing.B) {
	text := make([]byte, 10000)
	rand.Read(text)
	for _, ln := range []int{10, 30, 100} {
		b.Run(fmt.Sprint("buz", ln), func(b *testing.B) {
			h := NewBuz(ln)
			for b.Loop() {
				h.Write(text[:ln-1])
				var s uint64
				for _, b := range text[ln:] {
					h.WriteByte(b)
					s += h.Sum64()
				}
				if s == 0 {
					b.Error("placeholder to not optimize s out")
				}
			}
		})
		b.Run(fmt.Sprint("rabin", ln), func(b *testing.B) {
			h := NewRabinFingerprint64(ln)
			for b.Loop() {
				h.Write(text[:ln-1])
				var s uint64
				for _, b := range text[ln:] {
					h.WriteByte(b)
					s += h.Sum64()
				}
				if s == 0 {
					b.Error("placeholder to not optimize s out")
				}
			}
		})
		b.Run(fmt.Sprint("murmur3", ln), func(b *testing.B) {
			for b.Loop() {
				var s uint64
				for i := range text[ln:] {
					s += hashx.Bytes(text[i : i+ln])
				}
				if s == 0 {
					b.Error("placeholder to not optimize s out")
				}
			}
		})
	}
}
