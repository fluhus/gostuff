package clustering

import (
	"math"
	"reflect"
	"testing"
)

// TODO(amit): Add more test cases.

func TestClink(t *testing.T) {
	points := []float64{0, 1, 5}
	steps := []AggloStep{
		{0, 1, 1},
		{1, 2, 5},
	}
	agg := clink(len(points), func(i, j int) float64 {
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

func TestSlink(t *testing.T) {
	points := []float64{0, 1, 5}
	steps := []AggloStep{
		{0, 1, 1},
		{1, 2, 4},
	}
	agg := slink(len(points), func(i, j int) float64 {
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
