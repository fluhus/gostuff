// Wordnet parser and interface.
//
// !!! UNDER CONSTRUCTION !!!
//
// Basic usage
//
// The main entry point is the Wordnet type. It holds all the data of a
// wordnet dictionary.
//
// Parts of speech
//
// Some data refers to parts of speech (POS). Everywhere a POS is expected,
// it is a single letter as follows:
//  a: adjective
//  n: noun
//  r: adverb
//  v: verb
//
// Keys in Lemma field
//
// Keys are "pos.lemma". For example the key "n.back" relates to the noun
// "back", and the key "v.back" relates to the verb "back".
//
// Keys in Synset field
//
// These have no human-readable meaning, and should be used blindly for
// matching lemmas to synsets.
//
// Keys are "pos.byte_offset", where byte_offset if the field from the original
// data files. Here it has no meaning as real offset, but as a unique ID.
package wordnet

// Parses an entire wordnet directory. Path is the root of the directory.
// The parser will trverse it and parse the required files.
//
// The parser assumes directory structure is as published.
func Parse(path string) (*Wordnet, error) {
	result := &Wordnet{}
	var err error

	result.Synset, err = parseDataFiles(path)
	if err != nil {
		return nil, err
	}

	result.Lemma, err = parseIndexFiles(path)
	if err != nil {
		return nil, err
	}

	result.Exception, err = parseExceptionFiles(path)
	if err != nil {
		return nil, err
	}

	return result, nil
}
