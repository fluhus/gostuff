// Package jsonf provides convenience functions for saving and loading
// JSON-encoded values to and from the hard drive.
//
// Uses the [aio] package for I/O.
package jsonf

import (
	"encoding/json"

	"github.com/fluhus/gostuff/aio"
)

// Save saves v to the given file, encoded as JSON.
func Save(file string, v interface{}) error {
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

// Load loads a JSON encoded value from the given file and populates v with it.
func Load(file string, v interface{}) error {
	f, err := aio.Open(file)
	if err != nil {
		return err
	}
	defer f.Close()
	return json.NewDecoder(f).Decode(v)
}
