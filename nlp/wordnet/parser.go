package wordnet

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// TODO(amit): Convenience functions for pointers?

// ----- FILE LISTS -----------------------------------------------------------

var (
	dataFiles = map[string]string{
		"data.adj":  "a",
		"data.adv":  "r",
		"data.noun": "n",
		"data.verb": "v",
	}
	exceptionFiles = map[string]string{
		"adj.exc":  "a",
		"adv.exc":  "r",
		"noun.exc": "n",
		"verb.exc": "v",
	}
	indexFiles = []string{
		"index.adj",
		"index.adv",
		"index.noun",
		"index.verb",
	}
	exampleFile      = "sents.vrb"
	exampleIndexFile = "sentidx.vrb"
)

// ----- LEMMA INDEX PARSING --------------------------------------------------

// Parses the index files.
func parseIndexFiles(path string) (map[string][]string, error) {
	result := map[string][]string{}

	for _, file := range indexFiles {
		// Read index file.
		f, err := os.Open(filepath.Join(path, file))
		if err != nil {
			return nil, fmt.Errorf("%v: %v", file, err)
		}
		m, err := parseIndex(f)
		if err != nil {
			return nil, fmt.Errorf("%v: %v", file, err)
		}

		// Merge index with result.
		for lemma := range m {
			result[lemma] = m[lemma]
		}
	}

	return result, nil
}

// Parses the contents of an index file.
func parseIndex(r io.Reader) (map[string][]string, error) {
	result := map[string][]string{}
	scanner := bufio.NewScanner(r)

	lineNum := 0
	for scanner.Scan() {
		lineNum++
		if strings.HasPrefix(scanner.Text(), "  ") { // Copyright line.
			continue
		}

		line, err := parseIndexLine(scanner.Text())
		if err != nil {
			return nil, fmt.Errorf("Line %d: %v", lineNum, err)
		}

		if len(line.synset) == 1 {
			line.ranked = 1
		}
		for i := range line.synset {
			line.synset[i] = line.pos + line.synset[i]
		}
		if line.ranked > 0 {
			result[line.pos+"."+line.lemma] = line.synset[:line.ranked]
		}
	}

	return result, nil
}

// A single line in an index file.
type indexLine struct {
	lemma  string
	pos    string
	ptr    []string
	synset []string
	ranked int
}

// Parses an index file line.
func parseIndexLine(line string) (*indexLine, error) {
	result := &indexLine{}
	parts := strings.Split(strings.Trim(line, " "), " ")

	if len(parts) < 7 {
		return nil, fmt.Errorf("bad number of parts: %d, expected at least 7",
			len(parts))
	}

	result.lemma = parts[0]
	result.pos = parts[1]

	synsetCount, err := parseDeciUint(parts[2])
	if err != nil {
		return nil, fmt.Errorf("bad synset count: %s", parts[2])
	}
	ptrCount, err := parseDeciUint(parts[3])
	if err != nil {
		return nil, fmt.Errorf("bad pointer count: %s", parts[3])
	}

	parts = parts[4:]
	if len(parts) < ptrCount+2+synsetCount {
		return nil, fmt.Errorf("bad number of parts: %d, expected %d",
			len(parts)+4, ptrCount+synsetCount+6)
	}

	result.ptr = parts[:ptrCount]
	parts = parts[ptrCount:]

	result.ranked, err = parseDeciUint(parts[1])
	if err != nil {
		return nil, fmt.Errorf("Bad tagsense count: %s", parts[1])
	}

	result.synset = parts[2:]
	if result.ranked > len(result.synset) {
		return nil, fmt.Errorf("Bad tagsense-count: %d is greated than "+
			"synset count %d.", result.ranked, len(result.synset))
	}

	return result, nil
}

// ----- VERB EXAMPLE PARSING -------------------------------------------------

// Parses the verb example file.
func parseExampleFile(path string) (map[string]string, error) {
	f, err := os.Open(filepath.Join(path, exampleFile))
	if err != nil {
		return nil, fmt.Errorf("%s: %v", exampleFile, err)
	}
	return parseExamples(f)
}

// Parses a verb example file.
func parseExamples(r io.Reader) (map[string]string, error) {
	result := map[string]string{}
	scanner := bufio.NewScanner(r)

	lineNum := 0
	for scanner.Scan() {
		lineNum++
		parts := strings.Split(scanner.Text(), " ")
		if len(parts) == 0 {
			return nil, fmt.Errorf("line %d: No data to parse", lineNum)
		}
		_, err := parseDeciUint(parts[0])
		if err != nil {
			return nil, fmt.Errorf("line %d: %v", lineNum, err)
		}
		result[parts[0]] = strings.Join(parts[1:], " ")
	}

	return result, nil
}

// Parses the verb example index file.
func parseExampleIndexFile(path string) (map[string][]int, error) {
	f, err := os.Open(filepath.Join(path, exampleIndexFile))
	if err != nil {
		return nil, fmt.Errorf("%s: %v", exampleIndexFile, err)
	}
	return parseExampleIndex(f)
}

// Parses an entire verb example index file.
func parseExampleIndex(r io.Reader) (map[string][]int, error) {
	result := map[string][]int{}
	scanner := bufio.NewScanner(r)

	lineNum := 0
	for scanner.Scan() {
		lineNum++
		raw, err := parseExampleIndexLine(scanner.Text())
		if err != nil {
			return nil, fmt.Errorf("line %d: %v", lineNum, err)
		}
		key := fmt.Sprintf("%s.%d.%d", raw.lemma, raw.lexFileNum, raw.lexId)
		result[key] = raw.exampleIds
	}

	if scanner.Err() != nil {
		return nil, scanner.Err()
	}

	return result, nil
}

// Represents a single line in the verb example index file.
type rawExampleIndex struct {
	lemma      string
	pos        int
	lexFileNum int
	lexId      int
	headWord   string
	headId     int
	exampleIds []int
}

// Parses a single line in the lemma-example index file.
func parseExampleIndexLine(line string) (*rawExampleIndex, error) {
	result := &rawExampleIndex{}
	parts := strings.Split(line, " ")
	if len(parts) != 2 {
		return nil, fmt.Errorf("bad number of parts: %d, expected 2",
			len(parts))
	}

	// Parse sense.
	senseParts := strings.Split(parts[0], "%")
	if len(senseParts) != 2 {
		return nil, fmt.Errorf("bad number of sense-key parts: %d, expected"+
			" 2", len(senseParts))
	}

	result.lemma = senseParts[0]
	lexSenseParts := strings.Split(senseParts[1], ":")
	if len(lexSenseParts) != 5 {
		return nil, fmt.Errorf("bad number of lex-sense parts: %d, expected"+
			" 5", len(lexSenseParts))
	}

	// Parse lex-sense.
	var err error
	result.pos, err = parseDeciUint(lexSenseParts[0])
	if err != nil {
		return nil, err
	}
	result.lexFileNum, err = parseDeciUint(lexSenseParts[1])
	if err != nil {
		return nil, err
	}
	result.lexId, err = parseDeciUint(lexSenseParts[2])
	if err != nil {
		return nil, err
	}
	result.headWord = lexSenseParts[3]
	if result.headWord != "" {
		result.headId, err = parseDeciUint(lexSenseParts[4])
		if err != nil {
			return nil, err
		}
	}

	// Parse example numbers.
	if parts[1] != "" {
		numParts := strings.Split(parts[1], ",")
		nums := make([]int, len(numParts))
		for i := range numParts {
			nums[i], err = parseDeciUint(numParts[i])
			if err != nil {
				return nil, err
			}
		}
		result.exampleIds = nums
	}

	return result, nil
}

// ----- EXCEPTION PARSING ----------------------------------------------------

func parseExceptionFiles(path string) (map[string][]string, error) {
	result := map[string][]string{}
	for file, pos := range exceptionFiles {
		f, err := os.Open(filepath.Join(path, file))
		if err != nil {
			return nil, fmt.Errorf("%s: %v", file, err)
		}
		err = parseExceptionFile(f, pos, result)
		f.Close()
		if err != nil {
			return nil, fmt.Errorf("%s: %v", file, err)
		}
	}
	return result, nil
}

// Parses a single exception file. Adds keys to out that point to already
// existing values.
func parseExceptionFile(in io.Reader, pos string, out map[string][]string,
) error {
	scanner := bufio.NewScanner(in)

	// For each line.
	lineNum := 0
	for scanner.Scan() {
		lineNum++
		line := scanner.Text()
		parts := strings.Split(line, " ")
		if len(parts) < 2 {
			return fmt.Errorf("line %d: Bad number of fields: %d, expected 2",
				lineNum, len(parts))
		}

		for i := range parts {
			parts[i] = pos + "." + parts[i]
		}
		out[parts[0]] = parts[1:]
	}

	return scanner.Err()
}

// ----- DATA PARSING ---------------------------------------------------------

// Parses all the data files and returns the 'Synset' field for the Wordnet
// object. Path is data root directory. Example is a map from word sense to
// example IDs.
func parseDataFiles(path string, examples map[string][]int) (
	map[string]*Synset, error) {
	result := map[string]*Synset{}
	for file, pos := range dataFiles {
		f, err := os.Open(filepath.Join(path, file))
		if err != nil {
			return nil, fmt.Errorf("%s: %v", file, err)
		}
		err = parseDataFile(f, pos, examples, result)
		if err != nil {
			return nil, fmt.Errorf("%s: %v", file, err)
		}
	}
	return result, nil
}

// Parses a single data file. Path is the data file. Pos is the POS that this
// file represents. Example is a map from word sense to example IDs. Updates
// out with parsed data.
func parseDataFile(in io.Reader, pos string, examples map[string][]int,
	out map[string]*Synset) error {
	scanner := bufio.NewScanner(in)

	// For each line.
	lineNum := 0
	for scanner.Scan() {
		line := scanner.Text()
		lineNum++
		if strings.HasPrefix(line, "  ") { // Copyright line.
			continue
		}

		// Parse.
		raw, err := parseDataLine(line, pos == "v")
		if err != nil {
			return fmt.Errorf("Line %d: %v", lineNum, err)
		}

		// Assign.
		nice := rawSynsetToNiceSynset(raw)
		key := fmt.Sprintf("%v%v", pos, raw.synsetOffset)
		out[key] = nice

		// Handle examples.
		for i, word := range raw.word {
			key := fmt.Sprintf("%s.%d.%d", word.word, raw.lexFileNum,
				word.lexId)
			//fmt.Println(key)
			for _, exampleId := range examples[key] {
				nice.Example = append(nice.Example, &Example{i, exampleId})
			}
		}
	}

	return scanner.Err()
}

// Converts a raw parsed synset to the exported type.
func rawSynsetToNiceSynset(raw *rawSynset) *Synset {
	result := &Synset{
		raw.synsetOffset,
		raw.ssType,
		make([]string, len(raw.word)),
		make([]*Pointer, len(raw.ptr)),
		raw.frame,
		raw.gloss,
		nil,
	}
	for _, frame := range result.Frame {
		frame.WordNumber-- // Switch from 1-based to 0-based.
	}
	for i, word := range raw.word {
		result.Word[i] = word.word
	}
	for i, rawPtr := range raw.ptr {
		result.Pointer[i] = &Pointer{
			rawPtr.symbol,
			fmt.Sprintf("%v%v", rawPtr.pos, rawPtr.synsetOffset),
			rawPtr.source - 1, // Switch from 1-based to 0-based.
			rawPtr.target - 1, // Switch from 1-based to 0-based.
		}
	}

	return result
}

// Represents a single line in a data file.
type rawSynset struct {
	synsetOffset string
	lexFileNum   int
	ssType       string
	word         []*rawWord
	ptr          []*rawPointer
	frame        []*Frame
	gloss        string
}

type rawPointer struct {
	symbol       string
	synsetOffset string
	pos          string
	source       int // 1-based.
	target       int // 1-based.
}

type rawWord struct {
	word  string
	lexId int
}

// Accepted synset types.
var ssTypes = map[string]bool{
	"n": true,
	"v": true,
	"a": true,
	"s": true,
	"r": true,
}

// TODO(amit): Convert underscores in words to spaces.

// Parses a single line in a data file. hasFrames is true only for the verb
// file.
func parseDataLine(line string, hasFrames bool) (*rawSynset, error) {
	result := &rawSynset{}
	var err error
	parts := strings.Split(strings.Trim(line, " "), " ")
	if len(parts) < 6 {
		return nil, fmt.Errorf("too few fields: %d, expected at "+
			"least 6", len(parts))
	}

	// Parse beginning of line.
	result.synsetOffset = parts[0]
	result.lexFileNum, err = parseDeciUint(parts[1])
	if err != nil {
		return nil, err
	}

	if !ssTypes[parts[2]] {
		return nil, fmt.Errorf("unrecognized ss_type: %s", parts[2])
	}
	result.ssType = parts[2]

	// Parse words.
	wordCount, err := parseHexaUint(parts[3])
	if err != nil {
		return nil, err
	}
	parts = parts[4:]
	if len(parts) < 2*wordCount+2 {
		return nil, fmt.Errorf("too few fields for words: %d, expected at "+
			"least %d", len(parts), 2*wordCount+2)
	}
	result.word = make([]*rawWord, wordCount)

	for i := 0; i < wordCount; i++ {
		word := &rawWord{}
		word.word = parts[0]
		lexId, err := parseHexaUint(parts[1])
		if err != nil {
			return nil, err
		}
		word.lexId = lexId
		result.word[i] = word
		parts = parts[2:]
	}

	// Parse pointers.
	ptrCount, err := parseDeciUint(parts[0])
	if err != nil {
		return nil, err
	}
	parts = parts[1:]
	if len(parts) < 4*ptrCount+1 {
		return nil, fmt.Errorf("too few fields for pointers: %d, expected "+
			"at least %d", len(parts), 4*ptrCount+1)
	}
	result.ptr = make([]*rawPointer, ptrCount)

	for i := 0; i < ptrCount; i++ {
		ptr := &rawPointer{}
		ptr.symbol = parts[0]
		ptr.synsetOffset = parts[1]
		ptr.pos = parts[2]

		if len(parts[3]) != 4 {
			return nil, fmt.Errorf("bad pointer source/target field: %s",
				parts[3])
		}
		ptr.source, err = parseHexaUint(parts[3][:2])
		if err != nil {
			return nil, err
		}
		ptr.target, err = parseHexaUint(parts[3][2:])
		if err != nil {
			return nil, err
		}
		result.ptr[i] = ptr

		parts = parts[4:]
	}

	// Parse frames.
	if hasFrames {
		frameCount, err := parseDeciUint(parts[0])
		if err != nil {
			return nil, err
		}
		parts = parts[1:]
		if len(parts) < 3*frameCount+1 {
			return nil, fmt.Errorf("too few fields for frames: %d, expected "+
				"at least %d", len(parts), 3*frameCount+1)
		}

		result.frame = make([]*Frame, frameCount)
		for i := range result.frame {
			f, err := parseDeciUint(parts[1])
			if err != nil {
				return nil, err
			}
			w, err := parseHexaUint(parts[2])
			if err != nil {
				return nil, err
			}
			result.frame[i] = &Frame{w, f}
			parts = parts[3:]
		}
	}

	// Parse glossary.
	if parts[0] != "|" {
		return nil, fmt.Errorf("expected '|' at end of fields, but found "+
			"'%s'", parts[0])
	}
	result.gloss = strings.Join(parts[1:], " ")

	return result, nil
}

// ----- UTILS ----------------------------------------------------------------

// Now what in the world were they thinking when they put hexa and decimal in
// the same format? Academics and code. -_-

func parseHexaUint(s string) (int, error) {
	i, err := strconv.ParseUint(s, 16, 0)
	return int(i), err
}

func parseDeciUint(s string) (int, error) {
	i, err := strconv.ParseUint(s, 10, 0)
	return int(i), err
}

// Pointer symbol meanings.
const (
	Antonym                   = "!"
	Hypernym                  = "@"
	InstanceHypernym          = "@i"
	Hyponym                   = "~"
	InstanceHyponym           = "~i"
	MemberHolonym             = "#m"
	SubstanceHolonym          = "#s"
	PartHolonym               = "#p"
	MemberMeronym             = "%m"
	SubstanceMeronym          = "%s"
	PartMeronym               = "%p"
	Attribute                 = "="
	DerivationallyRelatedForm = "+"
	DomainOfSynsetTopic       = ";c"
	MemberOfThisDomainTopic   = "-c"
	DomainOfSynsetRegion      = ";r"
	MemberOfThisDomainRegion  = "-r"
	DomainOfSynsetUsage       = ";u"
	MemberOfThisDomainUsage   = "-u"
	Entailment                = "*"
	Cause                     = ">"
	AlsoSee                   = "^"
	VerbGroup                 = "$"
	SimilarTo                 = "&"
	ParticipleOfVerb          = "<"
	Pertainym                 = "\\"
	DerivedFromAdjective      = "\\"
)
