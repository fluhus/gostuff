package wordnet

// TODO(amit): Document exported types.

type Wordnet struct {
	Synset    map[string]*Synset
	Lemma     map[string]*Lemma
	Exception map[string][]string
}

type Synset struct {
	SsType string
	Word   []*DataWord
	Ptr    []*DataPtr
	Frame  []*DataFrame
	Gloss  string
}

type DataFrame struct {
	FrameNumber int
	WordNumber  int
}

type DataWord struct {
	Word  string
	LexId int
}

type DataPtr struct {
	Symbol string
	Synset string
	Source int // 1-based.
	Target int // 1-based.
}

type Lemma struct {
	PtrSymbol []string
	Synset    []string
}
