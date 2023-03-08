package rhash

import (
	"hash"
	"testing"
)

func TestBuz32(t *testing.T) {
	test32(t, func(i int) hash.Hash32 { return NewBuz32(i) })
}
