package wordnet

// TODO(amit): Complete documentation.

// An entire wordnet database.
type Wordnet struct {
	Synset    map[string]*Synset  // Maps from synset ID to synset.
	Lemma     map[string][]string // Maps from pos.lemma to synset IDs that contain it.
	Exception map[string][]string // Maps from exceptional word to its forms.
}

// A single synset.
type Synset struct {
	SsType string        // Part of speech, including 's' for adjective satellite.
	Word   []*SynsetWord // Words in this synset.
	Ptr    []*DataPtr    // Pointers to other synsets.
	Frame  []*Frame      // ???
	Gloss  string        // Word definition and usage examples.
}

// ???
type Frame struct {
	FrameNumber int
	WordNumber  int
}

// A word in a synset.
type SynsetWord struct {
	Word  string // The actual lemma.
	LexId int    // Index that uniquely identifies that sense of word.
}

// A pointer from one synset word to another.
type DataPtr struct {
	Symbol string // Relation between the 2 words.
	Synset string // Target synset.
	Source int    // 1-based index of word in source synset (0 for entire synset).
	Target int    // 1-based index of word in target synset (0 for entire synset).
}
