// Package jio provides convenience functions for saving and loading
// JSON-encoded values.
//
// Uses the [aio] package for I/O.
package jio

import (
	"encoding/json"
	"io"

	"github.com/fluhus/gostuff/aio"
)

// Write saves v to the given file, encoded as JSON.
func Write(file string, v interface{}) error {
	f, err := aio.Create(file)
	if err != nil {
		return err
	}
	e := json.NewEncoder(f)
	e.SetIndent("", "  ")
	if err := e.Encode(v); err != nil {
		f.Close()
		return err
	}
	return f.Close()
}

// Read loads a JSON encoded value from the given file and populates v with it.
func Read(file string, v interface{}) error {
	f, err := aio.Open(file)
	if err != nil {
		return err
	}
	defer f.Close()
	return json.NewDecoder(f).Decode(v)
}

// Iter returns an iterator over sequential JSON values in a file.
//
// Note: a file with several independent JSON values is not a valid JSON file.
func Iter[T any](file string) func(yield func(T, error) bool) {
	return func(yield func(T, error) bool) {
		f, err := aio.Open(file)
		if err != nil {
			var t T
			yield(t, err)
			return
		}
		defer f.Close()
		j := json.NewDecoder(f)
		for {
			var t T
			err := j.Decode(&t)
			if err == io.EOF {
				return
			}
			if err != nil {
				yield(t, err)
				return
			}
			if !yield(t, nil) {
				return
			}
		}
	}
}
