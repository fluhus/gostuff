package clustering

import (
	"fmt"
	"math"

	"github.com/fluhus/gostuff/gnum"
	"github.com/fluhus/gostuff/heaps"
)

// distPyramid is a distance half-matrix.
type distPyramid [][]float64

// dist returns the distance between a and b.
func (d distPyramid) dist(a, b int) float64 {
	if a > b {
		return d[a][b]
	}
	return d[b][a]
}

// makePyramid creates a distance half-matrix.
func makePyramid(n int, f func(int, int) float64) distPyramid {
	nn := n * (n - 1) / 2
	d := make([]float64, 0, nn)
	for i := 1; i < n; i++ {
		for j := 0; j < i; j++ {
			d = append(d, f(j, i))
		}
	}
	result := make([][]float64, n)
	j := 0
	for i := range result {
		result[i] = d[j : j+i]
		j += i
	}
	return result
}

// upgma is an implementation of UPGMA clustering. The distance between clusters
// is the average distance between pairs of their individual elements.
func upgma(n int, f func(int, int) float64) *AggloResult {
	pi := make([]int, n)         // Index of first merge target of each element.
	lambda := make([]float64, n) // Distance of first merge target of each element.

	// Last cluster does not get matched with anyone -> max distance.
	lambda[len(lambda)-1] = math.MaxFloat64

	// Calculate raw distances.
	d := makePyramid(n, f)
	heapss := make([]*heaps.Heap[upgmaCluster], n)
	for i := range heapss {
		heapss[i] = heaps.New(compareUpgmaClusters)
	}
	for i := 1; i < n; i++ {
		for j := 0; j < i; j++ {
			heapss[i].Push(upgmaCluster{j, d[i][j]})
			heapss[j].Push(upgmaCluster{i, d[i][j]})
		}
	}

	// Clustering.
	sizes := gnum.Ones[[]float64](n) // Cluster sizes
	// The identifier of each cluster = highest index of an element
	names := make([]int, n)
	for i := range names {
		names[i] = i
	}
	for i := 0; i < n-1; i++ {
		// Find lowest distance.
		min := math.MaxFloat64
		a, b := -1, -1
		for hi, h := range heapss {
			if h == nil {
				continue
			}
			// Clean up removed clusters.
			if h.Len() == 0 {
				panic(fmt.Sprintf("heap %d with length 0", hi))
			}
			for heapss[h.View()[0].i] == nil {
				h.Pop()
			}
			if h.View()[0].d < min {
				a = hi
				min = h.View()[0].d
				b = h.View()[0].i
			}
		}

		// Create agglo step.
		nmin := gnum.Min2(names[a], names[b])
		nmax := gnum.Max2(names[a], names[b])
		pi[nmin] = nmax
		lambda[nmin] = min

		// Merge clusters.
		names = append(names, nmax)
		sizes = append(sizes, sizes[a]+sizes[b])
		heapss[a] = nil
		heapss[b] = nil
		var cdist []float64
		cheap := heaps.New(compareUpgmaClusters)
		for hi, h := range heapss {
			if h == nil {
				cdist = append(cdist, 0)
				continue
			}
			da := d.dist(a, hi) * sizes[a]
			db := d.dist(b, hi) * sizes[b]
			dd := (da + db) / (sizes[a] + sizes[b])
			cdist = append(cdist, dd)
			h.Push(upgmaCluster{len(sizes) - 1, dd})
			cheap.Push(upgmaCluster{hi, dd})
		}
		d = append(d, cdist)
		heapss = append(heapss, cheap)
	}

	return newAggloResult(pi, lambda)
}

// Cluster info in UPGMA.
type upgmaCluster struct {
	i int     // Cluster index
	d float64 // Distance from cluster i
}

func compareUpgmaClusters(a, b upgmaCluster) bool {
	return a.d < b.d
}
