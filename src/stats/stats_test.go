package stats

import (
	"math"
	"testing"
)

func TestSum(t *testing.T) {
	assertEqualFloats(Sum([]float64{1, 0, 4, 3.5, -7}), 1.5, t)
}

func TestMean(t *testing.T) {
	assertEqualFloats(Mean([]float64{1, 0, 4, 3.5, -7}), 0.3, t)
}

func TestCov(t *testing.T) {
	a := []float64{1, 2, 3, 4}
	b := []float64{9, 8, 7, 6}
	assertEqualFloats(Cov(a, b), -1.25, t)
	assertEqualFloats(Cov(b, a), -1.25, t)
}

func TestVarStd(t *testing.T) {
	assertEqualFloats(Var([]float64{1, 2, 3, 4, 5}), 2, t)
	assertEqualFloats(Std([]float64{1, 2, 3, 4, 5}), math.Sqrt(2), t)
}

func TestMinMaxSpan(t *testing.T) {
	a := []float64{5, 3, 4, 1, 5, 2, 6, 5, 3, 1, 3, 5, 3}
	assertEqualFloats(Min(a), 1, t)
	assertEqualFloats(Max(a), 6, t)
	assertEqualFloats(Span(a), 5, t)
}

func TestEnt(t *testing.T) {
	assertEqualFloats(Ent([]float64{1, 1, 1, 1, 1, 1, 1, 1}), 3, t)
}

func assertEqualFloats(actual, expected float64, t *testing.T) {
	if actual != expected {
		t.Fatal("Bad value returned. Expected:", expected, "Actual:", actual)
	}
}
