// Package gzipf provides a united interface around gzip and non-gzip files.
// The functions Open, Create, and Append take the file's suffix in
// consideration and use gzip streams if needed.
//
// All I/O is buffered.
package gzipf

import (
	"bufio"
	"compress/gzip"
	"io"
	"os"
	"strings"
)

// A Reader reads data from a file. If the file has a .gz or .gzip suffix, it
// reads decompressed data.
type Reader struct {
	f *os.File
	z *gzip.Reader
	r io.Reader
}

func (r *Reader) Read(p []byte) (n int, err error) {
	return r.r.Read(p)
}

// Close closes the underlying file and gzip stream.
func (r *Reader) Close() error {
	if err := r.f.Close(); err != nil {
		return err
	}
	if r.z != nil {
		if err := r.z.Close(); err != nil {
			return err
		}
	}
	r.f, r.z, r.r = nil, nil, nil
	return nil
}

// Open opens a file for reading. If the file has a .gz or .gzip suffix, its
// data is decompressed.
func Open(path string) (*Reader, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	b := bufio.NewReader(f)
	if !isGzip(path) {
		return &Reader{f, nil, b}, nil
	}
	z, err := gzip.NewReader(b)
	if err != nil {
		f.Close()
		return nil, err
	}
	return &Reader{f, z, z}, nil
}

// A Writer writes data to a file. If the file has a .gz or .gzip suffix, it
// compresses its data.
type Writer struct {
	f *os.File
	b *bufio.Writer
	z *gzip.Writer
	w io.Writer
}

func (w *Writer) Write(p []byte) (n int, err error) {
	return w.w.Write(p)
}

// Close flushes any buffers and closes the underlying streams.
func (w *Writer) Close() error {
	// TODO(amit): Proceed to close all the streams even if there is an error.
	if w.z != nil {
		if err := w.z.Close(); err != nil {
			return err
		}
	}
	if err := w.b.Flush(); err != nil {
		return err
	}
	if err := w.f.Close(); err != nil {
		return err
	}
	return nil
}

// Create opens a file for writing, deleting any existing content.
// If the file has a .gz or .gzip suffix, it writes compressed data.
func Create(path string) (*Writer, error) {
	f, err := os.Create(path)
	if err != nil {
		return nil, err
	}
	b := bufio.NewWriter(f)
	if !isGzip(path) {
		return &Writer{f, b, nil, b}, nil
	}
	z, _ := gzip.NewWriterLevel(b, gzip.BestSpeed)
	return &Writer{f, b, z, z}, nil
}

// Append opens a file for appending, keeping its existing content.
// If the file has a .gz or .gzip suffix, it writes compressed data.
func Append(path string) (*Writer, error) {
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return nil, err
	}
	b := bufio.NewWriter(f)
	if !isGzip(path) {
		return &Writer{f, b, nil, b}, nil
	}
	z, _ := gzip.NewWriterLevel(b, gzip.BestSpeed)
	return &Writer{f, b, z, z}, nil
}

func isGzip(path string) bool {
	return strings.HasSuffix(path, ".gz") || strings.HasSuffix(path, ".gzip")
}
