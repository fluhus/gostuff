package wordnet

import (
	"sort"
)

// Indexes all words in the data.
func (wn *Wordnet) indexLemma() {
	wn.Lemma = map[string][]string{}

	// Sort synsets to keep index stable.
	ids := make([]string, 0, len(wn.Synset))
	for id := range wn.Synset {
		ids = append(ids, id)
	}
	sort.Strings(ids)

	for _, id := range ids {
		ss := wn.Synset[id]
		pos := id[0:1]
		for _, word := range ss.Word {
			w := pos + "." + word
			wn.Lemma[w] = append(wn.Lemma[w], id)
		}
	}
}
