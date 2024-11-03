package iterx

import (
	"iter"
	"reflect"
	"slices"
	"testing"
)

func TestUnreader_until(t *testing.T) {
	input := []int{1, 4, 2, 6, 8, 4, 5, 7}
	tests := []struct {
		until int
		want  []int
	}{
		{6, []int{1, 4, 2}}, {4, []int{6, 8}}, {1, []int{4, 5, 7}},
	}
	r := NewUnreader(Slice(input))
	for _, test := range tests {
		var got []int
		for i := range r.Until(func(j int) bool { return j == test.until }) {
			got = append(got, i)
		}
		if !slices.Equal(got, test.want) {
			t.Fatalf("New(%v).Until(%d)=%v, want %v", input, test.until, got, test.want)
		}
	}
}

func TestUnreader_groupBy(t *testing.T) {
	input := []int{1, 4, 2, 6, 9, 4, 5, 7}
	want := [][]int{{1}, {4, 2, 6}, {9}, {4}, {5, 7}}
	var got [][]int
	r := NewUnreader(Slice(input))
	for group := range r.GroupBy(func(i int, j int) bool {
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

func TestUnreader_stop(t *testing.T) {
	input := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}

	stopped := false
	broke := false
	r := NewUnreader(stopIter(input, &stopped))
	for x := range r.Until(func(i int) bool { return i == 5 }) {
		if x == 3 {
			broke = true
			break
		}
	}
	if !broke {
		t.Fatalf("Unreader(%v) never broke", input)
	}
	if !stopped {
		t.Fatalf("Unreader(%v) never called stop", input)
	}

	stopped, broke = false, false
	continued := true
	r = NewUnreader(stopIter(input, &stopped))
	for g := range r.GroupBy(func(i int, j int) bool {
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
		t.Fatalf("Unreader(%v) never broke", input)
	}
	if continued {
		t.Fatalf("Unreader(%v).GroupBy outer loop continued "+
			"after inner loop broke", input)
	}
	if !stopped {
		t.Fatalf("Unreader(%v) never called stop", input)
	}

	stopped = false
	toBreak := false
	r = NewUnreader(stopIter(input, &stopped))
	for g := range r.GroupBy(func(i int, j int) bool {
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
		t.Fatalf("Unreader(%v) never broke", input)
	}
	if !stopped {
		t.Fatalf("Unreader(%v) never called stop", input)
	}
}

func stopIter[T any](s []T, stopped *bool) iter.Seq[T] {
	return func(yield func(T) bool) {
		defer func() {
			*stopped = true
		}()
		for _, x := range s {
			if !yield(x) {
				return
			}
		}
	}
}
