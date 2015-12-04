package tsv

import (
	"strings"
	"testing"
)

func TestDecoder_struct(t *testing.T) {
	decoder := NewDecoder(strings.NewReader("Hello\t1\t-1\t3.14"), 0, 0)

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
	decoder := NewDecoder(strings.NewReader("2\t3\t5\t7\t11\t13"), 0, 0)
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
	decoder := NewDecoder(strings.NewReader("2\t-3\t5\t-7\t11\t-13"), 0, 0)
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
	decoder := NewDecoder(strings.NewReader("2\t-3.14\t5.5\t-7008\t0.11\t-1.3"),
		0, 0)
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
	decoder := NewDecoder(strings.NewReader("yar\thar\tfiddle\tdi\tdee"),
		0, 0)
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

// TODO(amit): Add tests for bad input.
