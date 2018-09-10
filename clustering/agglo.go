package clustering

import (
	"fmt"
	"math"
	"sort"
)

// How agglomerative clustering should calculate distance between clusters.
const (
	AggloMin = iota // Minimal distance between any pair of elements.
	AggloMax        // Maximal distance between any pair of elements.
)

// Agglo performs agglomerative clustering on the indexes 0 to n-1. d should
// return the distance between the i'th and j'th element, such that
// d(i,j)=d(j,i) and d(i,i)=0.
//
// clusterDist should be one of AggloMin or AggloMax.
//
// Works in O(n^2) time and makes O(n^2) calls to d.
func Agglo(n int, clusterDist int, d func(int, int) float64) *AggloResult {
	if n <= 0 {
		panic(fmt.Sprintf("Bad n: %d, must be positive", n))
	}

	switch clusterDist {
	case AggloMin:
		return slink(n, d)
	case AggloMax:
		return clink(n, d)
	default:
		panic(fmt.Sprintf("Unsupported cluster distance: %v, "+
			"want AggloMin or AggloMax", clusterDist))
	}
}

// slink is an implementation of the SLINK algorithm.
//
// Copied from:
// https://www.cs.ucsb.edu/~veronika/MAE/SLINK_sibson.pdf
func slink(n int, d func(int, int) float64) *AggloResult {
	// Implementation copied from paper, pardon the crap names.
	m := make([]float64, n)      // Distance of i'th element from elements/clusters.
	pi := make([]int, n)         // Index of first merge target of each element.
	lambda := make([]float64, n) // Distance of first merge target of each element.

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

// clink is an implementation of the CLINK algorithm.
//
// Copied from:
// https://academic.oup.com/comjnl/article-pdf/20/4/364/1108735/200364.pdf
func clink(n int, d func(int, int) float64) *AggloResult {
	// Implementation copied from paper, pardon the crap names.
	m := make([]float64, n)      // Distance of i'th element from elements/clusters.
	pi := make([]int, n)         // Index of first merge target of each element.
	lambda := make([]float64, n) // Distance of first merge target of each element.

	lambda[0] = math.MaxFloat64

	for i := 1; i < n; i++ {
		pi[i] = i
		lambda[i] = math.MaxFloat64

		for j := 0; j < i; j++ {
			m[j] = d(i, j)
		}

		for j := 0; j < i; j++ {
			if lambda[j] < m[j] {
				m[pi[j]] = math.Max(m[pi[j]], m[j])
				m[j] = math.MaxFloat64
			}
		}

		a := i - 1
		for j := 0; j < i; j++ {
			if lambda[i-j-1] >= m[pi[i-j-1]] {
				if m[i-j-1] < m[a] {
					a = i - j - 1
				}
			} else {
				m[i-j-1] = math.MaxFloat64
			}
		}

		b := pi[a]
		c := lambda[a]
		pi[a] = i
		lambda[a] = m[a]
		for a < i-1 {
			if b < i-1 {
				d := pi[b]
				e := lambda[b]
				pi[b] = i
				lambda[b] = c
				b = d
				c = e
			} else if b == i-1 {
				pi[b] = i
				lambda[b] = c
				break
			}
		}

		for j := 0; j < i; j++ {
			if pi[pi[j]] == i && lambda[j] >= lambda[pi[j]] {
				pi[j] = i
			}
		}
	}

	return newAggloResult(pi, lambda)
}

// AggloResult is an interactive agglomerative-clustering result.
type AggloResult struct {
	pi     []int
	lambda []float64
	perm   []int
	dict   []string
}

// Dict returns the string representations of elements in the clustering.
func (r *AggloResult) Dict() []string {
	return a.dict
}

// newAggloResult creates a new result.
func newAggloResult(pi []int, lambda []float64) *AggloResult {
	result := &AggloResult{pi, lambda, make([]int, len(pi)), nil}
	for i := range result.perm {
		result.perm[i] = i
	}
	sort.Sort((*aggloSorter)(result))
	return result
}

// aggloSorter is a sorting interface for AggloResult, for sorting by distance.
// This actually sorts the agglomerative steps by their order of occurance.
type aggloSorter AggloResult

// Len returns the number of elements in the sorter.
func (r *aggloSorter) Len() int {
	return len(r.perm)
}

// Less compares two steps by their order of occurance.
func (r *aggloSorter) Less(i, j int) bool {
	return r.lambda[r.perm[i]] < r.lambda[r.perm[j]]
}

// Swap swaps two steps.
func (r *aggloSorter) Swap(i, j int) {
	r.perm[i], r.perm[j] = r.perm[j], r.perm[i]
}

// SetDict sets the string representation of each element, for the String()
// function. Returns itself for chaining.
func (r *AggloResult) SetDict(dict []string) *AggloResult {
	if len(dict) != len(r.perm) {
		panic(fmt.Sprintf("Bad dictionary size: %d, expected %d",
			len(dict), len(r.perm)))
	}
	r.dict = dict
	return r
}

// String returns a representation of the clustering. If SetDict was not
// called, will use element numbers.
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
		strs[j] = fmt.Sprintf("[%s, %s]", strs[i], strs[j])
	}

	// TODO(amit): Panic if reached here.
	return ""
}

// Len returns the number of steps in this clustering. Equals the number of
// elements - 1.
func (r *AggloResult) Len() int {
	return len(r.perm) - 1
}

// An AggloStep is a single step in the clustering process.
// The index of a cluster is the greatest indexed element in it.
// C2 is always greater than C1.
type AggloStep struct {
	C1 int     // Index of the first merged cluster.
	C2 int     // Index of the second merged cluster.
	D  float64 // Distance between the clusters when merging.
}

// Step returns the i'th step in the clustering.
func (r *AggloResult) Step(i int) AggloStep {
	return AggloStep{r.perm[i], r.pi[r.perm[i]], r.lambda[r.perm[i]]}
}
