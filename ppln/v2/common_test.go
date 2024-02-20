package ppln

import (
	"slices"
	"testing"
)

func TestRangeInput(t *testing.T) {
	tests := []struct {
		from, to int
		want     []int
	}{
		{0, 5, []int{0, 1, 2, 3, 4}},
		{-2, 3, []int{-2, -1, 0, 1, 2}},
		{0, 0, nil},
		{1, 1, nil},
	}
	for _, test := range tests {
		var got []int
		for i, err := range RangeInput(test.from, test.to) {
			if err != nil {
				t.Fatalf("RangeInput(%d,%d) returned error: %v",
					test.from, test.to, err)
			}
			got = append(got, i)
		}
		if !slices.Equal(got, test.want) {
			t.Fatalf("RangeInput(%d,%d)=%v, want %v",
				test.from, test.to, got, test.want)
		}
	}
}

func TestSliceInput(t *testing.T) {
	want := []int{3, 5, 1, 7}
	var got []int
	for i, err := range SliceInput(slices.Clone(want)) {
		if err != nil {
			t.Fatalf("SliceInput(%v) returned error: %v",
				want, err)
		}
		got = append(got, i)
	}
	if !slices.Equal(got, want) {
		t.Fatalf("SliceInput(%v)=%v, want %v",
			want, got, want)
	}
}
