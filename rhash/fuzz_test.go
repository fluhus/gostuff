package rhash

import (
	"hash"
	"testing"

	"github.com/fluhus/gostuff/snm"
)

func FuzzBuz64(f *testing.F) {
	fuzz64(f, func(n int) hash.Hash64 { return NewBuz64(n) })
}

func FuzzRabin64(f *testing.F) {
	fuzz64(f, func(n int) hash.Hash64 { return NewRabinFingerprint64(n) })
}

func FuzzBuz32(f *testing.F) {
	fuzz32(f, func(n int) hash.Hash32 { return NewBuz32(n) })
}

func FuzzRabin32(f *testing.F) {
	fuzz32(f, func(n int) hash.Hash32 { return NewRabinFingerprint32(n) })
}

func fuzz64(f *testing.F, fn func(n int) hash.Hash64) {
	const m = 10
	prefix := snm.Slice(30, func(i int) byte { return byte(i) })

	f.Add([]byte{})
	f.Fuzz(func(t *testing.T, a []byte) {
		a = append(prefix, a...)
		h := fn(m)
		h.Write(a[len(a)-m:])
		want := h.Sum64()
		for i := range a[m:] {
			h.Write(a[i:])
			if got := h.Sum64(); got != want {
				t.Fatalf("Sum64(%v)=%v, want %v", a[i:], got, want)
			}
		}
	})
}

func fuzz32(f *testing.F, fn func(n int) hash.Hash32) {
	const m = 10
	prefix := snm.Slice(30, func(i int) byte { return byte(i) })

	f.Add([]byte{})
	f.Fuzz(func(t *testing.T, a []byte) {
		a = append(prefix, a...)
		h := fn(m)
		h.Write(a[len(a)-m:])
		want := h.Sum32()
		for i := range a[m:] {
			h.Write(a[i:])
			if got := h.Sum32(); got != want {
				t.Fatalf("Sum32(%v)=%v, want %v", a[i:], got, want)
			}
		}
	})
}
