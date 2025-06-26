package csvx

import (
	"encoding/csv"
	"io"
	"iter"

	"github.com/fluhus/gostuff/aio"
)

// Reader iterates over CSV entries from a reader.
// Applies the given modifiers before iteration.
func Reader(r io.Reader, mods ...ReaderModifier) iter.Seq2[[]string, error] {
	return func(yield func([]string, error) bool) {
		c := csv.NewReader(r)
		for _, mod := range mods {
			mod(c)
		}
		for {
			e, err := c.Read()
			if err == io.EOF {
				return
			}
			if !yield(e, nil) {
				return
			}
		}
	}
}

// File iterates over CSV entries from a file.
// Applies the given modifiers before iteration.
func File(file string, mods ...ReaderModifier) iter.Seq2[[]string, error] {
	return func(yield func([]string, error) bool) {
		f, err := aio.Open(file)
		if err != nil {
			yield(nil, err)
			return
		}
		defer f.Close()
		c := csv.NewReader(f)
		for _, mod := range mods {
			mod(c)
		}
		for {
			e, err := c.Read()
			if err == io.EOF {
				return
			}
			if !yield(e, nil) {
				return
			}
		}
	}
}

// ReaderModifier modifies the settings of a CSV reader
// before iteration starts.
type ReaderModifier = func(*csv.Reader)

// TSV makes the reader use tab as the delimiter.
func TSV(r *csv.Reader) {
	r.Comma = '\t'
}
