package hll

import (
	"fmt"
	"math"
	"testing"

	"github.com/fluhus/gostuff/hashx"
)

func TestCount_short(t *testing.T) {
	upto := 10000000
	if testing.Short() {
		upto = 1000
	}

	hll := NewComparable[int](16)
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
	hll := NewComparable[int](16)
	if count := hll.ApproxCount(); count != 0 {
		t.Fatalf("ApproxCount()=%v, want 0", count)
	}
}

func TestAddHLL(t *testing.T) {
	hll1 := NewComparable[int](16)
	for i := 1; i <= 5; i++ {
		hll1.Add(i)
	}
	if count := hll1.ApproxCount(); count != 5 {
		t.Fatalf("ApproxCount()=%v, want 5", count)
	}

	hll2 := NewComparable[int](16)
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
	b.Run("id", func(b *testing.B) {
		hll := New(16, func(i int) uint64 { return uint64(i) })
		for i := 0; b.Loop(); i++ {
			hll.Add(i)
		}
	})
	b.Run("hashx", func(b *testing.B) {
		hll := New[int](16, hashx.Int)
		for i := 0; b.Loop(); i++ {
			hll.Add(i)
		}
	})
	b.Run("comparable", func(b *testing.B) {
		hll := NewComparable[int](16)
		for i := 0; b.Loop(); i++ {
			hll.Add(i)
		}
	})
}

func BenchmarkCount(b *testing.B) {
	const nelements = 1000000
	for _, n := range []int{10, 12, 14, 16} {
		b.Run(fmt.Sprint(n), func(b *testing.B) {
			hll := NewComparable[int](n)
			for i := range nelements {
				hll.Add(i)
			}
			for b.Loop() {
				hll.ApproxCount()
			}
		})
	}
}
