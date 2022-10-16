package aio

import (
	"io"
)

// Wraps a reader and its underlying closer.
type readerWrapper struct {
	top    io.Reader
	bottom io.ReadCloser
}

func (r *readerWrapper) Read(p []byte) (int, error) {
	return r.top.Read(p)
}

func (r *readerWrapper) Close() error {
	return r.bottom.Close()
}

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

// A writer with a Flush method, to convert to a Closer.
type flusher interface {
	io.Writer
	Flush() error
}

// Converts a flusher to a Closer.
type flusherWrapper struct {
	f flusher
}

func (f *flusherWrapper) Write(p []byte) (int, error) {
	return f.f.Write(p)
}

func (f *flusherWrapper) Close() error {
	return f.f.Flush()
}
