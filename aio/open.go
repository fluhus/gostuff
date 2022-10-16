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
func OpenRaw(file string) (io.ReadCloser, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	return &readerWrapper{bufio.NewReader(f), f}, nil
}

// CreateRaw opens a file for writing, with a buffer.
// Erases any previously existing content.
func CreateRaw(file string) (io.WriteCloser, error) {
	f, err := os.Create(file)
	if err != nil {
		return nil, err
	}
	b := bufio.NewWriter(f)
	return &writerWrapper{&flusherWrapper{b}, f}, nil
}

var (
	rsuffixes = map[string]func(io.Reader) (io.Reader, error){}
	wsuffixes = map[string]func(io.WriteCloser) (io.WriteCloser, error){}
)

// Open opens a file for reading, with a buffer.
// Decompresses the data according to the file's suffix.
func Open(file string) (io.ReadCloser, error) {
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
	return &readerWrapper{ff, f}, nil
}

// Create opens a file for writing, with a buffer.
// Erases any previously existing content.
// Compresses the data according to the file's suffix.
func Create(file string) (io.WriteCloser, error) {
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
	return &writerWrapper{ff, f}, nil
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
