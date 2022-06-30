package minhash

import (
	"reflect"
	"sort"
	"testing"

	"github.com/fluhus/gostuff/gnum"
)

func TestCollection(t *testing.T) {
	tests := []struct {
		n     int
		input []uint64
		want  []uint64
	}{
		{
			3,
			[]uint64{1, 2, 2, 2, 2, 1, 1, 3, 3, 3, 1, 2, 3, 1, 3, 3, 2},
			[]uint64{1, 2, 3},
		},
		{
			3,
			[]uint64{1, 2, 3, 4, 5, 6, 7, 8, 9},
			[]uint64{1, 2, 3},
		},
		{
			3,
			[]uint64{9, 8, 7, 6, 5, 4, 3, 2, 1},
			[]uint64{1, 2, 3},
		},
		{
			5,
			[]uint64{40, 19, 55, 10, 32, 1, 100, 5, 99, 16, 16},
			[]uint64{1, 5, 10, 16, 19},
		},
	}
	for _, test := range tests {
		mh := New[uint64](test.n)
		for _, k := range test.input {
			mh.Push(k)
		}
		got := mh.View()
		sort.Slice(got, func(i, j int) bool {
			return got[i] < got[j]
		})
		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("New(%d).Push(%v)=%v, want %v",
				test.n, test.input, got, test.want)
		}
	}
}

func TestJSON(t *testing.T) {
	input := New[int](5)
	input.Push(1)
	input.Push(4)
	input.Push(9)
	input.Push(16)
	input.Push(25)
	input.Push(36)
	jsn, err := input.MarshalJSON()
	if err != nil {
		t.Fatalf("MinHash(1,4,9,16,25,36).MarshalJSON() failed: %v", err)
	}
	got := New[int](2)
	err = got.UnmarshalJSON(jsn)
	if err != nil {
		t.Fatalf("UnmarshalJSON(%q) failed: %v", jsn, err)
	}
	if !reflect.DeepEqual(got, input) {
		t.Fatalf("UnmarshalJSON(%q)=%v, want %v", jsn, got, input)
	}
}

func TestJaccard(t *testing.T) {
	tests := []struct {
		a, b []uint64
		k    int
		want float64
	}{
		{[]uint64{1, 2, 3}, []uint64{1, 2, 3}, 3, 1},
		{[]uint64{1, 2, 3}, []uint64{2, 3, 4}, 3, 2.0 / 3.0},
		{[]uint64{2, 3, 4}, []uint64{1, 2, 3}, 3, 2.0 / 3.0},
		{[]uint64{1, 2, 3, 4, 5}, []uint64{1, 3, 5}, 5, 0.6},
	}
	for _, test := range tests {
		a, b := New[uint64](test.k), New[uint64](test.k)
		for _, i := range test.a {
			a.Push(i)
		}
		for _, i := range test.b {
			b.Push(i)
		}
		a.Sort()
		b.Sort()
		if got := a.Jaccard(b); gnum.Abs(got-test.want) > 0.00001 {
			t.Errorf("Jaccard(%v,%v)=%f, want %f",
				test.a, test.b, got, test.want)
		}
	}
}
