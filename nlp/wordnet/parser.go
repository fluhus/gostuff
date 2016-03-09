package wordnet

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// ----- FILE LISTS -----------------------------------------------------------

var dataFiles = map[string]string{
	"data.adj":  "a",
	"data.adv":  "r",
	"data.noun": "n",
	"data.verb": "v",
}

// ----- INDEX PARSING --------------------------------------------------------

//type lemma struct {

//}

// ----- HIGH-LEVEL DATA PARSING ----------------------------------------------

// TODO(amit): Convert symbols to actual meaningful words?

// Parses all the data files and returns the 'Data' field for the Wordnet
// object. Path is data root directory.
func parseDataFiles(path string) (map[string]*Synset, error) {
	result := map[string]*Synset{}
	for file, pos := range dataFiles {
		err := parseDataFile(filepath.Join(path, file), pos, result)
		if err != nil {
			return nil, fmt.Errorf("%s: %v", file, err)
		}
	}
	return result, nil
}

// Parses a single data file. Path is the data file. Pos is the POS that this
// file represents. Updates out with parsed data.
func parseDataFile(path, pos string, out map[string]*Synset) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)

	// For each line.
	lineNum := 0
	for scanner.Scan() {
		line := scanner.Text()
		lineNum++
		if strings.HasPrefix(line, "  ") { // Copyright line.
			continue
		}

		raw, err := parseDataLine(line, pos == "v")
		if err != nil {
			return fmt.Errorf("Line %d: %v", lineNum, err)
		}

		key := pos + "." + raw.synsetOffset
		out[key] = rawSynsetToNiceSynset(raw)
	}

	return scanner.Err()
}

// Converts a raw parsed synset to the exported type.
func rawSynsetToNiceSynset(raw *rawSynset) *Synset {
	result := &Synset{
		raw.ssType,
		raw.word,
		make([]*DataPtr, len(raw.ptr)),
		raw.frame,
		raw.gloss,
	}

	for i, rawPtr := range raw.ptr {
		result.Ptr[i] = &DataPtr{
			rawPtr.symbol,
			rawPtr.pos + "." + rawPtr.synsetOffset,
			rawPtr.source,
			rawPtr.target,
		}
	}

	return result
}

// ----- RAW DATA PARSING -----------------------------------------------------

type rawSynset struct {
	synsetOffset string
	lexFilenum   int
	ssType       string
	word         []*DataWord
	ptr          []*rawDataPtr
	frame        []*DataFrame
	gloss        string
}

type rawDataPtr struct {
	symbol       string
	synsetOffset string
	pos          string
	source       int // 1-based.
	target       int // 1-based.
}

// Accepted POS types.
var ssTypes = map[string]bool{
	"n": true,
	"v": true,
	"a": true,
	"s": true,
	"r": true,
}

// Parses a single line in a data file. hasFrames is true only for the verb
// file.
func parseDataLine(line string, hasFrames bool) (*rawSynset, error) {
	result := &rawSynset{}
	var err error
	parts := strings.Split(line, " ")
	if len(parts) < 6 {
		return nil, fmt.Errorf("Line is too short: %d fields, expected at "+
			"least 6.", len(parts))
	}

	// Parse beginning of line.
	result.synsetOffset = parts[0]
	result.lexFilenum, err = parseHexaUint(parts[1])
	if err != nil {
		return nil, err
	}

	if !ssTypes[parts[2]] {
		return nil, fmt.Errorf("Unrecognized ss_type: %s", parts[2])
	}
	result.ssType = parts[2]

	// Parse words.
	wordCount, err := parseHexaUint(parts[3])
	if err != nil {
		return nil, err
	}
	parts = parts[4:]
	if len(parts) < 2*wordCount+2 {
		return nil, fmt.Errorf("Too few fields for words: %d, expected at "+
			"least %d.", len(parts), 2*wordCount+2)
	}
	result.word = make([]*DataWord, wordCount)

	for i := 0; i < wordCount; i++ {
		word := parts[0]
		lexId, err := parseHexaUint(parts[1])
		if err != nil {
			return nil, err
		}
		result.word[i] = &DataWord{word, lexId}
		parts = parts[2:]
	}

	// Parse pointers.
	ptrCount, err := parseDeciUint(parts[0])
	if err != nil {
		return nil, err
	}
	parts = parts[1:]
	if len(parts) < 4*ptrCount+1 {
		return nil, fmt.Errorf("Too few fields for pointers: %d, expected "+
			"at least %d.", len(parts), 4*ptrCount+1)
	}
	result.ptr = make([]*rawDataPtr, ptrCount)

	for i := 0; i < ptrCount; i++ {
		ptr := &rawDataPtr{}
		ptr.symbol = parts[0]
		ptr.synsetOffset = parts[1]
		ptr.pos = parts[2]

		if len(parts[3]) != 4 {
			return nil, fmt.Errorf("Bad pointer source/target field: %s",
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
			return nil, fmt.Errorf("Too few fields for frames: %d, expected "+
				"at least %d.", len(parts), 3*frameCount+1)
		}

		result.frame = make([]*DataFrame, frameCount)
		for i := range result.frame {
			f, err := parseDeciUint(parts[1])
			if err != nil {
				return nil, err
			}
			w, err := parseHexaUint(parts[2])
			if err != nil {
				return nil, err
			}
			result.frame[i] = &DataFrame{f, w}
			parts = parts[3:]
		}
	}

	// Parse glossary.
	if parts[0] != "|" {
		return nil, fmt.Errorf("Expected '|' at end of fields, but found "+
			"'%s'.", parts[0])
	}
	result.gloss = strings.Join(parts[1:], " ")

	return result, nil
}

// ----- UTILS ----------------------------------------------------------------

// WTF who uses hexa and decimal in the same data format??
// Who uses hexadecimal in textual data at all??
// Who let those academics write code??
// Ugh, man. Ugh.

func parseHexaUint(s string) (int, error) {
	i, err := strconv.ParseUint(s, 16, 0)
	return int(i), err
}

func parseDeciUint(s string) (int, error) {
	i, err := strconv.ParseUint(s, 10, 0)
	return int(i), err
}
