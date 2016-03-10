package wordnet

// TODO(amit): Complete documentation.

// An entire wordnet database.
type Wordnet struct {
	Synset    map[string]*Synset  // Maps from synset ID to synset.
	Lemma     map[string]*Lemma   // Maps from lemma to lemma object.
	Exception map[string][]string // Maps from exceptional word to its forms.
}

// A single synset.
type Synset struct {
	SsType string       // Part of speech.
	Word   []*DataWord  // Words in this synset.
	Ptr    []*DataPtr   // Pointers to other synsets.
	Frame  []*DataFrame // ???
	Gloss  string       // Word definition and usage examples.
}

// ???
type DataFrame struct {
	FrameNumber int
	WordNumber  int
}

// A word in a synset.
type DataWord struct {
	Word  string // The actual lemma.
	LexId int    // Index that uniquely identifies that sense of word.
}

// A pointer from one synset word to another.
type DataPtr struct {
	Symbol string // Relation between the 2 synsets.
	Synset string // Pointed-to synset.
	Source int    // 1-based index of word in source synset (0 for entire synset).
	Target int    // 1-based index of word in target synset (0 for entire synset).
}

// TODO(amit): Consider removing Lemma.PtrSymbol, as it may be redundant.

// A single lemma (word) in the glossary.
type Lemma struct {
	PtrSymbol []string // Symbols of pointers this lemma has in all synsets containing it.
	Synset    []string // IDs of synsets containing this lemma, as they appear in Wordnet.Synset.
}
