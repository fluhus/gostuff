// Package gobz provides convenience functions for serializing data.
// A gobz is simply a gzipped gob.
package gobz

import (
	"bufio"
	"compress/gzip"
	"encoding/gob"
	"fmt"
	"io"
	"os"
)

// Encode encodes a value to the given stream.
func Encode(w io.Writer, obj interface{}) error {
	// Open zip stream.
	z := gzip.NewWriter(w)
	defer z.Close()

	// Write data.
	err := gob.NewEncoder(z).Encode(obj)
	if err != nil {
		return fmt.Errorf("failed to encode object: %v", err)
	}
	return nil
}

// Decode decodes a value from the given stream.
func Decode(r io.Reader, obj interface{}) error {
	// Open zip stream.
	z, err := gzip.NewReader(r)
	if err != nil {
		return fmt.Errorf("failed to read gzip: %v", err)
	}

	// Read data.
	err = gob.NewDecoder(z).Decode(obj)
	if err != nil {
		return fmt.Errorf("failed to decode object: %v", err)
	}

	err = z.Close()
	if err != nil {
		return fmt.Errorf("failed to read gzip: %v", err)
	}

	return nil
}

// Save writes a value to the given file.
func Save(file string, obj interface{}) error {
	// Open file.
	f, err := os.Create(file)
	if err != nil {
		return fmt.Errorf("failed to open file: %v", err)
	}
	defer f.Close()

	b := bufio.NewWriter(f)
	defer b.Flush()

	return Encode(b, obj)
}

// Load reads a value from the given file.
func Load(file string, obj interface{}) error {
	// Open file.
	f, err := os.Open(file)
	if err != nil {
		return fmt.Errorf("failed to open file: %v", err)
	}
	defer f.Close()

	b := bufio.NewReader(f)

	return Decode(b, obj)
}
