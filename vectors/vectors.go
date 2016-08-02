// Basic vector operations.
//
// Example 1
//
//  a, b, c, d := {some vectors}
//
// Giving a nil result vector to Sum() allocates a new result, and thus does
// not change a, b, c or d.
//  dist := L2(Sum(nil, a, b), Sum(nil, c, d))
//
// Giving a non-nil result vector to Sum() changes it. The following call will
// give the same distance but will modify a and c.
//  dist := L2(Sum(a, b), Sum(c, d))
//
// Example 2
//
// Create a vector of -1's of length 10:
//  v := Mul(Ones(10), -1)
package vectors

import (
	"fmt"
	"math"
)

// L1 (Manhattan) distance. Equivalent to Lp(1).
func L1(a, b []float64) float64 {
	assertMatchingLengths(a, b)

	sum := 0.0
	for i := range a {
		sum += math.Abs(a[i] - b[i])
	}

	return sum
}

// L2 (Euclidean) distance. Equivalent to Lp(2).
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
			sum += math.Pow(math.Abs(a[i]-b[i]), fp)
		}

		return math.Pow(sum, 1/fp)
	}
}

// Adds b to a and returns a. b is unchanged. If a is nil, creates a new vector.
func Add(a []float64, b ...[]float64) []float64 {
	if a == nil {
		if len(b) == 0 {
			return nil
		} else {
			a = make([]float64, len(b[0]))
		}
	}

	for i := range b {
		assertMatchingLengths(a, b[i])
		for j := range a {
			a[j] += b[i][j]
		}
	}
	return a
}

// Subtracts b from a and returns a. b is unchanged. If a is nil, creates a new
// vector.
func Sub(a []float64, b ...[]float64) []float64 {
	if a == nil {
		if len(b) == 0 {
			return nil
		} else {
			a = make([]float64, len(b[0]))
		}
	}

	for i := range b {
		assertMatchingLengths(a, b[i])
		for j := range a {
			a[j] -= b[i][j]
		}
	}
	return a
}

// Multiplies the values of a by m and returns a.
func Mul(a []float64, m float64) []float64 {
	for i := range a {
		a[i] *= m
	}
	return a
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
		norm += v * v
	}
	return math.Sqrt(norm)
}

// Returns a slice of ones, of length n. Panics if n is negative.
func Ones(n int) []float64 {
	if n < 0 {
		panic(fmt.Sprintf("Bad vector length: %d", n))
	}
	a := make([]float64, n)
	for i := range a {
		a[i] = 1
	}
	return a
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
