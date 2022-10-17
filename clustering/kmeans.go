// Package clustering provides basic clustering functions.
package clustering

import (
	"fmt"
	"math"
	"math/rand"

	"github.com/fluhus/gostuff/gnum"
)

// Kmeans performs k-means clustering on the given data. Each vector is an
// element in the clustering. Returns the generated means, and the tag each
// element was given.
func Kmeans(vecs [][]float64, k int) (means [][]float64, tags []int) {
	// K must be at least 1.
	if k < 1 {
		panic(fmt.Sprint("Bad k:", k))
	}

	// Must have at least 1 vector.
	if len(vecs) == 0 {
		panic("Cannot cluster 0 vectors.")
	}

	// If k is too large - that's ok just reduce to avoid out-of-range.
	if k > len(vecs) {
		k = len(vecs)
	}

	// First tagging.
	means = initialMeans(vecs, k)
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

// tag tags each row with the index of its nearest centroid. The old tags are
// used for optimization.
func tag(vecs, means [][]float64, oldTags []int) []int {
	if len(means) == 0 {
		panic("Cannot tag on 0 centroids.")
	}

	// Create a distance matrix of means from one another.
	meansd := make([][]float64, len(means))
	for i := range meansd {
		meansd[i] = make([]float64, len(means))
		for j := range means {
			meansd[i][j] = gnum.L2(means[i], means[j])
		}
	}

	tags := make([]int, len(vecs))

	// Go over vectors.
	for i := range vecs {
		// Find nearest centroid.
		tags[i] = oldTags[i]
		d := gnum.L2(means[oldTags[i]], vecs[i])

		for j := 0; j < len(means); j++ {
			// Use triangle inequality to skip means that are too distant.
			if j == tags[i] || meansd[j][tags[i]] >= 2*d {
				continue
			}

			dj := gnum.L2(means[j], vecs[i])
			if dj < d {
				d = dj
				tags[i] = j
			}
		}
	}

	return tags
}

// findMeans calculates the new means, according to average of tagged rows in
// each group.
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
		gnum.Add(means[tags[i]], vecs[i])
	}

	// Divide by count.
	for i := range means {
		if counts[i] != 0 {
			gnum.Mul1(means[i], 1/float64(counts[i]))
		}
	}

	return means
}

// initialMeans picks the initial means with the K-means++ algorithm.
func initialMeans(vecs [][]float64, k int) [][]float64 {
	result := make([][]float64, k)
	perm := rand.Perm(len(vecs))
	numTrials := 2 + int(math.Log(float64(k)))

	probs := make([]float64, len(vecs))    // Probability of each vector.
	nearest := make([]int, len(vecs))      // Index of nearest mean to each vector.
	distance := make([]float64, len(vecs)) // Distance to nearest mean.
	mdistance := make([][]float64, k)      // Distance between means.
	for i := range mdistance {
		mdistance[i] = make([]float64, k)
	}

	// Pick each mean.
	for i := range result {
		result[i] = make([]float64, len(vecs[0]))

		// First mean is first vector.
		if i == 0 {
			copy(result[0], vecs[perm[0]])
			for _, j := range perm {
				distance[j] = gnum.L2(vecs[j], result[0])
			}
			continue
		}

		// Find next mean.
		bestCandidate := -1
		bestImprovement := -math.MaxFloat64

		for t := 0; t < numTrials; t++ { // Make a few attempts.
			sum := 0.0
			for _, j := range perm {
				probs[j] = distance[j] * distance[j]
				sum += probs[j]
			}
			// Pick element with probability relative to d^2.
			r := rand.Float64() * sum
			newMean := 0
			for r > probs[newMean] {
				r -= probs[newMean]
				newMean++
			}
			copy(result[i], vecs[newMean])

			// Update distances from new mean to other means.
			for j := range mdistance[:i] {
				mdistance[j][i] = gnum.L2(result[i], result[j])
				mdistance[i][j] = mdistance[j][i]
			}

			// Check improvement.
			newImprovement := 0.0
			for j := range vecs {
				if mdistance[i][nearest[j]] < 2*distance[j] {
					d := gnum.L2(vecs[j], result[i])
					d = math.Min(distance[j], d)
					newImprovement += distance[j] - d
				}
			}
			if newImprovement > bestImprovement {
				bestCandidate = newMean
				bestImprovement = newImprovement
			}
		}

		copy(result[i], vecs[bestCandidate])

		// Update distances.
		for j := range mdistance[:i] { // From new mean to other means.
			mdistance[j][i] = gnum.L2(result[i], result[j])
			mdistance[i][j] = mdistance[j][i]
		}
		for j := range vecs { // From vecs to nearest means.
			if mdistance[i][nearest[j]] < 2*distance[j] {
				d := gnum.L2(vecs[j], result[i])
				if d < distance[j] {
					distance[j] = math.Min(distance[j], d)
					nearest[j] = i
				}
			}
		}
	}

	return result
}

// MeanSquaredError calculates the average squared-distance of elements from
// their assigned means.
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
		dist := gnum.L2(means[tags[i]], vecs[i])
		d += dist * dist
	}

	return d / float64(len(vecs))
}
