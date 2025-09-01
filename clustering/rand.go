package clustering

import (
	"fmt"
	"slices"

	"github.com/fluhus/gostuff/sets"
)

// AdjustedRandIndex compares 2 taggings of the data for similarity. A score of
// 1 means identical, a score of 0 means as good as random, and a negative
// score means worse than random.
func AdjustedRandIndex(tags1, tags2 []int) float64 {
	// Check input.
	if len(tags1) != len(tags2) {
		panic(fmt.Sprintf("Mismatching lengths: %d, %d",
			len(tags1), len(tags2)))
	}

	sets1 := tagsToSets(tags1)
	sets2 := tagsToSets(tags2)

	r := randIndex(sets1, sets2)
	e := expectedRandIndex(sets1, sets2)
	m := maxRandIndex(sets1, sets2)
	return (r - e) / (m - e)
}

// randIndex returns the RI part of the adjusted index.
func randIndex(tags1, tags2 [][]int) float64 {
	r := 0
	for _, t1 := range tags1 {
		for _, t2 := range tags2 {
			r += choose2(sets.SortedIntersectionLen(t1, t2))
		}
	}
	return float64(r)
}

// expectedRandIndex returns the expected index according to hypergeometrical
// distribution.
func expectedRandIndex(tags1, tags2 [][]int) float64 {
	p1 := 0
	n := 0
	for _, tags := range tags1 {
		n += len(tags)
		p1 += choose2(len(tags))
	}
	p2 := 0
	for _, tags := range tags2 {
		p2 += choose2(len(tags))
	}
	p := float64(choose2(n))
	return float64(p1) * float64(p2) / p
}

// maxRandIndex returns the maximal possible index.
func maxRandIndex(tags1, tags2 [][]int) float64 {
	p := 0
	for _, tags := range tags1 {
		p += choose2(len(tags))
	}
	for _, tags := range tags2 {
		p += choose2(len(tags))
	}
	return float64(p) / 2
}

func choose2(n int) int {
	return n * (n - 1) / 2
}

// tagsToSets converts a list of tags to a list of sets of indexes, one list
// for each tag.
func tagsToSets(tags []int) [][]int {
	// Make map from tag to its set.
	sets := map[int][]int{}
	for i, tag := range tags {
		sets[tag] = append(sets[tag], i)
	}

	// Convert map to slice.
	result := make([][]int, 0, len(sets))
	for _, set := range sets {
		slices.Sort(set)
		result = append(result, set)
	}

	return result
}
