package clustering

import (
	"math"
	"reflect"
	"testing"
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
