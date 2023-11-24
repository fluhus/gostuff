package clustering

import (
	"fmt"
	"math"
	"reflect"
	"testing"

	"github.com/fluhus/gostuff/gnum"
)

func TestUPGMA(t *testing.T) {
	points := []float64{1, 4, 6, 10}
	steps := []AggloStep{
		{1, 2, 2},
		{0, 2, 4},
		{2, 3, 19.0 / 3.0},
	}
	agg := upgma(len(points), func(i, j int) float64 {
		return math.Abs(points[i] - points[j])
	})
	if agg.Len() != len(points)-1 {
		t.Fatalf("Len()=%v, want %v", agg.Len(), len(points)-1)
	}
	for i := range steps {
		if step := agg.Step(i); !reflect.DeepEqual(steps[i], step) {
			t.Errorf("Step(%v)=%v, want %v", i, step, steps[i])
		}
	}
}

func TestUPGMA_more(t *testing.T) {
	points := []float64{1, 3, 8, 12, 20, 28}
	steps := []AggloStep{
		{0, 1, 2},
		{2, 3, 4},
		{1, 3, 8},
		{4, 5, 8},
		{3, 5, 18},
	}
	agg := upgma(len(points), func(i, j int) float64 {
		return math.Abs(points[i] - points[j])
	})
	if agg.Len() != len(points)-1 {
		t.Fatalf("Len()=%v, want %v", agg.Len(), len(points)-1)
	}
	for i := range steps {
		if step := agg.Step(i); !reflect.DeepEqual(steps[i], step) {
			t.Errorf("Step(%v)=%v, want %v", i, step, steps[i])
		}
	}
}

func BenchmarkUPGMA(b *testing.B) {
	for _, n := range []int{10, 30, 100} {
		b.Run(fmt.Sprint(n), func(b *testing.B) {
			nums := make([]float64, n)
			for i := range nums {
				nums[i] = 1.0 / float64(i+1)
				if i%2 == 0 {
					nums[i] += 10
				}
			}
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				upgma(n, func(i1, i2 int) float64 {
					return gnum.Abs(nums[i1] - nums[i2])
				})
			}
		})
	}
}

func FuzzUPGMA(f *testing.F) {
	f.Add(1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0)
	f.Fuzz(func(t *testing.T, a float64, b float64, c float64,
		d float64, e float64, f float64, g float64, h float64, i float64) {
		nums := []float64{a, b, c, d, e, f, g, h, i}
		upgma(len(nums), func(i1, i2 int) float64 {
			return gnum.Abs(nums[i1] - nums[i2])
		})
	})
}
