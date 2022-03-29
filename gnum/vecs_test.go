package gnum

import (
	"testing"

	"golang.org/x/exp/slices"
)

func TestOnes(t *testing.T) {
	want := []int8{1, 1, 1, 1}
	got := Ones[[]int8](4)
	if !slices.Equal(got, want) {
		t.Errorf("Ones(4) = %v, want %v", got, want)
	}
}

func TestAdd(t *testing.T) {
	a := []uint{4, 6}
	b := []uint{2, 3}
	got := Add(nil, a, b)
	want := []uint{6, 9}
	if !slices.Equal(got, want) {
		t.Errorf("Add(nil, %v, %v)=%v, want %v", a, b, got, want)
	}
	if want := []uint{4, 6}; !slices.Equal(a, want) {
		t.Errorf("a=%v, want %v", a, want)
	}
	if want := []uint{2, 3}; !slices.Equal(b, want) {
		t.Errorf("b=%v, want %v", b, want)
	}
}

func TestAdd_inplace(t *testing.T) {
	a := []uint{4, 6}
	b := []uint{2, 3}
	got := Add(a, b)
	want := []uint{6, 9}
	if !slices.Equal(got, want) {
		t.Errorf("Add(nil, %v, %v)=%v, want %v", a, b, got, want)
	}
	if !slices.Equal(a, want) {
		t.Errorf("a=%v, want %v", a, want)
	}
	if want := []uint{2, 3}; !slices.Equal(b, want) {
		t.Errorf("b=%v, want %v", b, want)
	}
}
