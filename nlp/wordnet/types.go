package wordnet

// TODO(amit): Document exported types.

type Wordnet struct {
	Data map[string]*Synset
}

type DataFrame struct {
	FrameNumber int
	WordNumber  int
}

type DataWord struct {
	Word  string
	LexId int
}

type Synset struct {
	SsType string
	Word   []*DataWord
	Ptr    []*DataPtr
	Frame  []*DataFrame
	Gloss  string
}

type DataPtr struct {
	Symbol string
	Synset string
	Source int // 1-based.
	Target int // 1-based.
}


