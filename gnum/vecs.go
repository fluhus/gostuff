package gnum

import (
	"fmt"
	"math"

	"golang.org/x/exp/constraints"
)

// L1 returns the L1 (Manhattan) distance between a and b.
// Equivalent to Lp(1) but returns the same type.
func L1[S ~[]N, N Number](a, b S) N {
	assertMatchingLengths(a, b)
	var sum N
	for i := range a {
		sum += Diff(a[i], b[i])
	}
	return sum
}

// L2 returns the L2 (Euclidean) distance between a and b.
// Equivalent to Lp(2).
func L2[S ~[]N, N Number](a, b S) float64 {
	assertMatchingLengths(a, b)
	var sum N
	for i := range a {
		d := (a[i] - b[i])
		sum += d * d
	}
	return math.Sqrt(float64(sum))
}

// Lp returns an Lp distance function. Lp is calculated as follows:
//
//	Lp(v) = (sum_i(v[i]^p))^(1/p)
func Lp[S ~[]N, N Number](p int) func(S, S) float64 {
	if p < 1 {
		panic(fmt.Sprintf("invalid p: %d", p))
	}

	if p == 1 {
		return func(a, b S) float64 {
			return float64(L1(a, b))
		}
	}
	if p == 2 {
		return L2[S, N]
	}

	return func(a, b S) float64 {
		assertMatchingLengths(a, b)
		fp := float64(p)
		var sum float64
		for i := range a {
			sum += math.Pow(float64(Diff(a[i], b[i])), fp)
		}
		return math.Pow(sum, 1/fp)
	}
}

// Add adds b to a and returns a. b is unchanged. If a is nil, creates a new
// vector.
func Add[S ~[]N, N Number](a S, b ...S) S {
	if a == nil {
		if len(b) == 0 {
			return nil
		}
		a = make(S, len(b[0]))
	}
	for i := range b {
		assertMatchingLengths(a, b[i])
		for j := range a {
			a[j] += b[i][j]
		}
	}
	return a
}

// Sub subtracts b from a and returns a. b is unchanged. If a is nil, creates a
// new vector.
func Sub[S ~[]N, N Number](a S, b ...S) S {
	if a == nil {
		if len(b) == 0 {
			return nil
		}
		a = make(S, len(b[0]))
	}
	for i := range b {
		assertMatchingLengths(a, b[i])
		for j := range a {
			a[j] -= b[i][j]
		}
	}
	return a
}

// Mul multiplies a by b and returns a. b is unchanged. If a is nil, creates a
// new vector.
func Mul[S ~[]N, N Number](a S, b ...S) S {
	if a == nil {
		if len(b) == 0 {
			return nil
		}
		a = Ones[S](len(b[0]))
	}
	for i := range b {
		assertMatchingLengths(a, b[i])
		for j := range a {
			a[j] -= b[i][j]
		}
	}
	return a
}

// Add1 adds m to a and returns a.
func Add1[S ~[]N, N Number](a S, m N) S {
	for i := range a {
		a[i] += m
	}
	return a
}

// Sub1 subtracts m from a and returns a.
func Sub1[S ~[]N, N Number](a S, m N) S {
	for i := range a {
		a[i] -= m
	}
	return a
}

// Mul1 multiplies the values of a by m and returns a.
func Mul1[S ~[]N, N Number](a S, m N) S {
	for i := range a {
		a[i] *= m
	}
	return a
}

// Dot returns the dot product of the input vectors.
func Dot[S ~[]N, N Number](a, b S) N {
	assertMatchingLengths(a, b)
	var sum N
	for i := range a {
		sum += a[i] * b[i]
	}
	return sum
}

// Norm returns the L2 norm of the vector.
func Norm[S ~[]N, N constraints.Float](a S) float64 {
	var norm N
	for _, v := range a {
		norm += v * v
	}
	return math.Sqrt(float64(norm))
}

// Ones returns a slice of n ones. Panics if n is negative.
func Ones[S ~[]N, N Number](n int) S {
	if n < 0 {
		panic(fmt.Sprintf("bad vector length: %d", n))
	}
	a := make(S, n)
	for i := range a {
		a[i] = 1
	}
	return a
}

// Copy returns a copy of the given slice.
func Copy[S ~[]N, N any](a S) S {
	result := make(S, len(a))
	copy(result, a)
	return result
}

// Cast casts the values of a and places them in a new slice.
func Cast[S ~[]N, T ~[]M, N Number, M Number](a S) T {
	t := make(T, len(a))
	for i, s := range a {
		t[i] = M(s)
	}
	return t
}

// Panics if the input vectors are of different lengths.
func assertMatchingLengths[S ~[]N, N any](a, b S) {
	if len(a) != len(b) {
		panic(fmt.Sprintf("mismatching lengths: %d, %d", len(a), len(b)))
	}
}
