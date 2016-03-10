// Provides a parser for wordnet's data and basic search operations.
//
// !!! UNDER CONSTRUCTION !!!
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
