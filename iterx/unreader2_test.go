package iterx

import (
	"iter"
	"reflect"
	"slices"
	"testing"

	"github.com/fluhus/gostuff/ppln"
)

func TestUnreader2_until(t *testing.T) {
	input := []int{1, 4, 2, 6, 8, 4, 5, 7}
	tests := []struct {
		until int
		want  []int
	}{
		{6, []int{1, 4, 2}}, {4, []int{6, 8}}, {1, []int{4, 5, 7}},
	}
	r := NewUnreader2(ppln.SliceInput(input))
	for _, test := range tests {
		var got []int
		for i, err := range r.Until(func(j int, err error) bool { return j == test.until }) {
			if err != nil {
				t.Fatalf("New(%v).Until(%d) failed: %v", input, test.until, err)
			}
			got = append(got, i)
		}
		if !slices.Equal(got, test.want) {
			t.Fatalf("New(%v).Until(%d)=%v, want %v", input, test.until, got, test.want)
		}
	}
}

func TestUnreader2_groupBy(t *testing.T) {
	input := []int{1, 4, 2, 6, 9, 4, 5, 7}
	want := [][]int{{1}, {4, 2, 6}, {9}, {4}, {5, 7}}
	var got [][]int
	r := NewUnreader2(ppln.SliceInput(input))
	for group := range r.GroupBy(func(i int, ie error, j int, je error) bool {
		return i%2 == j%2
	}) {
		var got1 []int
		for i := range group {
			got1 = append(got1, i)
		}
		got = append(got, got1)
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("New(%v).GroupBy(...)=%v, want %v", input, got, want)
	}
}

func TestUnreader2_stop(t *testing.T) {
	input := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}

	stopped := false
	broke := false
	r := NewUnreader2(stopIter2(input, &stopped))
	for x := range r.Until(func(i int, err error) bool { return i == 5 }) {
		if x == 3 {
			broke = true
			break
		}
	}
	if !broke {
		t.Fatalf("Unreader2(%v) never broke", input)
	}
	if !stopped {
		t.Fatalf("Unreader2(%v) never called stop", input)
	}

	stopped, broke = false, false
	continued := true
	r = NewUnreader2(stopIter2(input, &stopped))
	for g := range r.GroupBy(func(i int, ie error, j int, je error) bool {
		return i/3 == j/3
	}) {
		continued = true
		for x := range g {
			if x == 5 {
				broke = true
				continued = false
				break
			}
		}
	}
	if !broke {
		t.Fatalf("Unreader2(%v) never broke", input)
	}
	if continued {
		t.Fatalf("Unreader2(%v).GroupBy outer loop continued "+
			"after inner loop broke", input)
	}
	if !stopped {
		t.Fatalf("Unreader2(%v) never called stop", input)
	}

	stopped = false
	toBreak := false
	r = NewUnreader2(stopIter2(input, &stopped))
	for g := range r.GroupBy(func(i int, ie error, j int, je error) bool {
		return i/3 == j/3
	}) {
		for x := range g {
			if x == 5 {
				toBreak = true
			}
		}
		if toBreak {
			break
		}
	}
	if !toBreak {
		t.Fatalf("Unreader2(%v) never broke", input)
	}
	if !stopped {
		t.Fatalf("Unreader2(%v) never called stop", input)
	}
}

func stopIter2[T any](s []T, stopped *bool) iter.Seq2[T, error] {
	return func(yield func(T, error) bool) {
		defer func() {
			*stopped = true
		}()
		for _, x := range s {
			if !yield(x, nil) {
				return
			}
		}
	}
}
