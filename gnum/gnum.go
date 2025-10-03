// Package gnum provides generic numerical functions.
package gnum

import (
	"fmt"
	"math"

	"golang.org/x/exp/constraints"
)

// Number is a constraint that contains comparable numbers.
type Number interface {
	constraints.Float | constraints.Integer
}

// Max returns the maximal value in the slice or zero if the slice is empty.
func Max[S ~[]N, N constraints.Ordered](s S) N {
	if len(s) == 0 {
		var zero N
		return zero
	}
	e := s[0]
	for _, v := range s[1:] {
		e = max(e, v)
	}
	return e
}

// Min returns the maximal value in the slice or zero if the slice is empty.
func Min[S ~[]N, N constraints.Ordered](s S) N {
	if len(s) == 0 {
		var zero N
		return zero
	}
	e := s[0]
	for _, v := range s[1:] {
		e = min(e, v)
	}
	return e
}

// ArgMax returns the index of the maximal value in the slice or -1 if the slice is empty.
func ArgMax[S ~[]E, E constraints.Ordered](s S) int {
	if len(s) == 0 {
		return -1
	}
	imax, max := 0, s[0]
	for i, v := range s {
		if v > max {
			imax, max = i, v
		}
	}
	return imax
}

// ArgMin returns the index of the minimal value in the slice or -1 if the slice is empty.
func ArgMin[S ~[]E, E constraints.Ordered](s S) int {
	if len(s) == 0 {
		return -1
	}
	imin, min := 0, s[0]
	for i, v := range s {
		if v < min {
			imin, min = i, v
		}
	}
	return imin
}

// Abs returns the absolute value of n.
//
// For floats use [math.Abs].
func Abs[N constraints.Signed](n N) N {
	if n < 0 {
		return -n
	}
	return n
}

// Diff returns the non-negative difference between a and b.
func Diff[N Number](a, b N) N {
	if a > b {
		return a - b
	}
	return b - a
}

// Sum returns the sum of the slice.
func Sum[S ~[]N, N Number](a S) N {
	var sum N
	for _, v := range a {
		sum += v
	}
	return sum
}

// Mean returns the average of the slice.
func Mean[S ~[]N, N Number](a S) float64 {
	return float64(Sum(a)) / float64(len(a))
}

// ExpMean returns the exponential average of the slice.
// Non-positive values result in NaN.
func ExpMean[S ~[]N, N Number](a S) float64 {
	if len(a) == 0 {
		return math.NaN()
	}
	sum := 0.0
	for _, v := range a {
		sum += math.Log(float64(v))
	}
	return math.Exp(sum / float64(len(a)))
}

// Cov returns the covariance of a and b.
func Cov[S ~[]N, N Number](a, b S) float64 {
	assertMatchingLengths(a, b)
	ma := Mean(a)
	mb := Mean(b)
	cov := 0.0
	for i := range a {
		cov += (float64(a[i]) - ma) * (float64(b[i]) - mb)
	}
	cov /= float64(len(a))
	return cov
}

// Var returns the variance of a.
func Var[S ~[]N, N Number](a S) float64 {
	return Cov(a, a)
}

// Std returns the standard deviation of a.
func Std[S ~[]N, N Number](a S) float64 {
	return math.Sqrt(Var(a))
}

// Corr returns the Pearson correlation between the a and b.
func Corr[S ~[]N, N Number](a, b S) float64 {
	return Cov(a, b) / Std(a) / Std(b)
}

// Entropy returns the Shannon-entropy of a.
// The elements in a don't have to sum up to 1.
func Entropy[S ~[]N, N Number](a S) float64 {
	sum := float64(Sum(a))
	result := 0.0
	for i, v := range a {
		if v < 0 {
			panic(fmt.Sprintf("negative value at position %d: %v",
				i, v))
		}
		if v == 0 {
			continue
		}
		p := float64(v) / sum
		result -= p * math.Log2(p)
	}
	return result
}

// Idiv divides a by b, rounded to the nearest integer.
func Idiv[T constraints.Integer](a, b T) T {
	return T(math.Round(float64(a) / float64(b)))
}

// Quantiles returns the elements that divide the given slice
// at the given ratios.
//
// For example, 0.5 returns the middle element,
// 0.25 returns the element at a quarter of the length, etc.
// 0 and 1 return the first and last element, respectively.
func Quantiles[T any](a []T, qq ...float64) []T {
	if len(qq) == 0 {
		return nil
	}
	if len(a) == 0 {
		panic("input slice cannot be empty")
	}
	result := make([]T, 0, len(qq))
	n := float64(len(a) - 1)
	for _, q := range qq {
		i := int(math.Round(q * n))
		result = append(result, a[i])
	}
	return result
}

// NQuantiles returns the elements that divide
// the given slice into n equal parts (up to rounding),
// including the first and last elements.
//
// For example, for n=2 it returns the first element,
// the middle, and the last element.
func NQuantiles[T any](a []T, n int) []T {
	q := make([]float64, n+1)
	for i := range q {
		q[i] = float64(i) / float64(n)
	}
	return Quantiles(a, q...)
}

// LogFactorial returns an approximation of log(n!),
// calculated in constant time.
func LogFactorial(n int) float64 {
	if n < 0 {
		panic(fmt.Sprintf("n cannot be negative: %v", n))
	}
	if n == 0 || n == 1 {
		return 0
	}
	// Stirling's approximation.
	const halfLog2pi = 0x1.d67f1c864beb4p-01 // 0.5*math.Log(2*math.Pi)
	nf := float64(n)
	logn := math.Log(nf)
	return halfLog2pi + 0.5*logn + nf*(logn-1)
}
