package graphs

import (
	"cmp"
	"slices"
	"testing"

	"github.com/fluhus/gostuff/snm"
)

func TestBFS(t *testing.T) {
	edges := [][2]string{
		{"a", "b"},
		{"a", "c"},
		{"b", "d"},
		{"b", "e"},
		{"c", "f"},
		{"e", "g"},
		{"e", "h"},
	}
	want := [][]string{
		{"a"},
		{"b", "c"},
		{"d", "e", "f"},
		{"g", "h"},
	}

	snm.Shuffle(edges)
	g := New[string]()
	for _, e := range edges {
		g.AddEdge(e[0], e[1])
	}

	got := slices.Collect(g.BFS("a"))
	if wantLen := sumLens(want); len(got) != wantLen {
		t.Fatalf("BFS(...) len=%v, want %v", len(got), wantLen)
	}
	ggot := groupLike(got, want)
	for _, g := range ggot {
		slices.Sort(g)
	}
	if !slices.EqualFunc(ggot, want, slices.Equal) {
		t.Fatalf("BFS(...)=%v, want %v", got, want)
	}
}

func TestBFS_loop(t *testing.T) {
	edges := [][2]string{
		{"a", "b"},
		{"b", "c"},
		{"b", "d"},
		{"c", "e"},
		{"d", "f"},
		{"e", "g"},
		{"f", "g"},
	}
	want := [][]string{
		{"a"},
		{"b"},
		{"c", "d"},
		{"e", "f"},
		{"g"},
	}

	snm.Shuffle(edges)
	g := New[string]()
	for _, e := range edges {
		g.AddEdge(e[0], e[1])
	}

	got := slices.Collect(g.BFS("a"))
	if wantLen := sumLens(want); len(got) != wantLen {
		t.Fatalf("BFS(...) len=%v, want %v", len(got), wantLen)
	}
	ggot := groupLike(got, want)
	for _, g := range ggot {
		slices.Sort(g)
	}
	if !slices.EqualFunc(ggot, want, slices.Equal) {
		t.Fatalf("BFS(...)=%v, want %v", got, want)
	}
}

func TestDFS(t *testing.T) {
	edges := [][2]string{
		{"a", "b"},
		{"a", "c"},
		{"b", "d"},
		{"b", "e"},
		{"c", "f"},
		{"c", "g"},
	}
	want := [][]string{
		{"a"},
		{"b", "d", "e"},
		{"c", "f", "g"},
	}

	snm.Shuffle(edges)
	g := New[string]()
	for _, e := range edges {
		g.AddEdge(e[0], e[1])
	}

	got := slices.Collect(g.DFS("a"))
	if wantLen := sumLens(want); len(got) != wantLen {
		t.Fatalf("DFS(...) len=%v, want %v", len(got), wantLen)
	}
	ggot := groupLike(got, want)
	for _, g := range ggot {
		slices.Sort(g[1:])
	}
	slices.SortFunc(ggot, func(a, b []string) int {
		return cmp.Compare(a[0], b[0])
	})
	if !slices.EqualFunc(ggot, want, slices.Equal) {
		t.Fatalf("DFS(...)=%v, want %v", got, want)
	}
}

func TestDFS_loop(t *testing.T) {
	edges := [][2]string{
		{"a", "b"},
		{"b", "c"},
		{"c", "d"},
		{"d", "a"},
	}
	want := [][]string{
		{"b", "c", "d", "a"},
		{"b", "a", "d", "c"},
	}

	snm.Shuffle(edges)
	g := New[string]()
	for _, e := range edges {
		g.AddEdge(e[0], e[1])
	}

	got := slices.Collect(g.DFS("b"))
	found := slices.IndexFunc(want, func(a []string) bool {
		return slices.Equal(a, got)
	})
	if found == -1 {
		t.Fatalf("DFS(...)=%v, want one of %v", got, want)
	}
}

// Returns the sum of lengths of slices.
func sumLens(a [][]string) int {
	i := 0
	for _, x := range a {
		i += len(x)
	}
	return i
}

// Returns a sliced like like.
func groupLike(a []string, like [][]string) [][]string {
	var result [][]string
	for _, x := range like {
		result = append(result, a[:len(x)])
		a = a[len(x):]
	}
	return result
}
