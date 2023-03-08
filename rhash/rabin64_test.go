package rhash

import (
	"hash"
	"testing"
)

func TestRabin64(t *testing.T) {
	test64(t, func(n int) hash.Hash64 { return NewRabinFingerprint64(n) })
}
