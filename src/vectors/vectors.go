// Handles basic vector operations.
package vectors

import (
	"math"
	"fmt"
)

// L1 (Manhattan) distance. Equivalent to Lp(1) but more efficient.
func L1(a, b []float64) float64 {
	assertMatchingLengths(a, b)

	sum := 0.0
	for i := range a {
		sum += math.Abs(a[i] - b[i])
	}

	return sum
}

// L2 (Euclidean) distance. Equivalent to Lp(2) but more efficient.
func L2(a, b []float64) float64 {
	assertMatchingLengths(a, b)

	sum := 0.0
	for i := range a {
		d := (a[i] - b[i])
		sum += d * d
	}

	return math.Sqrt(sum)
}

// Returns an Lp distance function. Lp is calculated as follows:
//  Lp(v) = the p'th root of sum_i(v[i]^p)
//
// For convenience, L1 (Manhattan)and L2 (Euclidean) are prepared package
// variables.
func Lp(p int) func([]float64, []float64) float64 {
	if p < 1 {
		panic(fmt.Sprintf("Invalid p: %d", p))
	}
	
	// Prepared functions.
	if p == 1 {
		return L1
	}
	
	if p == 2 {
		return L2
	}

	return func(a, b []float64) float64 {
		assertMatchingLengths(a, b)

		fp := float64(p)
		sum := 0.0
		for i := range a {
			sum += math.Pow(math.Abs(a[i] - b[i]), fp)
		}

		return math.Pow(sum, 1/fp)
	}
}

// Adds b to a. b is unchanged.
func Add(a, b []float64) {
	assertMatchingLengths(a, b)
	for i := range a {
		a[i] += b[i]
	}
}

// Subtracts b from a. b is unchanged.
func Sub(a, b []float64) {
	assertMatchingLengths(a, b)
	for i := range a {
		a[i] -= b[i]
	}
}

// Multiplies the values of a by m.
func Mul(a []float64, m float64) {
	for i := range a {
		a[i] *= m
	}
}

// Returns the dot product of the input vectors.
func Dot(a, b []float64) float64 {
	assertMatchingLengths(a, b)
	
	sum := 0.0
	for i := range a {
		sum += a[i] * b[i]
	}
	return sum
}

// Returns the norm of the vector, in L2.
func Norm(a []float64) float64 {
	norm := 0.0
	for _, v := range a {
		norm += v*v
	}
	return math.Sqrt(norm)
}

// Returns a copy of the given vector.
func Copy(a []float64) []float64 {
	result := make([]float64, len(a))
	copy(result, a)
	return result
}

// Panics if 2 vectors are of inequal lengths.
func assertMatchingLengths(a, b []float64) {
	if len(a) != len(b) {
		panic(fmt.Sprintf("Mismatching lengths: %d, %d.", len(a), len(b)))
	}
}
