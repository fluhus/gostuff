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
	SsType  string        // Part of speech, including 's' for adjective satellite.
	Word    []*SynsetWord // Words in this synset.
	Pointer []*Pointer    // Pointers to other synsets.
	Frame   []*Frame      // Sentence frames for verbs.
	Gloss   string        // Word definition and usage examples.
}

// A frame is a generic phrase that illustrates how to use a verb.
//
// See the list of frames here:
// https://wordnet.princeton.edu/man/wninput.5WN.html#sect4
type Frame struct {
	FrameNumber int // Frame number on the WordNet site.
	WordNumber  int // 1-based index of word in the containing synset, 0 for entire synset.
}

// A word in a synset.
type SynsetWord struct {
	Word  string // The actual lemma.
	LexId int    // Index that uniquely identifies that sense of word.
}

// A pointer from one synset word to another.
//
// See list of pointer symbols here:
// https://wordnet.princeton.edu/man/wninput.5WN.html#sect3
type Pointer struct {
	Symbol string // Relation between the 2 words. Target is <symbol> to source.
	Synset string // Target synset.
	Source int    // 1-based index of word in source synset, 0 for entire synset.
	Target int    // 1-based index of word in target synset, 0 for entire synset.
}
