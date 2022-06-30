// Package stats provides basic statistic functions.
//
// Deprecated: use gnum instead.
package stats

import (
	"fmt"
	"math"
	"sort"
)

// Sum returns the sum of values in a sample.
func Sum(a []float64) float64 {
	sum := 0.0
	for _, v := range a {
		sum += v
	}
	return sum
}

// Mean returns the mean value of the sample.
func Mean(a []float64) float64 {
	return Sum(a) / float64(len(a))
}

// Cov returns the covariance of the 2 samples. Panics if lengths don't match.
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

// Var returns the variance of the sample.
func Var(a []float64) float64 {
	return Cov(a, a)
}

// Std returns the standard deviation of the sample.
func Std(a []float64) float64 {
	return math.Sqrt(Var(a))
}

// Corr returns the correlation between the samples. Panics if lengths don't
// match.
func Corr(a, b []float64) float64 {
	return Cov(a, b) / Std(a) / Std(b)
}

// Min returns the minimal element in the sample.
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

// Max returns the maximal element in the sample.
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

// Span returns the span of the sample (max - min).
func Span(a []float64) float64 {
	return Max(a) - Min(a)
}

// Ent returns the entropy for the given distribution.
// The distribution does not have to sum up to 1.
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

// Hist creates a histogram of counts of values in the given slice. 'counts'
// maps a unique value from 'a' to its count. 'byValue' holds uniqe values
// sorted. 'byCount' holds unique values sorted by their counts, most common
// first.
func Hist(a []float64) (counts map[float64]int, byValue []float64,
	byCount []float64) {
	// Create raw counts.
	counts = map[float64]int{}
	for _, value := range a {
		counts[value]++
	}

	// Take unique values.
	byValue = make([]float64, 0, len(counts))
	for value := range counts {
		byValue = append(byValue, value)
	}
	sort.Sort(sort.Float64Slice(byValue))

	// Sort by counts.
	byCount = make([]float64, len(byValue))
	copy(byCount, byValue)
	sort.Sort(&histSorter{counts, byCount})

	return
}

// Type for histogram sorting.
type histSorter struct {
	counts map[float64]int
	values []float64
}

func (h *histSorter) Len() int {
	return len(h.values)
}
func (h *histSorter) Less(i, j int) bool {
	return h.counts[h.values[i]] > h.counts[h.values[j]]
}
func (h *histSorter) Swap(i, j int) {
	h.values[i], h.values[j] = h.values[j], h.values[i]
}

// Panics if 2 vectors are of inequal lengths.
func assertMatchingLengths(a, b []float64) {
	if len(a) != len(b) {
		panic(fmt.Sprintf("Mismatching lengths: %d, %d.", len(a), len(b)))
	}
}
