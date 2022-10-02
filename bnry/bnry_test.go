package bnry

import (
	"reflect"
	"testing"

	"golang.org/x/exp/slices"
)

func TestMarshal(t *testing.T) {
	a := byte(113)
	b := uint64(2391278932173219)
	c := "amit"
	d := int16(10000)
	e := []int32{1, 11, 100, 433223}
	f := true
	g := false
	buf, err := MarshalBinary(a, b, c, d, e, f, g)
	if err != nil {
		t.Fatalf("MarshalBinary(%v,%v,%v,%v,%v,%v,%v) failed: %v",
			a, b, c, d, e, f, g, err)
	}
	var aa byte
	var bb uint64
	var cc string
	var dd int16
	var ee []int32
	var ff bool
	var gg bool
	if err := UnmarshalBinary(
		buf, &aa, &bb, &cc, &dd, &ee, &ff, &gg); err != nil {
		t.Fatalf("UnmarshalBinary(%v) failed: %v", buf, err)
	}
	if aa != a {
		t.Errorf("UnmarshalBinary(...)=%v, want %v", aa, a)
	}
	if bb != b {
		t.Errorf("UnmarshalBinary(...)=%v, want %v", bb, b)
	}
	if cc != c {
		t.Errorf("UnmarshalBinary(...)=%v, want %v", cc, c)
	}
	if dd != d {
		t.Errorf("UnmarshalBinary(...)=%v, want %v", dd, d)
	}
	if !slices.Equal(ee, e) {
		t.Errorf("UnmarshalBinary(...)=%v, want %v", ee, e)
	}
	if ff != f {
		t.Errorf("UnmarshalBinary(...)=%v, want %v", dd, d)
	}
	if gg != g {
		t.Errorf("UnmarshalBinary(...)=%v, want %v", dd, d)
	}
}

func FuzzMarshal(f *testing.F) {
	f.Add(uint8(1), int16(1), uint32(1), int64(1), "", true, 1.0, float32(1))
	f.Fuzz(func(t *testing.T, a uint8, b int16, c uint32, d int64, e string,
		g bool, h float64, i float32) {
		buf, err := MarshalBinary(a, b, c, d, e, g, h, i)
		if err != nil {
			t.Fatal(err)
		}
		var (
			aa uint8
			bb int16
			cc uint32
			dd int64
			ee string
			gg bool
			hh float64
			ii float32
		)
		err = UnmarshalBinary(buf, &aa, &bb, &cc, &dd, &ee, &gg, &hh, &ii)
		if err != nil {
			t.Fatal(err)
		}
		if aa != a {
			t.Fatalf("got %v, want %v", aa, a)
		}
		if bb != b {
			t.Fatalf("got %v, want %v", bb, b)
		}
		if cc != c {
			t.Fatalf("got %v, want %v", cc, c)
		}
		if dd != d {
			t.Fatalf("got %v, want %v", dd, d)
		}
		if ee != e {
			t.Fatalf("got %v, want %v", ee, e)
		}
		if gg != g {
			t.Fatalf("got %v, want %v", gg, g)
		}
		if hh != h {
			t.Fatalf("got %v, want %v", hh, h)
		}
		if ii != i {
			t.Fatalf("got %v, want %v", ii, i)
		}
	})
}

func TestMarshal_slices(t *testing.T) {
	a := []uint32{321321, 213, 4944}
	b := []string{"this", "is", "", "a", "slice"}
	c := []int8{100, 9, 0, -21}
	d := []bool{true, false, false, true, true}
	buf, err := MarshalBinary(slices.Clone(a), slices.Clone(b),
		slices.Clone(c), slices.Clone(d))
	if err != nil {
		t.Fatal("MarshalBinary(...) failed:", err)
	}
	var (
		aa []uint32
		bb []string
		cc []int8
		dd []bool
	)
	err = UnmarshalBinary(buf, &aa, &bb, &cc, &dd)
	if err != nil {
		t.Fatal("UnmarshalBinary(...) failed:", err)
	}
	inputs := []any{a, b, c, d}
	outputs := []any{aa, bb, cc, dd}
	for i := range inputs {
		if !reflect.DeepEqual(inputs[i], outputs[i]) {
			t.Fatalf("UnmarshalBinary(...)=%v, want %v", outputs[i], inputs[i])
		}
	}
}
