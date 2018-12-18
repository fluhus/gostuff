package wordnet

// WordNet is an entire wordnet database.
type WordNet struct {
	// Maps from synset ID to synset.
	Synset map[string]*Synset `json:"synset"`

	// Maps from pos.lemma to synset IDs that contain it.
	Lemma map[string][]string `json:"lemma"`

	// Like Lemma, but synsets are ordered from the most frequently used to the
	// least. Only a subset of the synsets are ranked, so LemmaRanked has less
	// synsets.
	LemmaRanked map[string][]string `json:"lemmaRanked"`

	// Maps from exceptional word to its forms.
	Exception map[string][]string `json:"exception"`

	// Maps from example ID to sentence template. Using string keys for JSON
	// compatibility.
	Example map[string]string `json:"example"`
}

// Synset is a set of synonymous words.
type Synset struct {
	// Synset offset, also used as an identifier.
	Offset string `json:"offset"`

	// Part of speech, including 's' for adjective satellite.
	Pos string `json:"pos"`

	// Words in this synset.
	Word []string `json:"word"`

	// Pointers to other synsets.
	Pointer []*Pointer `json:"pointer"`

	// Sentence frames for verbs.
	Frame []*Frame `json:"frame"`

	// Lexical definition.
	Gloss string `json:"gloss"`

	// Usage examples for words in this synset. Verbs only.
	Example []*Example `json:"example"`
}

// A Frame links a synset word to a generic phrase that illustrates how to use
// it. Applies to verbs only.
//
// See the list of frames here:
// https://wordnet.princeton.edu/man/wninput.5WN.html#sect4
type Frame struct {
	// Index of word in the containing synset, -1 for entire synset.
	WordNumber int `json:"wordNumber"`

	// Frame number on the WordNet site.
	FrameNumber int `json:"frameNumber"`
}

// A Pointer denotes a semantic relation between one synset/word to another.
//
// See list of pointer symbols here:
// https://wordnet.princeton.edu/man/wninput.5WN.html#sect3
type Pointer struct {
	// Relation between the 2 words. Target is <symbol> to source. See
	// package constants for meaning of symbols.
	Symbol string `json:"symbol"`

	// Target synset ID.
	Synset string `json:"synset"`

	// Index of word in source synset, -1 for entire synset.
	Source int `json:"source"`

	// Index of word in target synset, -1 for entire synset.
	Target int `json:"target"`
}

// An Example links a synset word to an example sentence. Applies to verbs only.
type Example struct {
	// Index of word in the containing synset.
	WordNumber int `json:"wordNumber"`

	// Number of template in the WordNet.Example field.
	TemplateNumber int `json:"templateNumber"`
}
