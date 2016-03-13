// Wordnet parser and interface.
//
// !!! UNDER CONSTRUCTION !!!
//
// Basic usage
//
// The main entry point is the Wordnet type. It holds all the data of a
// wordnet dictionary, and provides search methods.
//
// For example, to search for the noun meanings of 'cat':
//  wn, _ := wordnet.Parse(...)
//  catNouns := wn.Search("cat")["n"]
// Will return the synsets that contain the word "cat" and are nouns.
//
// Parts of speech
//
// Some data refers to parts of speech (POS). Everywhere a part of speech is
// expected, it is a single letter as follows:
//  a: adjective
//  n: noun
//  r: adverb
//  v: verb
//
// Keys in Lemma field
//
// Keys are "pos.lemma". For example the key "n.back" relates to the noun
// "back", and the key "v.back" relates to the verb "back".
//
// Keys in Synset field
//
// These have no human-readable meaning, and should be used blindly for
// matching lemmas to synsets.
//
// Keys are "pos.byte_offset", where byte_offset if the field from the original
// data files. Here it has no meaning as real offset, but as a unique ID.
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
// Should be equivalent to NLTK's path_similarity function.
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
// input synset.
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
