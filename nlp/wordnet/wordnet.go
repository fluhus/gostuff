// Wordnet parser and interface.
//
// !!! UNDER CONSTRUCTION !!!
//
// Basic usage
//
// The main entry point is the Wordnet type. It holds all the data of a
// wordnet dictionary, and provides search methods.
//
// To search for the noun meanings of 'cat':
//  wn, _ := wordnet.Parse(...)
//  catNouns := wn.Search("cat")["n"]
//  // = slice of all synsets that contain the word "cat" and are nouns.
//
// To calculate similarity between words:
//  wn, _ := wordnet.Parse(...)
//  cat := wn.Search("cat")["n"][0]
//  dog := wn.Search("dog")["n"][0]
//  similarity := wn.PathSimilarity(cat, dog, false)
//  // = 0.2
//
// Parts of speech
//
// Some data refers to parts of speech (POS). Everywhere a part of speech is
// expected, it is a single letter as follows:
//  a: adjective
//  n: noun
//  r: adverb
//  v: verb
package wordnet

import (
	"math"
)

// Parses an entire wordnet directory. Path is the root of the directory.
// The parser will trverse it and parse the required files, assuming
// directory structure is as published.
func Parse(path string) (*Wordnet, error) {
	result := &Wordnet{}
	var err error

	result.Synset, err = parseDataFiles(path)
	if err != nil {
		return nil, err
	}

	result.indexLemma()

	result.Exception, err = parseExceptionFiles(path)
	if err != nil {
		return nil, err
	}

	return result, nil
}

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
	// TODO(amit): Search in exceptions too?
	return result
}

// Returns a score denoting how similar two word senses are, based on the
// shortest path that connects the senses in the is-a (hypernym/hypnoym)
// taxonomy. The score is in the range 0 to 1, where 1 means identity and 0
// means completely disjoint.
//
// If simulateRoot is true, will create a common fake root for the top of each
// synset's hierarchy if no common ancestor was found.
//
// Should be equivalent to NLTK's path_similarity function.
func (wn *Wordnet) PathSimilarity(from, to *Synset, simulateRoot bool) float64 {
	hypFrom := wn.hypernyms(from)
	hypTo := wn.hypernyms(to)
	shortest := math.MaxInt32

	// Find common ancestor that gives the shortest path.
	for s := range hypFrom {
		if _, ok := hypTo[s]; ok {
			distance := hypFrom[s] + hypTo[s]
			if distance < shortest {
				shortest = distance
			}
		}
	}

	// If no common ancestor, make a fake root.
	if shortest == math.MaxInt32 {
		if simulateRoot {
			depthFrom := 0
			depthTo := 0
			for _, d := range hypFrom {
				if d > depthFrom {
					depthFrom = d
				}
			}
			for _, d := range hypTo {
				if d > depthTo {
					depthTo = d
				}
			}
			shortest = depthFrom + depthTo + 2 // 2 for fake root.
		} else {
			return 0
		}
	}

	return 1.0 / (float64(shortest) + 1.0)
}

// Returns the hypernym hierarchy of the synset, with their distance from the
// input synset.
func (wn *Wordnet) hypernyms(ss *Synset) map[*Synset]int {
	result := map[*Synset]int{}
	next := map[*Synset]struct{}{ss: struct{}{}}
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
