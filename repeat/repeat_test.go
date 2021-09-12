package repeat

import (
	"io"
	"testing"
)

func TestReader(t *testing.T) {
	r := NewReader([]byte("amit"), 2)
	buf := make([]byte, 3)

	n, err := r.Read(buf)
	if err != nil {
		t.Fatalf("Read() failed: %v", err)
	}
	if n != 3 {
		t.Fatalf("Read() n=%v, want 3", n)
	}
	if string(buf) != "ami" {
		t.Fatalf("Read()=%q, want %q", buf, "ami")
	}

	n, err = r.Read(buf)
	if err != nil {
		t.Fatalf("Read() failed: %v", err)
	}
	if n != 3 {
		t.Fatalf("Read() n=%v, want 3", n)
	}
	if string(buf) != "tam" {
		t.Fatalf("Read()=%q, want %q", buf, "tam")
	}

	n, err = r.Read(buf)
	if err != nil {
		t.Fatalf("Read() failed: %v", err)
	}
	if n != 2 {
		t.Fatalf("Read() n=%v, want 2", n)
	}
	if string(buf[:n]) != "it" {
		t.Fatalf("Read()=%q, want %q", buf[:n], "it")
	}

	_, err = r.Read(buf)
	if err != io.EOF {
		t.Fatalf("Read() err=%v, want EOF", err)
	}
}

func TestReader_infinite(t *testing.T) {
	r := NewReader([]byte("amit"), -1)
	buf := make([]byte, 3)

	n, err := r.Read(buf)
	if err != nil {
		t.Fatalf("Read() failed: %v", err)
	}
	if n != 3 {
		t.Fatalf("Read() n=%v, want 3", n)
	}
	if string(buf) != "ami" {
		t.Fatalf("Read()=%q, want %q", buf, "ami")
	}

	n, err = r.Read(buf)
	if err != nil {
		t.Fatalf("Read() failed: %v", err)
	}
	if n != 3 {
		t.Fatalf("Read() n=%v, want 3", n)
	}
	if string(buf) != "tam" {
		t.Fatalf("Read()=%q, want %q", buf, "tam")
	}

	n, err = r.Read(buf)
	if err != nil {
		t.Fatalf("Read() failed: %v", err)
	}
	if n != 3 {
		t.Fatalf("Read() n=%v, want 3", n)
	}
	if string(buf[:n]) != "ita" {
		t.Fatalf("Read()=%q, want %q", buf[:n], "ita")
	}
}
