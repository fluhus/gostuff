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

func TestHist(t *testing.T) {
	a := []float64{1, 2, 3, 2, 3, 2, 1, 3, 3, 3, 3, 3, 3, 2, 1, 1, 1}
	x, y, z := Hist(a)

	if len(x) != 3 || x[1] != 5 || x[2] != 4 || x[3] != 8 {
		t.Fatal("Bad value returned. Expected:", map[float64]int{1: 5, 2: 4, 3: 8},
			"Actual:", x)
	}

	assertEqualFloatSlices(y, []float64{1, 2, 3}, t)
	assertEqualFloatSlices(z, []float64{3, 1, 2}, t)
}

func assertEqualFloats(actual, expected float64, t *testing.T) {
	if actual != expected {
		t.Fatal("Bad value returned. Expected:", expected, "Actual:", actual)
	}
}

func assertEqualFloatSlices(actual, expected []float64, t *testing.T) {
	if len(actual) != len(expected) {
		t.Fatal("Bad value returned. Expected:", expected, "Actual:", actual)
	}
	for i := range actual {
		if actual[i] != expected[i] {
			t.Fatal("Bad value returned. Expected:", expected, "Actual:", actual)
		}
	}
}
