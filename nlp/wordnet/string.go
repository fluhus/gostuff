package wordnet

import (
	"bytes"
	"fmt"
)

// String functions for types.

func (s *Synset) String() string {
	result := bytes.NewBuffer(make([]byte, 0, 100))
	fmt.Fprintf(result, "(%s)", s.SsType)
	for i, word := range s.Word {
		if i > 0 {
			fmt.Fprintf(result, ",")
		}
		fmt.Fprintf(result, " %v", word)
	}
	fmt.Fprintf(result, ": %s", s.Gloss)
	return result.String()
}

func (w SynsetWord) String() string {
	return fmt.Sprintf("%s.%d", w.Word, w.LexId)
}
