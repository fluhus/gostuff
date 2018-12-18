package nlp

// TF-IDF functionality.

import (
	"math"
)

// TfIdf returns the TF-IDF scores of the given corpus. For each documet,
// returns a map from token to TF-IDF score.
//
// TF = count(token in document) / count(all tokens in document)
//
// IDF = log(count(documents) / count(documents with token))
func TfIdf(docTokens [][]string) []map[string]float64 {
	tf := make([]map[string]float64, len(docTokens))
	idf := map[string]float64{}

	// Collect TF and DF.
	for i := range docTokens {
		tf[i] = map[string]float64{}
		for j := range docTokens[i] {
			tf[i][docTokens[i][j]]++
		}
		for token := range tf[i] {
			tf[i][token] /= float64(len(docTokens[i]))
			idf[token]++
		}
	}

	// Turn DF to IDF.
	for token, df := range idf {
		idf[token] = math.Log(float64(len(docTokens)) / df)
	}

	// Turn TF to TF-IDF.
	for i := range tf {
		for token := range tf[i] {
			tf[i][token] *= idf[token]
		}
	}

	return tf
}
