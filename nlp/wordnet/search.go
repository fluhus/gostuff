package wordnet

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
