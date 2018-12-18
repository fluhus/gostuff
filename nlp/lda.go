package nlp

import (
	"fmt"
	"math/rand"
	"sort"
	"time"
)

// LdaVerbose determines whether progress information should be printed during
// LDA. For debugging.
var LdaVerbose = false

// ----- INTERFACE FUNCTIONS ---------------------------------------------------

// Lda performs LDA on the given data. docTokens should contain tokenized
// documents, such that docTokens[i][j] is the j'th token in the i'th document.
// k is the number of topics. Returns the topics and token-topic assignment,
// respective to docTokens.
//
// Topics are returned in a map from word to a probability vector, such that
// the i'th position is the probability of the i'th topic generating that word.
// For each i, the i'th position of all words sum to 1.
func Lda(docTokens [][]string, k int) (map[string][]float64, [][]int) {
	return LdaThreads(docTokens, k, 1)
}

// LdaThreads is like the function Lda but runs on multiple subroutines.
// Calling this function with 1 thread is equivalent to calling Lda.
func LdaThreads(docTokens [][]string, k, numThreads int) (map[string][]float64,
	[][]int) {
	// Check input.
	if k < 1 {
		panic(fmt.Sprintf("k must be positive. Got %d.", k))
	}
	if numThreads < 1 {
		panic(fmt.Sprintf("Number of threads must be positive. Got %d.",
			numThreads))
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
	if LdaVerbose {
		fmt.Println("LDA:", len(words), "words in dictionary")
	}

	// Convert tokens to indexes.
	docs := make([][]int, len(docTokens))
	for i := range docs {
		docs[i] = make([]int, len(docTokens[i]))
		for j := range docs[i] {
			docs[i][j] = words[docTokens[i][j]]
		}
	}

	topics := newDists(k, len(words), 0.1/float64(len(words)))

	// Initial topic assignment.
	doct := make([][]int, len(docs))
	for i := range docs {
		doct[i] = make([]int, len(docs[i]))
		for j := range doct[i] {
			t := rand.Intn(k)
			doct[i][j] = t
			topics[t].add(docs[i][j])
		}
	}

	lastChange := 0 // How many words changed their topic in the last iteration.
	for _, t := range doct {
		lastChange += len(t)
	}
	breakSignals := 0

	// Fun part!
	for {
		newTopics := newDists(k, len(words), 0.1/float64(len(words)))

		// Big buffers for speed.
		push := make(chan int, numThreads*1000)
		pull := make(chan int, numThreads*1000)
		change := make(chan int, numThreads)
		done := make(chan int, numThreads)

		// Pusher thread - pushes documnet index to threads.
		go func() {
			for i := range docs {
				push <- i
			}
			close(push)
		}()

		// Puller thread - updates new topics with done documents.
		go func() {
			count := 0
			progress := -1
			for i := range pull {
				// Print progress if verbose.
				if LdaVerbose {
					count++
					newProgress := count * 100 / len(doct)
					if newProgress > progress {
						progress = newProgress
						fmt.Printf("\rLDA: [%d%%]", progress)
					}
				}

				// Update document.
				for j := range doct[i] {
					newTopics[doct[i][j]].add(docs[i][j])
				}
			}

			if LdaVerbose {
				fmt.Println()
			}
			done <- 0
		}()

		// changeCount thread - counts how many word changed their topic.
		changeCount := 0
		go func() {
			for count := range change {
				changeCount += count
			}
			done <- 0
		}()

		// Worker threads.
		for thread := 0; thread < numThreads; thread++ {
			go func() {
				// Make a local copy of topics.
				myTopics := copyDists(topics)
				myChangeCount := 0
				myRand := newRand()      // Thread-local random to prevent waiting on rand's default source.
				ts := make([]float64, k) // Reusable slice for randomly picking topics.

				// For each document.
				for i := range push {
					// Create distribution of profiles.
					d := newDist(k, 0.1/float64(k))
					for j := range doct[i] {
						d.add(doct[i][j])
					}

					// Reassign each word.
					for j := range doct[i] {
						t := doct[i][j]
						word := docs[i][j]

						// Unassign.
						d.sub(t)
						myTopics[t].sub(word)

						// Pick new topic.
						for k := range ts {
							ts[k] = d.p(k) * myTopics[k].p(word)
						}
						t2 := pickRandom(ts, myRand)
						if t2 != doct[i][j] {
							myChangeCount++
						}

						// Assign.
						doct[i][j] = t2
						d.add(t2)
						myTopics[t2].add(word)
					}

					// Report this doc is done.
					pull <- i
				}

				change <- myChangeCount
				done <- 0
			}()
		}

		// Wait for threads.
		for i := 0; i < numThreads; i++ {
			<-done
		}
		close(pull)
		close(change)
		<-done
		<-done

		// Update topics.
		topics = newTopics

		// Check halting condition.
		if changeCount >= lastChange {
			breakSignals++
			if breakSignals == 5 {
				break
			}
		}

		if LdaVerbose {
			fmt.Printf("LDA: Changes: %d (%d) %.3f\n", changeCount, breakSignals,
				float64(changeCount)/float64(lastChange))
		}
		lastChange = changeCount
	}

	// Make return values.
	topicDists := make([][]float64, len(topics))
	for i := range topicDists {
		topicDists[i] = topics[i].dist()
	}

	dict := map[string][]float64{}
	for word, i := range words {
		d := make([]float64, k)
		for j := range d {
			d[j] = topicDists[j][i]
		}
		dict[word] = d
	}

	return dict, doct
}

// ----- HELPERS ---------------------------------------------------------------

// dist is a distribution on elements by counts.
type dist struct {
	sum    float64
	count  []float64
	alpha  float64
	alphas float64
}

// newDist creates a new empty distribution.
func newDist(n int, alpha float64) *dist {
	return &dist{0, make([]float64, n), alpha, alpha * float64(n)}
}

// newDists creates a slice of empty distributions.
func newDists(k, n int, alpha float64) []*dist {
	result := make([]*dist, k)
	for i := range result {
		result[i] = newDist(n, alpha)
	}
	return result
}

// p returns the probability of i, considering alpha.
func (d *dist) p(i int) float64 {
	if d.sum == 0 {
		return 0
	}
	return (d.count[i] + d.alpha*d.sum) / (d.sum + d.alphas*d.sum)
}

// add increments i by 1.
func (d *dist) add(i int) {
	d.count[i]++
	d.sum++
}

// sun decrements i by 1.
func (d *dist) sub(i int) {
	d.count[i]--
	d.sum--

	if d.count[i] < 0 {
		panic(fmt.Sprintf("Reached negative count for i=%d.", i))
	}
}

// dist returns the counts of this distribution, normalized by its sum.
func (d *dist) dist() []float64 {
	result := make([]float64, len(d.count))
	copy(result, d.count)
	if d.sum != 0 {
		for i := range result {
			result[i] /= d.sum
		}
	}
	return result
}

// copy deep-copies a distribution.
func (d *dist) copy() *dist {
	count := make([]float64, len(d.count))
	for i := range count {
		count[i] = d.count[i]
	}
	return &dist{d.sum, count, d.alpha, d.alphas}
}

// copyDists deep-copies a slice of distributions.
func copyDists(dists []*dist) []*dist {
	result := make([]*dist, len(dists))
	for i := range result {
		result[i] = dists[i].copy()
	}
	return result
}

// top returns the n most likely items in the distribution.
func (d *dist) top(n int) []int {
	s := newDistSorter(d)
	sort.Sort(s)
	if n > len(s.perm) {
		n = len(s.perm)
	}
	return s.perm[:n]
}

// distSorter is a distribution sorting interface.
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

// newRand creates a new random generator.
func newRand() *rand.Rand {
	return rand.New(rand.NewSource(time.Now().UnixNano()))
}

// pickRandom picks a random index from a, with a probability proportional to
// its value. Using a local random-generator to prevent waiting on rand's
// default source.
func pickRandom(a []float64, rnd *rand.Rand) int {
	if len(a) == 0 {
		panic("Cannot pick element from an empty distribution.")
	}

	sum := float64(0)
	for i := range a {
		if a[i] < 0 {
			panic(fmt.Sprintf("Got negative value in distribution: %v", a[i]))
		}
		sum += a[i]
	}
	if sum == 0 {
		return rnd.Intn(len(a))
	}

	r := rnd.Float64() * sum
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
