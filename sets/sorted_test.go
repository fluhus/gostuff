package sets

import (
	"slices"
	"testing"
)

func TestSorted(t *testing.T) {
	tests := []struct {
		a, b, u, i []int
	}{
		{nil, nil, nil, nil},
		{[]int{1}, nil, []int{1}, nil},
		{nil, []int{2}, []int{2}, nil},
		{[]int{1}, []int{2}, []int{1, 2}, nil},
		{[]int{2}, []int{1}, []int{1, 2}, nil},
		{[]int{1, 3, 5}, []int{3, 4, 5, 6}, []int{1, 3, 4, 5, 6}, []int{3, 5}},
	}
	for _, test := range tests {
		i := SortedIntersection(test.a, test.b)
		u := SortedUnion(test.a, test.b)
		il := SortedIntersectionLen(test.a, test.b)
		ul := SortedUnionLen(test.a, test.b)
		if !slices.Equal(i, test.i) {
			t.Fatalf("SortedIntersection(%v,%v)=%v, want %v",
				test.a, test.b, i, test.i)
		}
		if !slices.Equal(u, test.u) {
			t.Fatalf("SortedUnion(%v,%v)=%v, want %v",
				test.a, test.b, u, test.u)
		}
		if il != len(i) {
			t.Fatalf("SortedIntersectionLen(%v,%v)=%v, want %v",
				test.a, test.b, il, len(i))
		}
		if ul != len(u) {
			t.Fatalf("SortedUnionLen(%v,%v)=%v, want %v",
				test.a, test.b, ul, len(u))
		}
	}
}
