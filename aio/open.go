// Package aio provides buffered file I/O.
package aio

import (
	"bufio"
	"compress/gzip"
	"io"
	"os"
	"path/filepath"
)

// BUG(amit): Append is not yet supported.

const (
	// If true, .gz files are automatically compressed/decompressed.
	GZipSupport = true
)

// OpenRaw opens a file for reading, with a buffer.
func OpenRaw(file string) (*Reader, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	return &Reader{*bufio.NewReader(f), f}, nil
}

// CreateRaw opens a file for writing, with a buffer.
// Erases any previously existing content.
func CreateRaw(file string) (*Writer, error) {
	f, err := os.Create(file)
	if err != nil {
		return nil, err
	}
	return &Writer{*bufio.NewWriter(f), f}, nil
}

// AppendRaw opens a file for writing, with a buffer.
// Appends to previously existing content if any.
func AppendRaw(file string) (*Writer, error) {
	f, err := os.OpenFile(file, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0o644)
	if err != nil {
		return nil, err
	}
	return &Writer{*bufio.NewWriter(f), f}, nil
}

var (
	rsuffixes = map[string]func(io.Reader) (io.Reader, error){}
	wsuffixes = map[string]func(io.WriteCloser) (io.WriteCloser, error){}
)

// Open opens a file for reading, with a buffer.
// Decompresses the data according to the file's suffix.
func Open(file string) (*Reader, error) {
	f, err := OpenRaw(file)
	if err != nil {
		return nil, err
	}
	fn := rsuffixes[filepath.Ext(file)]
	if fn == nil {
		return f, nil
	}
	ff, err := fn(f)
	if err != nil {
		return nil, err
	}
	return &Reader{*bufio.NewReader(ff), f}, nil
}

// Create opens a file for writing, with a buffer.
// Erases any previously existing content.
// Compresses the data according to the file's suffix.
func Create(file string) (*Writer, error) {
	f, err := CreateRaw(file)
	if err != nil {
		return nil, err
	}
	fn := wsuffixes[filepath.Ext(file)]
	if fn == nil {
		return f, nil
	}
	ff, err := fn(f)
	if err != nil {
		return nil, err
	}
	wrapper := &writerWrapper{ff, f}
	return &Writer{*bufio.NewWriter(ff), wrapper}, nil
}

// Append opens a file for writing, with a buffer.
// Appends to previously existing content if any.
// Compresses the data according to the file's suffix.
func Append(file string) (*Writer, error) {
	f, err := AppendRaw(file)
	if err != nil {
		return nil, err
	}
	fn := wsuffixes[filepath.Ext(file)]
	if fn == nil {
		return f, nil
	}
	ff, err := fn(f)
	if err != nil {
		return nil, err
	}
	wrapper := &writerWrapper{ff, f}
	return &Writer{*bufio.NewWriter(ff), wrapper}, nil
}

// AddReadSuffix adds a supported suffix for automatic decompression.
// suffix should include the dot. f should take a raw reader and return a reader
// that decompresses the data.
func AddReadSuffix(suffix string, f func(io.Reader) (io.Reader, error)) {
	rsuffixes[suffix] = f
}

// AddWriteSuffix adds a supported suffix for automatic compression.
// suffix should include the dot. f should take a raw writer and return a writer
// that compresses the data.
func AddWriteSuffix(suffix string, f func(io.WriteCloser) (
	io.WriteCloser, error)) {
	wsuffixes[suffix] = f
}

func init() {
	if GZipSupport {
		AddReadSuffix(".gz", func(r io.Reader) (io.Reader, error) {
			z, err := gzip.NewReader(r)
			return z, err
		})
		AddWriteSuffix(".gz", func(w io.WriteCloser) (io.WriteCloser, error) {
			z, err := gzip.NewWriterLevel(w, 1)
			return z, err
		})
	}
}
