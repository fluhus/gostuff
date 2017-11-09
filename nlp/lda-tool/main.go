// Command lda-tool performs LDA on the input documents.
package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"regexp"
	"runtime"
	"strings"

	"github.com/fluhus/gostuff/flug"
	"github.com/fluhus/gostuff/maps"
	"github.com/fluhus/gostuff/nlp"
)

func main() {
	parseArgs()

	// Read input and perform LDA.
	fmt.Fprintln(os.Stdout, "Reading documents from stdin. Run with no arguments for usage help.")
	docs, err := readDocs(os.Stdin)
	if err != nil {
		die("Error: failed to read input:", err)
	}
	lda, _ := nlp.LdaThreads(docs, args.K, args.NumThreads)

	// Print output.
	if args.Json {
		j, _ := json.MarshalIndent(lda, "", "\t")
		fmt.Println(string(j))
	} else {
		for _, w := range maps.Keys(lda).([]string) {
			fmt.Print(w)
			for _, x := range lda[w] {
				fmt.Printf(" %v", x)
			}
			fmt.Println()
		}
	}
}

// readDocs reads documents, one per line, from the input reader.
// It splits and lowercases the documents, and returns them as a 2d slice.
func readDocs(r io.Reader) ([][]string, error) {
	wordsRe := regexp.MustCompile("\\w+")
	scanner := bufio.NewScanner(r)
	var result [][]string
	for scanner.Scan() {
		w := wordsRe.FindAllString(strings.ToLower(scanner.Text()), -1)
		
		// Copy line to a lower capacity slice, to reduce memory usage.
		result = append(result, make([]string, len(w)))
		copy(result[len(result)-1], w)
	}
	if scanner.Err() != nil {
		return nil, scanner.Err()
	}
	return result, nil
}

// die reports an error message and exits with error code 2.
// Arguments are treated like Println.
func die(a ...interface{}) {
	fmt.Fprintln(os.Stderr, a...)
	os.Exit(2)
}

var args = struct {
	K          int  `flug:"k,Number of topics."`
	NumThreads int  `flug:"t,Number of therads to use. (default: number of CPUs)"`
	Json       bool `flug:"j,Output as JSON instead of default format."`
}{0, 0, false}

// parseArgs parses the program's arguments and validates them.
// Exits with an error message upon validation error.
func parseArgs() {
	flug.Register(&args)
	flag.Parse()
	if len(os.Args) == 1 {
		fmt.Fprintln(os.Stderr, help)
		flag.PrintDefaults()
		os.Exit(1)
	}
	if args.K < 1 {
		die("Error: invalid k:", args.K)
	}
	if args.NumThreads < 0 {
		die("Error: invalid number of threads:", args.NumThreads)
	}
	if args.NumThreads == 0 {
		args.NumThreads = runtime.NumCPU()
	}
}

var help = `Performs LDA on the given documents.

Input is read from the standard input. Format is one document per line.
Documents will be lowercased and normalized (spaces and punctuation omitted).

Output is printed to the standard output. Format is one word per line.
Each word is followed by K numbers, the i'th number represents the likelihood
of the i'th topic to emit that word.
`
