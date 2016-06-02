package clustering

import (
	"fmt"
	"math"
	"sort"
)

// Performs agglomerative clustering on the indexes 0 to n-1. d should return
// the distance between the i'th and j'th element, such that d(i,j)=d(j,i) and
// d(i,i)=0.
//
// Works in O(n^2) time and makes O(n^2) calls to d.
func Agglo(n int, d func(int, int) float64) *AggloResult {
	if n <= 0 {
		panic(fmt.Sprintf("Bad n: %d, must be positive", n))
	}

	// An implementation of the SLINK algorithm.
	// Copied from paper, pardon the crap names.
	m := make([]float64, n)
	pi := make([]int, n)
	lambda := make([]float64, n)
	lambda[0] = math.MaxFloat64

	for i := 1; i < n; i++ {
		pi[i] = i
		lambda[i] = math.MaxFloat64

		for j := 0; j < i; j++ {
			m[j] = d(i, j)
		}

		for j := 0; j < i; j++ {
			if m[j] <= lambda[j] {
				m[pi[j]] = math.Min(m[pi[j]], lambda[j])
				lambda[j] = m[j]
				pi[j] = i
			} else {
				m[pi[j]] = math.Min(m[pi[j]], m[j])
			}
		}

		for j := 0; j < i; j++ {
			if lambda[j] >= lambda[pi[j]] {
				pi[j] = i
			}
		}
	}

	return newAggloResult(pi, lambda)
}

// An interactive agglomerative-clustering result.
type AggloResult struct {
	pi     []int
	lambda []float64
	perm   []int
	dict   []string
}

// Creates a new result.
func newAggloResult(pi []int, lambda []float64) *AggloResult {
	result := &AggloResult{pi, lambda, make([]int, len(pi)), nil}
	for i := range result.perm {
		result.perm[i] = i
	}
	sort.Sort((*aggloSorter)(result))
	return result
}

// Sorting interface for AggloResult, for sorting by distance. This actually
// sorts the agglomerative steps by their order of occurance.
type aggloSorter AggloResult

func (r *aggloSorter) Len() int {
	return len(r.perm)
}
func (r *aggloSorter) Less(i, j int) bool {
	return r.lambda[r.perm[i]] < r.lambda[r.perm[j]]
}
func (r *aggloSorter) Swap(i, j int) {
	r.perm[i], r.perm[j] = r.perm[j], r.perm[i]
}

// Sets the string representation of each element, for the String() function.
// Returns itself for chaining.
func (r *AggloResult) SetDict(dict []string) *AggloResult {
	if len(dict) != len(r.perm) {
		panic(fmt.Sprintf("Bad dictionary size: %d, expected %d",
			len(dict), len(r.perm)))
	}
	r.dict = dict
	return r
}

// String representation of the clustering. If SetDict was not called, will use
// element numbers.
func (r *AggloResult) String() string {
	strs := make([]string, len(r.perm))
	for i := range strs {
		if r.dict == nil {
			strs[i] = fmt.Sprint(i)
		} else {
			strs[i] = r.dict[i]
		}
	}
	for _, i := range r.perm {
		j := r.pi[i]
		if i == j { // Reached the end.
			return strs[i]
		}
		strs[j] = fmt.Sprintf("[%s, %s, %.1f]", strs[i], strs[j], r.lambda[i])
	}

	// TODO(amit): Panic if reached here.
	return ""
}

// Returns the number of steps in this clustering. Equals the number of
// elements - 1.
func (r *AggloResult) Len() int {
	return len(r.perm) - 1
}

// Returns the i'th step in the clustering. Returns the indexes of the merged
// clusters and their distance just before merging. The index of a cluster is
// the greatest indexed element in it.
func (r *AggloResult) Step(i int) (int, int, float64) {
	return r.perm[i], r.pi[r.perm[i]], r.lambda[r.perm[i]]
}
