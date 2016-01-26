// Convenience package for serializing data. A gobz is simply a gzipped gob.
// This package provides simple functions for handling them.
package gobz

import (
	"fmt"
	"os"
	"io"
	"bufio"
	"compress/gzip"
	"encoding/gob"
)

// Encodes a value to the given stream.
func Encode(w io.Writer, obj interface{}) error {
	// Open zip stream.
	z := gzip.NewWriter(w)
	defer z.Close()
	
	// Write data.
	err := gob.NewEncoder(z).Encode(obj)
	if err != nil {
		return fmt.Errorf("Could not encode object: %v", err)
	} else {
		return nil
	}
}

// Decodes a value from the given stream.
func Decode(r io.Reader, obj interface{}) error {
	// Open zip stream.
	z, err := gzip.NewReader(r)
	if err != nil {
		return fmt.Errorf("Could not read gzip: %v", err)
	}
	
	// Read data.
	err = gob.NewDecoder(z).Decode(obj)
	if err != nil {
		return fmt.Errorf("Could not decode object: %v", err)
	}
	
	err = z.Close()
	if err != nil {
		return fmt.Errorf("Could not read gzip: %v", err)
	}
	
	return nil
}

// Writes a value to the given file.
func Save(file string, obj interface{}) error {
	// Open file.
	f, err := os.Create(file)
	if err != nil {
		return fmt.Errorf("Could not open file: %v", err)
	}
	defer f.Close()
	
	b := bufio.NewWriter(f)
	defer b.Flush()
	
	return Encode(b, obj)
}

// Reads a value from the given file.
func Load(file string, obj interface{}) error {
	// Open file.
	f, err := os.Open(file)
	if err != nil {
		return fmt.Errorf("Could not open file: %v", err)
	}
	defer f.Close()
	
	b := bufio.NewReader(f)
	
	return Decode(b, obj)
}
