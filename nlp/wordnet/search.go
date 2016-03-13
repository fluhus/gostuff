package wordnet

import (
	"math"
)

// Searches for a word in the dictionary. Returns a map from part of speech
// (a, n, r, v) to all synsets that contain that word.
func (wn *Wordnet) Search(word string) map[string][]*Synset {
	result := map[string][]*Synset{}
	for _, pos := range [...]string{"a", "n", "r", "v"} {
		ids := wn.Lemma[pos+"."+word]
		result[pos] = make([]*Synset, len(ids))
		for i, id := range ids {
			result[pos][i] = wn.Synset[id]
		}
	}
	// TODO(amit): Search in exceptions too.
	return result
}

func (wn *Wordnet) PathSimilarity(from, to *Synset) float64 {
	hypFrom := wn.hypernyms(from)
	hypTo := wn.hypernyms(to)
	best := math.MaxInt32
	
	for s := range hypFrom {
		if _, ok := hypTo[s]; ok {
			distance := hypFrom[s] + hypTo[s]
			if distance < best {
				best = distance
			}
		}
	}
	
	if best == math.MaxInt32 { // Found no common ancestor.
		return 0
	}
	
	return 1.0 / (float64(best) + 1.0)
}

// Returns the hypernym hierarchy of the synset, with their distance from the
// synset.
func (wn *Wordnet) hypernyms(ss *Synset) map[*Synset]int {
	result := map[*Synset]int{}
	next := map[*Synset]struct{}{ss:struct{}{}}
	level := 0
	for len(next) > 0 {
		newNext := map[*Synset]struct{}{}
		for s := range next {
			result[s] = level
			for _, ptr := range s.Pointer {
				if ptr.Symbol[:1] == "@" { // Hypernym relation.
					newNext[wn.Synset[ptr.Synset]] = struct{}{}
				}
			}
		}
		level++
		next = newNext
	}
	
	return result
}

