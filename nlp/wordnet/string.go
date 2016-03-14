package wordnet

import (
	"bytes"
	"fmt"
)

// String functions for types.

// TODO(amit): Consider removing String() functions to simplify the API.

// Returns a compact string representation of the WordNet data collection, for
// debugging.
func (w *WordNet) String() string {
	return fmt.Sprintf("WordNet[%d lemmas, %d synsets, %d exceptions,"+
		" %d examples]",
		len(w.Lemma), len(w.Synset), len(w.Exception), len(w.Example))
}

// Returns a string representation of the synset, for debugging.
func (s *Synset) String() string {
	result := bytes.NewBuffer(make([]byte, 0, 100))
	fmt.Fprintf(result, "Synset[%s.", s.Pos)
	for i, word := range s.Word {
		if i > 0 {
			fmt.Fprintf(result, ",")
		}
		fmt.Fprintf(result, " %v", word)
	}
	fmt.Fprintf(result, ": %s]", s.Gloss)
	return result.String()
}
