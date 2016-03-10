package wordnet

// Indexes all words in the data.
func (wn *Wordnet) indexLemma() {
	wn.Lemma = map[string][]string{}
	for id, ss := range wn.Synset {
		pos := id[0:1]
		for _, word := range ss.Word {
			w := pos + "." + word.Word
			wn.Lemma[w] = append(wn.Lemma[w], id)
		}
	}
}
