package hll

import (
	"math"
	"testing"

	"github.com/fluhus/gostuff/bnry"
	"github.com/spaolacci/murmur3"
)

func TestCount_short(t *testing.T) {
	upto := 10000000
	if testing.Short() {
		upto = 1000
	}

	hll := newIntHLL()
	next := 1
	ratioSum := 0.0
	ratioCount := 0.0
	for i := 1; i <= upto; i++ {
		hll.Add(i)
		if i != next { // Check only a sample.
			continue
		}
		next = (next + 1) * 21 / 20
		ratio := float64(i) / float64(hll.ApproxCount())
		ratioSum += math.Abs(math.Log(ratio))
		ratioCount++
	}
	avg := math.Exp(ratioSum / ratioCount)
	want := 1.003
	if avg > want {
		t.Errorf("average error=%f, want at most %f",
			avg, want)
	}
}

func TestCount_zero(t *testing.T) {
	hll := newIntHLL()
	if count := hll.ApproxCount(); count != 0 {
		t.Fatalf("ApproxCount()=%v, want 0", count)
	}
}

func TestAddHLL(t *testing.T) {
	hll1 := newIntHLL()
	for i := 1; i <= 5; i++ {
		hll1.Add(i)
	}
	if count := hll1.ApproxCount(); count != 5 {
		t.Fatalf("ApproxCount()=%v, want 5", count)
	}

	hll2 := newIntHLL()
	for i := 4; i <= 9; i++ {
		hll2.Add(i)
	}
	if count := hll2.ApproxCount(); count != 6 {
		t.Fatalf("ApproxCount()=%v, want 6", count)
	}

	hll1.AddHLL(hll2)
	if count := hll1.ApproxCount(); count != 9 {
		t.Fatalf("ApproxCount()=%v, want 9", count)
	}
}

func BenchmarkAdd(b *testing.B) {
	hll := New2(16, func(i int) uint64 { return uint64(i) })
	for i := 0; i < b.N; i++ {
		hll.Add(i)
	}
}

func BenchmarkAdd_intHLL(b *testing.B) {
	hll := newIntHLL()
	for i := 0; i < b.N; i++ {
		hll.Add(i)
	}
}

func BenchmarkCount(b *testing.B) {
	const nelements = 1000000
	hll := newIntHLL()
	for i := 0; i < nelements; i++ {
		hll.Add(i)
	}
	b.Run("", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			hll.ApproxCount()
		}
	})
}

func newIntHLL() *HLL2[int] {
	h := murmur3.New64()
	w := bnry.NewWriter(h)
	return New2(16, func(i int) uint64 {
		h.Reset()
		w.Write(i)
		return h.Sum64()
	})
}
