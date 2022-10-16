package aio

import (
	"bufio"
	"io"
)

// Wraps a writer and its underlying closer.
type writerWrapper struct {
	top, bottom io.WriteCloser
}

func (w *writerWrapper) Write(p []byte) (int, error) {
	return w.top.Write(p)
}

func (w *writerWrapper) Close() error {
	if err := w.top.Close(); err != nil {
		return err
	}
	return w.bottom.Close()
}

type Reader struct {
	bufio.Reader
	r io.ReadCloser
}

func (r *Reader) Close() error {
	return r.r.Close()
}

type Writer struct {
	bufio.Writer
	w io.WriteCloser
}

func (w *Writer) Close() error {
	if err := w.Flush(); err != nil {
		return err
	}
	return w.w.Close()
}
