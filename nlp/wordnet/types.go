package wordnet

// An entire wordnet database.
type WordNet struct {
	// Maps from synset ID to synset.
	Synset map[string]*Synset

	// Maps from pos.lemma to synset IDs that contain it.
	Lemma map[string][]string

	// Maps from exceptional word to its forms.
	Exception map[string][]string

	// Maps from example ID to sentence template.
	Example map[int]string
}

// A single synset.
type Synset struct {
	// Part of speech, including 's' for adjective satellite.
	Pos string

	// Words in this synset.
	Word []string

	// Pointers to other synsets.
	Pointer []*Pointer

	// Sentence frames for verbs.
	Frame []*Frame

	// Word definition and usage examples.
	Gloss string

	// Usage examples for words in this synset. Verbs only.
	Example []*Example
}

// A frame is a generic phrase that illustrates how to use a verb.
//
// See the list of frames here:
// https://wordnet.princeton.edu/man/wninput.5WN.html#sect4
type Frame struct {
	// Frame number on the WordNet site.
	FrameNumber int

	// Index of word in the containing synset, -1 for entire synset.
	WordNumber int
}

// Denotes a semantic relation between one synset/word to another.
//
// See list of pointer symbols here:
// https://wordnet.princeton.edu/man/wninput.5WN.html#sect3
type Pointer struct {
	// Relation between the 2 words. Target is <symbol> to source. See
	// package constants for meaning of symbols.
	Symbol string

	// Target synset ID.
	Synset string

	// Index of word in source synset, -1 for entire synset.
	Source int

	// Index of word in target synset, -1 for entire synset.
	Target int
}

// Links a synset word to an example sentence. Applies to verbs only.
type Example struct {
	// Index of word in the containing synset.
	WordNumber int

	// Number of template in the WordNet.Example field.
	TemplateNumber int
}
