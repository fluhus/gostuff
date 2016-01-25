package nlp

import (
	"fmt"
	"math/rand"
	"sort"
)

// ----- INTERFACE FUNCTION ----------------------------------------------------

func Lda(docTokens [][]string, k int) ([][]float32, [][]int, []string) {
	// Check input.
	if k < 1 {
		panic(fmt.Sprintf("k must be positive. Got %d.", k))
	}

	// Create word map.
	words := map[string]int{}
	for _, doc := range docTokens {
		for _, word := range doc {
			if _, ok := words[word]; !ok {
				words[word] = len(words)
			}
		}
	}
	if len(words) == 0 {
		panic("Found 0 words in documents.")
	}

	// Convert tokens to indexes.
	docs := make([][]int, len(docTokens))
	for i := range docs {
		docs[i] = make([]int, len(docTokens[i]))
		for j := range docs[i] {
			docs[i][j] = words[docTokens[i][j]]
		}
	}

	topics := newDists(k, len(words), 0.1/float32(len(words)))

	// Initial assignment.
	doct := make([][]int, len(docs))
	for i := range docs {
		doct[i] = make([]int, len(docs[i]))
		for j := range doct[i] {
			t := rand.Intn(k)
			doct[i][j] = t
			topics[t].add(docs[i][j])
		}
	}

	// Fun part.
	ts := make([]float32, k) // Reusable slice for randomly picking topics.
	lastChange := len(words)
	breakSignals := 0
	for iter := 0; iter > -1; iter++ {
		changeMap := map[int]bool{}
		for i := range doct {
			// Create distribution of profiles.
			d := newDist(k, 0.1/float32(k))
			for j := range doct[i] {
				d.add(doct[i][j])
			}

			// Reassign each word.
			for j := range doct[i] {
				t := doct[i][j]
				word := docs[i][j]

				// Unassign.
				d.sub(t)
				topics[t].sub(word)

				// Pick new topic.
				for k := range ts {
					ts[k] = topics[k].p(word) * d.p(k)
				}
				t = pickRandom(ts)
				if t != doct[i][j] {
					changeMap[word] = true
				}

				// Assign.
				doct[i][j] = t
				d.add(t)
				topics[t].add(word)
			}
		}

		if len(changeMap) >= lastChange {
			breakSignals++
			if breakSignals == 5 {
				break
			}
		}
		lastChange = len(changeMap)
	}

	// Make return values.
	sdrow := make([]string, len(words))
	for word, i := range words {
		sdrow[i] = word
	}

	topicDists := make([][]float32, len(topics))
	for i := range topicDists {
		topicDists[i] = topics[i].dist()
	}

	return topicDists, doct, sdrow
}

// ----- HELPERS ---------------------------------------------------------------

// A distribution on elements by counts.
type dist struct {
	sum    float32
	count  []float32
	alpha  float32
	alphas float32
}

// Creates a new empty distribution.
func newDist(n int, alpha float32) *dist {
	return &dist{0, make([]float32, n), alpha, alpha * float32(n)}
}

// Creates a slice of empty distributions.
func newDists(k, n int, alpha float32) []*dist {
	result := make([]*dist, k)
	for i := range result {
		result[i] = newDist(n, alpha)
	}
	return result
}

// Returns the probability of i, considering alpha.
func (d *dist) p(i int) float32 {
	if d.sum == 0 {
		return 0
	}
	return (d.count[i] + d.alpha*d.sum) / (d.sum + d.alphas*d.sum)
}

// Increments i by 1.
func (d *dist) add(i int) {
	d.count[i]++
	d.sum++
}

// Decrements i by 1.
func (d *dist) sub(i int) {
	d.count[i]--
	d.sum--

	if d.count[i] < 0 {
		panic(fmt.Sprintf("Reached negative count for i=%d.", i))
	}
}

// Returns the counts of this distribution, normalized by its sum.
func (d *dist) dist() []float32 {
	result := make([]float32, len(d.count))
	for i := range result {
		result[i] = d.count[i] / d.sum
	}
	return result
}

// Returns the n most likely items in the distribution.
func (d *dist) top(n int) []int {
	s := newDistSorter(d)
	sort.Sort(s)
	if n > len(s.perm) {
		n = len(s.perm)
	}
	return s.perm[:n]
}

// Distribution sorting interface.
type distSorter struct {
	*dist
	perm []int
}

func newDistSorter(d *dist) *distSorter {
	s := &distSorter{d, make([]int, len(d.count))}
	for i := range s.perm {
		s.perm[i] = i
	}
	return s
}

func (d *distSorter) Len() int {
	return len(d.perm)
}

func (d *distSorter) Less(i, j int) bool {
	return d.count[d.perm[i]] > d.count[d.perm[j]]
}

func (d *distSorter) Swap(i, j int) {
	d.perm[i], d.perm[j] = d.perm[j], d.perm[i]
}

// Picks a random index from a, with a probability proportional to its value.
func pickRandom(a []float32) int {
	if len(a) == 0 {
		panic("Cannot pick element from an empty distribution.")
	}

	sum := float32(0)
	for i := range a {
		if a[i] < 0 {
			panic(fmt.Sprintf("Got negative value in distribution: %v", a[i]))
		}
		sum += a[i]
	}
	if sum == 0 {
		return rand.Intn(len(a))
	}

	r := rand.Float32() * sum
	i := 0
	for i < len(a) && r > a[i] {
		r -= a[i]
		i++
	}
	if i == len(a) {
		i--
	}
	return i
}
