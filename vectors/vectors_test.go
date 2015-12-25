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

func TestAdd(t *testing.T) {
	v1 := []float64{1, 2, 3}
	v2 := []float64{4, 5, 6}
	v3 := []float64{5, 7, 9}
	assertEqualFloatSlices(Add(v1, v2), v3, t)
	assertEqualFloatSlices(v1, v3, t)
}

func TestSub(t *testing.T) {
	v1 := []float64{1, 2, 3}
	v2 := []float64{4, 5, 6}
	v3 := []float64{-3, -3, -3}
	assertEqualFloatSlices(Sub(v1, v2), v3, t)
	assertEqualFloatSlices(v1, v3, t)
}

func TestMul(t *testing.T) {
	v1 := []float64{1, 2, 3}
	v2 := []float64{11, 22, 33}
	assertEqualFloatSlices(Mul(v1, 11), v2, t)
	assertEqualFloatSlices(v1, v2, t)
}

func TestDot(t *testing.T) {
	v1 := []float64{1, 2, 3}
	v2 := []float64{11, 22, 33}
	assertEqualFloats(Dot(v1, v2), 154, t)
}

func TestOnes(t *testing.T) {
	assertEqualFloatSlices(Ones(5), []float64{1, 1, 1, 1, 1}, t)
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
