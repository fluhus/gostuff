package nlp

import (
	"github.com/agonopol/go-stem"
	"regexp"
	"strings"
)

// Splits text into tokens.
var tokenSplitter = regexp.MustCompile("\\w([\\w']*\\w)?")

// Splits a given text to a slice of stemmed, lowercase words. If keepStopWords
// is false, will drop stop words.
func Tokenize(s string, keepStopWords bool) []string {
	s = correctUtf8Punctuation(s)
	s = strings.ToLower(s)
	words := tokenSplitter.FindAllString(s, -1)
	var result []string
	for _, word := range words {
		if !keepStopWords && StopWords[word] {
			continue
		}
		result = append(result, Stem(word))
	}

	return result
}

// Porter-stems the given word.
func Stem(s string) string {
	if strings.HasSuffix(s, "'s") {
		s = s[:len(s)-2]
	}
	return string(stemmer.Stem([]byte(s)))
}

// Translates or removes non-ASCII punctuation characters.
func correctUtf8Punctuation(s string) string {
	return strings.Replace(s, "’", "'", -1)
	// TODO(amit): Improve this function with more characters.
}
