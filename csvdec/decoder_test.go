package csvdec

import (
	"reflect"
	"strings"
	"testing"
)

// TODO(amit): Change output style to got, want.
// TODO(amit): Add tests for bad input.

func TestDecoder_struct(t *testing.T) {
	decoder := NewDecoder(strings.NewReader("Hello,1,-1,3.14"))

	act := struct {
		S string
		U uint
		I int
		F float64
	}{}
	exp := struct {
		S string
		U uint
		I int
		F float64
	}{"Hello", 1, -1, 3.14}

	err := decoder.Decode(&act)
	if err != nil {
		t.Fatal("Error decoding struct:", t)
	}

	if act.S != exp.S || act.U != exp.U || act.I != exp.I || act.F != exp.F {
		t.Fatal("Wrong decoding. Expected:", exp, "Actual:", act)
	}
}

func TestDecoder_uintSlice(t *testing.T) {
	decoder := NewDecoder(strings.NewReader("2,3,5,7,11,13"))
	var act []uint
	exp := []uint{2, 3, 5, 7, 11, 13}

	err := decoder.Decode(&act)
	if err != nil {
		t.Fatal("Error decoding struct:", t)
	}

	if len(act) != len(exp) {
		t.Fatal("Wrong decoding. Expected:", exp, "Actual:", act)
	}

	for i := range act {
		if act[i] != exp[i] {
			t.Fatal("Wrong decoding. Expected:", exp, "Actual:", act)
		}
	}
}

func TestDecoder_intSlice(t *testing.T) {
	decoder := NewDecoder(strings.NewReader("2,-3,5,-7,11,-13"))
	var act []int
	exp := []int{2, -3, 5, -7, 11, -13}

	err := decoder.Decode(&act)
	if err != nil {
		t.Fatal("Error decoding struct:", t)
	}

	if len(act) != len(exp) {
		t.Fatal("Wrong decoding. Expected:", exp, "Actual:", act)
	}

	for i := range act {
		if act[i] != exp[i] {
			t.Fatal("Wrong decoding. Expected:", exp, "Actual:", act)
		}
	}
}

func TestDecoder_floatSlice(t *testing.T) {
	decoder := NewDecoder(strings.NewReader("2,-3.14,5.5,-7008,0.11,-1.3"))
	var act []float64
	exp := []float64{2, -3.14, 5.5, -7008, 0.11, -1.3}

	err := decoder.Decode(&act)
	if err != nil {
		t.Fatal("Error decoding struct:", t)
	}

	if len(act) != len(exp) {
		t.Fatal("Wrong decoding. Expected:", exp, "Actual:", act)
	}

	for i := range act {
		if act[i] != exp[i] {
			t.Fatal("Wrong decoding. Expected:", exp, "Actual:", act)
		}
	}
}

func TestDecoder_stringSlice(t *testing.T) {
	decoder := NewDecoder(strings.NewReader("yar,har,fiddle,di,dee"))
	var act []string
	exp := []string{"yar", "har", "fiddle", "di", "dee"}

	err := decoder.Decode(&act)
	if err != nil {
		t.Fatal("Error decoding struct:", t)
	}

	if len(act) != len(exp) {
		t.Fatal("Wrong decoding. Expected:", exp, "Actual:", act)
	}

	for i := range act {
		if act[i] != exp[i] {
			t.Fatal("Wrong decoding. Expected:", exp, "Actual:", act)
		}
	}
}

func TestDecoder_structWithSlice(t *testing.T) {
	type T struct {
		S string
		I []int
	}
	want := T{"hello", []int{5, 4, 3, 2, 1}}
	got := T{}
	input := "hello,5,4,3,2,1"
	decoder := NewDecoder(strings.NewReader(input))

	if err := decoder.Decode(&got); err != nil {
		t.Fatalf("Decode(%q) failed: %v, want success", input, err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Decode(%q)=%v, want %v", input, got, want)
	}
}
