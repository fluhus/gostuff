package csvdec

import (
	"reflect"
	"slices"
	"strings"
	"testing"
)

// TODO(amit): Change output style to got, want.
// TODO(amit): Add tests for bad input.

func TestReadr_struct(t *testing.T) {
	reader := New(strings.NewReader("Hello,1,-1,3.14"))

	got := struct {
		S string
		U uint
		I int
		F float64
	}{}
	want := struct {
		S string
		U uint
		I int
		F float64
	}{"Hello", 1, -1, 3.14}

	err := reader.ReadInto(&got)
	if err != nil {
		t.Fatal("Error decoding struct:", t)
	}

	if got.S != want.S || got.U != want.U || got.I != want.I || got.F != want.F {
		t.Fatal("Wrong decoding. Expected:", want, "Actual:", got)
	}
}

func TestReader_uintSlice(t *testing.T) {
	reader := New(strings.NewReader("2,3,5,7,11,13"))
	var got []uint
	want := []uint{2, 3, 5, 7, 11, 13}

	err := reader.ReadInto(&got)
	if err != nil {
		t.Fatal("Error decoding struct:", t)
	}

	if len(got) != len(want) {
		t.Fatal("Wrong decoding. Expected:", want, "Actual:", got)
	}

	for i := range got {
		if got[i] != want[i] {
			t.Fatal("Wrong decoding. Expected:", want, "Actual:", got)
		}
	}
}

func TestReader_intSlice(t *testing.T) {
	reader := New(strings.NewReader("2,-3,5,-7,11,-13"))
	var got []int
	want := []int{2, -3, 5, -7, 11, -13}

	err := reader.ReadInto(&got)
	if err != nil {
		t.Fatal("Error decoding struct:", t)
	}

	if len(got) != len(want) {
		t.Fatal("Wrong decoding. Expected:", want, "Actual:", got)
	}

	for i := range got {
		if got[i] != want[i] {
			t.Fatal("Wrong decoding. Expected:", want, "Actual:", got)
		}
	}
}

func TestReader_floatSlice(t *testing.T) {
	reader := New(strings.NewReader("2,-3.14,5.5,-7008,0.11,-1.3"))
	var got []float64
	want := []float64{2, -3.14, 5.5, -7008, 0.11, -1.3}

	err := reader.ReadInto(&got)
	if err != nil {
		t.Fatal("Error decoding struct:", t)
	}

	if len(got) != len(want) {
		t.Fatal("Wrong decoding. Expected:", want, "Actual:", got)
	}

	for i := range got {
		if got[i] != want[i] {
			t.Fatal("Wrong decoding. Expected:", want, "Actual:", got)
		}
	}
}

func TestReader_stringSlice(t *testing.T) {
	reader := New(strings.NewReader("yar,har,fiddle,di,dee"))
	var got []string
	want := []string{"yar", "har", "fiddle", "di", "dee"}

	err := reader.ReadInto(&got)
	if err != nil {
		t.Fatal("Error decoding struct:", t)
	}

	if len(got) != len(want) {
		t.Fatal("Wrong decoding. Expected:", want, "Actual:", got)
	}

	for i := range got {
		if got[i] != want[i] {
			t.Fatal("Wrong decoding. Expected:", want, "Actual:", got)
		}
	}
}

func TestReader_structWithSlice(t *testing.T) {
	type T struct {
		S string
		I []int
	}
	want := T{"hello", []int{5, 4, 3, 2, 1}}
	got := T{}
	input := "hello,5,4,3,2,1"
	reader := New(strings.NewReader(input))

	if err := reader.ReadInto(&got); err != nil {
		t.Fatalf("Read(%q) failed: %v, want success", input, err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Read(%q)=%v, want %v", input, got, want)
	}
}

func TestIter(t *testing.T) {
	input := "yar,har,fiddle,di,dee\nbeing,a,pirate,is,alright,to,be"
	reader := New(strings.NewReader(input))
	reader.FieldsPerRecord = -1
	want := [][]string{{"yar", "har", "fiddle", "di", "dee"},
		{"being", "a", "pirate", "is", "alright", "to", "be"}}
	var got [][]string
	for row, err := range Iter[[]string](reader) {
		if err != nil {
			t.Fatalf("Iter(%q) failed: %v", input, err)
		}
		got = append(got, row)
	}
	if !slices.EqualFunc(got, want, slices.Equal) {
		t.Fatalf("Iter(%q)=%v, want %v", input, got, want)
	}
}
