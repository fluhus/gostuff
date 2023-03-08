package rhash

import (
	"hash"
	"testing"
)

func TestBuz64(t *testing.T) {
	test64(t, func(i int) hash.Hash64 { return NewBuz64(i) })
}
