package vectors

import (
	"testing"
)

func TestL1(t *testing.T) {
	v1 := []float64{0, 0, 0}
	v2 := []float64{1, 0, 0}
	v3 := []float64{0, 1, 1}
	v4 := []float64{1, 1, 1}

	assertEqualFloats(L1(v1, v1), 0.0, t)
	assertEqualFloats(L1(v1, v2), 1.0, t)
	assertEqualFloats(L1(v1, v3), 2.0, t)
	assertEqualFloats(L1(v1, v4), 3.0, t)
	assertEqualFloats(L1(v2, v4), 2.0, t)
	assertEqualFloats(L1(v3, v4), 1.0, t)
}

func TestL2(t *testing.T) {
	v1 := []float64{0, 0}
	v2 := []float64{0, 1}
	v3 := []float64{4, 4}

	assertEqualFloats(L2(v1, v1), 0.0, t)
	assertEqualFloats(L2(v1, v2), 1.0, t)
	assertEqualFloats(L2(v3, v2), 5.0, t)
}

// TODO(amit): Add more tests.

func assertEqualFloats(actual, expected float64, t *testing.T) {
	if actual != expected {
		t.Fatal("Bad value returned. Expected:", expected, "Actual:", actual)
	}
}
