// Package iterx provides convenience functions for iterators.
package iterx

import (
	"bufio"
	"encoding/csv"
	"io"
	"iter"

	"github.com/fluhus/gostuff/aio"
)

// LinesReader iterates over text lines from a reader.
func LinesReader(r io.Reader) iter.Seq2[string, error] {
	return func(yield func(string, error) bool) {
		sc := bufio.NewScanner(r)
		for sc.Scan() {
			if !yield(sc.Text(), nil) {
				return
			}
		}
		if err := sc.Err(); err != nil {
			yield("", err)
		}
	}
}

// LinesFile iterates over text lines from a reader.
func LinesFile(file string) iter.Seq2[string, error] {
	return func(yield func(string, error) bool) {
		f, err := aio.Open(file)
		if err != nil {
			yield("", err)
			return
		}
		defer f.Close()
		sc := bufio.NewScanner(f)
		for sc.Scan() {
			if !yield(sc.Text(), nil) {
				return
			}
		}
		if err := sc.Err(); err != nil {
			yield("", err)
		}
	}
}

// CSVReader iterates over CSV entries from a reader.
// fn is an optional function for modifying the CSV parser,
// for example for changing the delimiter.
func CSVReader(r io.Reader, fn func(*csv.Reader)) iter.Seq2[[]string, error] {
	return func(yield func([]string, error) bool) {
		c := csv.NewReader(r)
		if fn != nil {
			fn(c)
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

// CSVFile iterates over CSV entries from a file.
// fn is an optional function for modifying the CSV parser,
// for example for changing the delimiter.
func CSVFile(file string, fn func(*csv.Reader)) iter.Seq2[[]string, error] {
	return func(yield func([]string, error) bool) {
		f, err := aio.Open(file)
		if err != nil {
			yield(nil, err)
			return
		}
		defer f.Close()
		c := csv.NewReader(f)
		if fn != nil {
			fn(c)
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

// CollectErr collects the given T's in a slice until the error is non-nil.
func CollectErr[T any](it iter.Seq2[T, error]) ([]T, error) {
	var a []T
	for t, err := range it {
		if err != nil {
			return a, err
		}
		a = append(a, t)
	}
	return a, nil
}
