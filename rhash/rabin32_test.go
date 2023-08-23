package rhash

import (
	"hash"
	"testing"
)

func TestRabin32(t *testing.T) {
	test32(t, func(n int) hash.Hash32 { return NewRabinFingerprint32(n) })
}
