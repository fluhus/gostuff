package clustering

import (
	"fmt"
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
func randIndex(tags1, tags2 []intSet) float64 {
	r := 0
	for _, t1 := range tags1 {
		for _, t2 := range tags2 {
			r += choose2(t1.intersect(t2))
		}
	}
	return float64(r)
}

// expectedRandIndex returns the expected index according to hypergeometrical
// distribution.
func expectedRandIndex(tags1, tags2 []intSet) float64 {
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
func maxRandIndex(tags1, tags2 []intSet) float64 {
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

// ----- INT SET --------------------------------------------------------------

// intSet is a set of integers.
type intSet map[int]struct{}

// tagsToSets converts a list of tags to a list of sets of indexes, one list
// for each tag.
func tagsToSets(tags []int) []intSet {
	// Make map from tag to its set.
	sets := map[int]intSet{}
	for i, tag := range tags {
		if sets[tag] == nil {
			sets[tag] = intSet{}
		}
		sets[tag].add(i)
	}

	// Convert map to slice.
	result := make([]intSet, 0, len(sets))
	for _, set := range sets {
		result = append(result, set)
	}

	return result
}

// add adds a number to the set.
func (is intSet) add(i int) {
	is[i] = struct{}{}
}

// contains checks if a set contains the given element.
func (is intSet) contains(i int) bool {
	_, ok := is[i]
	return ok
}

// intersect returns the size of the intersection of the 2 sets.
func (is intSet) intersect(other intSet) int {
	result := 0
	for i := range is {
		if other.contains(i) {
			result++
		}
	}
	return result
}
