package wordnet

// Parses an entire wordnet directory. Path is the root of the directory.
// The parser will trverse it and parse the required files.
//
// The parser assumes directory structure is as published.
func Parse(path string) (*Wordnet, error) {
	result := &Wordnet{}
	var err error

	result.Data, err = parseDataFiles(path)
	if err != nil {
		return nil, err
	}
	
	return result, nil
}

