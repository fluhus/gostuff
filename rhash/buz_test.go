package rhash

import (
	"hash"
	"testing"
)

func TestBuz64(t *testing.T) {
	test64(t, func(i int) hash.Hash64 { return NewBuz(i) })
}

func TestBuz32(t *testing.T) {
	test32(t, func(i int) hash.Hash32 { return NewBuz(i) })
}

func TestBuz64_seed(t *testing.T) {
	seed := BuzRandomSeed()
	test64(t, func(i int) hash.Hash64 {
		return NewBuzWithSeed(i, seed)
	})
}

func TestBuz64_modifySeed(t *testing.T) {
	seed := &BuzSeed{}
	for i := range seed {
		seed[i] = uint64(i)
	}
	input := "dsfhjkdfhdjsfjksdjadi"
	h := NewBuzWithSeed(len(input), seed)
	h.Write([]byte(input))
	want := h.Sum64()
	for i := range seed {
		seed[i] = 0
	}
	h.Write([]byte(input))
	got := h.Sum64()
	if got != want {
		t.Fatalf("Buz.Sum64(%q)=%v, want %v", input, got, want)
	}
}
