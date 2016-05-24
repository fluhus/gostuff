// Performs k-means analysis on data.
package clustering

import (
	"fmt"
	"math/rand"

	"github.com/fluhus/gostuff/vectors"
)

// Performs k-means clustering on the given data. Each vector is an element in
// the clustering. Returns the tags each element was given, and the generated
// means.
func Kmeans(vecs [][]float64, k int) (tags []int, means [][]float64) {
	// K must be at least 1.
	if k < 1 {
		panic(fmt.Sprint("Bad k:", k))
	}

	// Must have at least 1 vector.
	if len(vecs) == 0 {
		panic("Cannot cluster 0 vectors.")
	}

	// Use m & n to ease syntax.
	m, n := len(vecs), len(vecs[0])

	// If k is too large - that's ok just reduce to avoid out-of-range.
	if k > m {
		k = m
	}

	// Pick random elements to be the means.
	initMeans := rand.Perm(m)[:k]
	means = make([][]float64, k)
	for i := range means {
		means[i] = make([]float64, n)
		copy(means[i], vecs[initMeans[i]])
	}

	// First tagging.
	tags = tag(vecs, means, make([]int, len(vecs)))
	dist := MeanSquaredError(vecs, means, tags)
	distOld := 2 * dist

	// Iterate until converged.
	for dist > distOld || dist/distOld < 0.999 {
		distOld = dist
		means = findMeans(vecs, tags, k)
		tags = tag(vecs, means, tags)
		dist = MeanSquaredError(vecs, means, tags)
	}

	return
}

// Tags each row with the index of its nearest centroid. The old tags are used
// for optimization.
func tag(vecs, means [][]float64, oldTags []int) []int {
	if len(means) == 0 {
		panic("Cannot tag on 0 centroids.")
	}

	// Create a distance matrix of means from one another.
	meansd := make([][]float64, len(means))
	for i := range meansd {
		meansd[i] = make([]float64, len(means))
		for j := range means {
			meansd[i][j] = vectors.L2(means[i], means[j])
		}
	}

	tags := make([]int, len(vecs))

	// Go over vectors.
	for i := range vecs {
		// Find nearest centroid.
		tags[i] = oldTags[i]
		d := vectors.L2(means[oldTags[i]], vecs[i])

		for j := 0; j < len(means); j++ {
			// Use triangle inequality to skip means that are too distant.
			if j == tags[i] || meansd[j][tags[i]] >= 2*d {
				continue
			}

			dj := vectors.L2(means[j], vecs[i])
			if dj < d {
				d = dj
				tags[i] = j
			}
		}
	}

	return tags
}

// Calculates the new means, according to average of tagged rows in each
// group.
func findMeans(vecs [][]float64, tags []int, k int) [][]float64 {
	// Initialize new arrays.
	means := make([][]float64, k)
	for i := range means {
		means[i] = make([]float64, len(vecs[0]))
	}
	counts := make([]int, k)

	// Sum all vectors according to tags.
	for i := range vecs {
		counts[tags[i]]++
		vectors.Add(means[tags[i]], vecs[i])
	}

	// Divide by count.
	for i := range means {
		if counts[i] != 0 {
			vectors.Mul(means[i], 1/float64(counts[i]))
		}
	}

	return means
}

// Calculates the average squared-distance of elements from their assigned
// means.
func MeanSquaredError(vecs, means [][]float64, tags []int) float64 {
	if len(tags) != len(vecs) {
		panic(fmt.Sprintf("Non-matching lengths of matrix and tags: %d, %d",
			len(vecs), len(tags)))
	}
	if len(vecs) == 0 {
		return 0
	}

	d := 0.0
	for i := range tags {
		dist := vectors.L2(means[tags[i]], vecs[i])
		d += dist * dist
	}

	return d / float64(len(vecs))
}
