// Package jsonf provides convenience functions for saving and loading
// JSON-encoded values to and from the hard drive.
//
// Automatically uses gzip streams for gzip files, using the gzipf package.
package jsonf

import (
	"encoding/json"

	"github.com/fluhus/gostuff/gzipf"
)

// Save encodes v as JSON and saves it to the given file.
func Save(file string, v interface{}) error {
	f, err := gzipf.Create(file)
	if err != nil {
		return err
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	return enc.Encode(v)
}

// Load loads a JSON encoded value from the given file and populates v with it.
func Load(file string, v interface{}) error {
	f, err := gzipf.Open(file)
	if err != nil {
		return err
	}
	defer f.Close()
	return json.NewDecoder(f).Decode(v)
}
