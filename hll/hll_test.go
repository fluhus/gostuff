package hll

import (
	"fmt"
	"math"
	"math/rand"
	"testing"

	"github.com/fluhus/gostuff/binio"
)

func TestCount_short(t *testing.T) {
	buf := make([]byte, 8)
	hll := New()
	next := 1
	for i := 1; i <= 100_000; i++ {
		binio.Uint64ToBytes(uint64(i), buf)
		hll.Add(buf)
		if i != next { // Check only a sample.
			continue
		}
		next *= 3
		count := hll.ApproxCount()
		wantMin := int(math.Round(float64(i) * 0.99))
		wantMax := int(math.Round(float64(i) * 1.01))
		if count < wantMin || count > wantMax {
			t.Errorf("ApproxCount(%v)=%v, want %v-%v", i, count, wantMin, wantMax)
		}
	}
}

func TestCount_long(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping long test")
	}

	buf := make([]byte, 8)
	hll := New()
	next := 1
	checked := 0
	failed := 0
	for i := 1; i <= 1000_000_000; i++ {
		binio.Uint64ToBytes(uint64(i), buf)
		hll.Add(buf)
		if i != next { // Check only a sample.
			continue
		}
		next *= 3
		checked++
		count := hll.ApproxCount()
		wantMin := int(math.Round(float64(i) * 0.99))
		wantMax := int(math.Round(float64(i) * 1.01))
		if count < wantMin || count > wantMax {
			t.Logf("ApproxCount(%v)=%v, want %v-%v", i, count, wantMin, wantMax)
			failed++
		}
	}
	if failed > 1 {
		t.Error("Checked", checked, "failed", failed)
	}
}

func TestCount_zero(t *testing.T) {
	hll := New()
	if count := hll.ApproxCount(); count != 0 {
		t.Fatalf("ApproxCount()=%v, want 0", count)
	}
}

func BenchmarkAdd(b *testing.B) {
	for _, n := range []int{10, 30, 100} {
		b.Run(fmt.Sprintf("%v byte elements", n), func(b *testing.B) {
			hll := New()
			r := rand.New(rand.NewSource(0))
			bufs := make([][]byte, b.N) // Elements to add.
			for i := 0; i < b.N; i++ {
				bufs[i] = make([]byte, n)
				r.Read(bufs[i])
			}
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				hll.Add(bufs[i])
			}
		})
	}
}

func BenchmarkCount(b *testing.B) {
	const nelements = 1000000
	hll := New()
	buf := make([]byte, 8)
	for i := 0; i < nelements; i++ { // Add some random elements.
		rand.Read(buf)
		hll.Add(buf)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		hll.ApproxCount()
	}
}
