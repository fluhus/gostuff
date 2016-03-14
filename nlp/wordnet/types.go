package wordnet

// An entire wordnet database.
//
// Keys in Lemma field are "pos.lemma". For example the key "n.back" relates to
// the noun "back", and the key "v.back" relates to the verb "back".
//
// Keys in Synset field have no human-readable meaning, and should be used
// blindly for matching lemmas to synsets.
//
// Keys in Exception field are plain words.
type WordNet struct {
	Synset    map[string]*Synset  // Maps from synset ID to synset.
	Lemma     map[string][]string // Maps from pos.lemma to synset IDs that contain it.
	Exception map[string][]string // Maps from exceptional word to its forms.
	Example   map[int]string      // Maps from example ID to sentence template.
}

// A single synset.
type Synset struct {
	Pos     string     // Part of speech, including 's' for adjective satellite.
	Word    []string   // Words in this synset.
	Pointer []*Pointer // Pointers to other synsets.
	Frame   []*Frame   // Sentence frames for verbs.
	Gloss   string     // Word definition and usage examples.
}

// A frame is a generic phrase that illustrates how to use a verb.
//
// See the list of frames here:
// https://wordnet.princeton.edu/man/wninput.5WN.html#sect4
type Frame struct {
	FrameNumber int // Frame number on the WordNet site.
	WordNumber  int // Index of word in the containing synset, -1 for entire synset.
}

// Denotes a semantic relation between one synset/word to another.
//
// See list of pointer symbols here:
// https://wordnet.princeton.edu/man/wninput.5WN.html#sect3
type Pointer struct {
	Symbol string // Relation between the 2 words. Target is <symbol> to source.
	Synset string // Target synset ID.
	Source int    // Index of word in source synset, -1 for entire synset.
	Target int    // Index of word in target synset, -1 for entire synset.
}
