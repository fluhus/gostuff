// Package repeat implements the repeat-reader. A repeat-reader outputs a given
// constant byte sequence repeatedly.
//
// This package was originally written for testing and profiling parsers.
// A typical use may look something like this:
//
//	input := "some line to be parsed\n"
//	r := NewReader([]byte(input), 1000)
//	parser := myparse.NewParser(r)
//
//	(start profiling)
//	for range parser.Items() {} // Exhaust parser.
//	(stop profiling)
package repeat

import "io"

// Reader outputs a given constant byte sequence repeatedly.
type Reader struct {
	data []byte
	i    int
	n    int
}

// NewReader returns a reader that outputs data n times. If n is negative, repeats
// infinitely. Copies the contents of data.
func NewReader(data []byte, n int) *Reader {
	cp := append(make([]byte, 0, len(data)), data...)
	return &Reader{data: cp, n: n}
}

// Read fills p with repetitions of the reader's data. Writes until p is full or
// until the last repetition was written. Subsequent calls to Read resume from where
// the last repetition stopped. When no more bytes are available, returns 0, EOF.
// Otherwise the error is nil.
func (r *Reader) Read(p []byte) (int, error) {
	if len(p) == 0 {
		return 0, nil
	}
	if r.i == 0 && r.n == 0 {
		return 0, io.EOF
	}
	m := 0
	for {
		n := copy(p, r.data[r.i:])
		r.i += n
		if r.i == len(r.data) {
			r.i = 0
			if r.n > 0 {
				r.n--
			}
		}
		p = p[n:]
		m += n
		if len(p) == 0 || (r.i == 0 && r.n == 0) {
			break
		}
	}
	return m, nil
}

// Close is a no-op. Implements [io.ReadCloser].
func (r *Reader) Close() error {
	return nil
}
