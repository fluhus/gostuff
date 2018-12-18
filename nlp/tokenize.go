package nlp

import (
	"github.com/agonopol/go-stem"
	"regexp"
	"strings"
)

// Tokenizer splits text into tokens. This regexp represents a single word.
// Changing this regexp will affect the Tokenize function.
var Tokenizer = regexp.MustCompile("\\w([\\w']*\\w)?")

// Tokenize splits a given text to a slice of stemmed, lowercase words. If
// keepStopWords is false, will drop stop words.
func Tokenize(s string, keepStopWords bool) []string {
	s = correctUtf8Punctuation(s)
	s = strings.ToLower(s)
	words := Tokenizer.FindAllString(s, -1)
	var result []string
	for _, word := range words {
		if !keepStopWords && StopWords[word] {
			continue
		}
		result = append(result, Stem(word))
	}

	return result
}

// Stem porter-stems the given word.
func Stem(s string) string {
	if strings.HasSuffix(s, "'s") {
		s = s[:len(s)-2]
	}
	return string(stemmer.Stem([]byte(s)))
}

// correctUtf8Punctuation translates or removes non-ASCII punctuation characters.
func correctUtf8Punctuation(s string) string {
	return strings.Replace(s, "’", "'", -1)
	// TODO(amit): Improve this function with more characters.
}
