package graphs

import (
	"reflect"
	"slices"
	"testing"

	"github.com/fluhus/gostuff/snm"
)

func TestComponents(t *testing.T) {
	edges := [][2]int{
		{0, 1}, {1, 2}, {5, 7}, {6, 9}, {9, 10}, {8, 10}, {7, 8},
	}
	want := [][]int{
		{0, 1, 2}, {3}, {4}, {5, 6, 7, 8, 9, 10}, {11},
	}
	g := New[int]()
	for i := range 12 {
		g.AddVertices(i)
	}
	for _, e := range edges {
		g.AddEdge(e[0], e[1])
	}
	got := g.ConnectedComponents()
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("components(...)=%v, want %v", got, want)
	}
}

func TestComponents_string(t *testing.T) {
	edges := [][2]string{
		{"a", "bb"}, {"eeeee", "dddd"}, {"bb", "ccc"}, {"dddd", "eeeee"},
	}
	want := [][]string{
		{"ffffff"}, {"a", "bb", "ccc"}, {"eeeee", "dddd"},
	}
	g := New[string]()
	g.AddVertices("ffffff")
	for _, e := range edges {
		g.AddEdge(e[0], e[1])
	}
	got := g.ConnectedComponents()
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("components(...)=%v, want %v", got, want)
	}
}

func TestVerticesEdges(t *testing.T) {
	vertices := []string{"ffffff", "bb"}
	edges := [][2]string{
		{"a", "bb"}, {"eeeee", "dddd"}, {"bb", "ccc"}, {"dddd", "eeeee"},
	}

	wantVertices := []string{"a", "bb", "ccc", "dddd", "eeeee", "ffffff"}
	wantEdges := [][2]string{
		{"bb", "a"}, {"bb", "ccc"}, {"eeeee", "dddd"},
	}

	g := New[string]()
	g.AddVertices(vertices...)
	for _, e := range edges {
		g.AddEdge(e[0], e[1])
	}

	gotVertices := snm.Sorted(slices.Collect(g.Vertices()))
	if !slices.Equal(gotVertices, wantVertices) {
		t.Errorf("Vertices()=%q, want %q", gotVertices, wantVertices)
	}

	var gotEdges [][2]string
	for a, b := range g.Edges() {
		gotEdges = append(gotEdges, [2]string{a, b})
	}
	slices.SortFunc(gotEdges, func(a, b [2]string) int {
		return slices.Compare(a[:], b[:])
	})
	if !slices.Equal(gotEdges, wantEdges) {
		t.Errorf("Edges()=%q, want %q", gotEdges, wantEdges)
	}
}
