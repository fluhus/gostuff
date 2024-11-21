package repeat

import (
	"io"
	"testing"
)

func TestReader(t *testing.T) {
	r := NewReader([]byte("amit"), 2)
	buf := make([]byte, 3)
	want := []string{"ami", "tam", "it"}

	for _, w := range want {
		n, err := r.Read(buf)
		if err != nil {
			t.Fatalf("Read() failed: %v", err)
		}
		if got := string(buf[:n]); got != w {
			t.Fatalf("Read()=%q, want %q", got, "ami")
		}
	}
	if _, err := r.Read(buf); err != io.EOF {
		t.Fatalf("Read() err=%v, want EOF", err)
	}
}

func TestReader_infinite(t *testing.T) {
	r := NewReader([]byte("amit"), -1)
	buf := make([]byte, 3)
	want := []string{"ami", "tam", "ita", "mit", "ami", "tam", "ita", "mit"}

	for i, w := range want {
		n, err := r.Read(buf)
		if err != nil {
			t.Fatalf("#%v: Read() failed: %v", i, err)
		}
		if got := string(buf[:n]); got != w {
			t.Fatalf("#%v: Read()=%q, want %q", i, got, "ami")
		}
	}
}
