// Provides statistic functions.
package stats

import (
	"fmt"
	"math"
)

// Returns the sum of the values in a sample.
func Sum(a []float64) float64 {
	sum := 0.0
	for _, v := range a {
		sum += v
	}
	return sum
}

// Returns the mean value of the sample.
func Mean(a []float64) float64 {
	return Sum(a) / float64(len(a))
}

// Returns the covariance of the 2 samples. Panics if lengths don't match.
func Cov(a, b []float64) float64 {
	assertMatchingLengths(a, b)

	ma := Mean(a)
	mb := Mean(b)
	cov := 0.0
	for i := range a {
		cov += (a[i] - ma) * (b[i] - mb)
	}
	cov /= float64(len(a))

	return cov
}

// Returns the variance of the sample.
func Var(a []float64) float64 {
	return Cov(a, a)
}

// Returns the standard deviation of the sample.
func Std(a []float64) float64 {
	return math.Sqrt(Var(a))
}

// Returns the correlation between the samples.
func Corr(a, b []float64) float64 {
	return Cov(a, b) / Std(a) / Std(b)
}

// Returns the minimal element in the sample.
func Min(a []float64) float64 {
	if len(a) == 0 {
		return math.NaN()
	}

	min := a[0]
	for _, v := range a {
		if v < min {
			min = v
		}
	}

	return min
}

// Returns the maximal element in the sample.
func Max(a []float64) float64 {
	if len(a) == 0 {
		return math.NaN()
	}

	max := a[0]
	for _, v := range a {
		if v > max {
			max = v
		}
	}

	return max
}

// Returns the span of the sample (max - min).
func Span(a []float64) float64 {
	return Max(a) - Min(a)
}

// Returns the entropy for the given distribution.
// The distribution does not have to sum up to 1, for it will be normalized
// anyway.
func Ent(distribution []float64) float64 {
	// Sum of the distribution.
	sum := Sum(distribution)

	// Go over each bucket.
	result := 0.0
	for _, v := range distribution {
		// Negative values are not allowed.
		if v < 0.0 {
			return math.NaN()
		}

		// Ignore zeros.
		if v == 0.0 {
			continue
		}

		// Probability.
		p := v / sum

		// Entropy.
		result -= p * math.Log2(p)
	}

	return result
}

// Panics if 2 vectors are of inequal lengths.
func assertMatchingLengths(a, b []float64) {
	if len(a) != len(b) {
		panic(fmt.Sprintf("Mismatching lengths: %d, %d.", len(a), len(b)))
	}
}
