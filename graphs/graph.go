// Package graphs implements a simple graph.
package graphs

import (
	"cmp"
	"iter"
	"slices"

	"github.com/fluhus/gostuff/sets"
	"github.com/fluhus/gostuff/snm"
)

// Graph is a simple graph.
// It contains a set of vertices and a set of edges,
// which are pairs of vertices.
type Graph[T comparable] struct {
	v snm.Enumerator[T] // Value to ID.
	e sets.Set[[2]int]  // Pairs of IDs.
}

// New returns an empty graph.
func New[T comparable]() *Graph[T] {
	return &Graph[T]{snm.Enumerator[T]{}, sets.Set[[2]int]{}}
}

// AddVertices adds the given values as vertices.
// Values that already exist are ignored.
func (g *Graph[T]) AddVertices(t ...T) {
	for _, v := range t {
		g.v.IndexOf(v)
	}
}

// NumVertices returns the current number of vertices.
func (g *Graph[T]) NumVertices() int {
	return len(g.v)
}

// NumEdges returns the current number of edges.
func (g *Graph[T]) NumEdges() int {
	return len(g.e)
}

// Edges iterates over current set of edges.
func (g *Graph[T]) Edges() iter.Seq2[T, T] {
	return func(yield func(T, T) bool) {
		flat := g.v.Elements()
		for e := range g.e {
			if !yield(flat[e[0]], flat[e[1]]) {
				return
			}
		}
	}
}

// Vertices iterates over current set of vertices,
// by order of addition to the graph.
func (g *Graph[T]) Vertices() iter.Seq[T] {
	return func(yield func(T) bool) {
		for _, x := range g.v.Elements() {
			if !yield(x) {
				return
			}
		}
	}
}

// HasEdge returns whether there is an edge between a and b.
func (g *Graph[T]) HasEdge(a, b T) bool {
	return g.e.Has(g.toEdge(a, b))
}

// AddEdge adds a and b to the vertex set and adds an edge between them.
// The edge is undirected, meaning that AddEdge(a,b) is equivalent
// to AddEdge(b,a).
// If the edge already exists, this is a no-op.
func (g *Graph[T]) AddEdge(a, b T) {
	g.e.Add(g.toEdge(a, b))
}

// DeleteEdge removes the edge between a and b,
// while keeping them in the vertex set.
func (g *Graph[T]) DeleteEdge(a, b T) {
	delete(g.e, g.toEdge(a, b))
}

// ConnectedComponents returns a slice of connected components.
// In each component, the elements are ordered by order of addition to the
// graph.
// The components are ordered by the order of addition of their
// first elements.
func (g *Graph[T]) ConnectedComponents() [][]T {
	edges := g.edgeSlices()
	m := snm.Slice(g.NumVertices(), func(i int) int { return -1 })
	queue := &snm.Queue[int]{}

	for i := range g.NumVertices() {
		if m[i] != -1 {
			continue
		}
		m[i] = i
		queue.Enqueue(i)
		for queue.Len() > 0 {
			e := queue.Dequeue()
			for _, j := range edges[e] {
				if m[j] == -1 {
					m[j] = i
					queue.Enqueue(j)
				}
			}
			edges[e] = nil
		}
	}

	comps := map[int][]int{}
	for k, v := range m {
		comps[v] = append(comps[v], k)
	}
	poncs := make([][]int, 0, len(comps))
	for _, v := range comps {
		poncs = append(poncs, snm.Sorted(v))
	}
	slices.SortFunc(poncs, func(a, b []int) int {
		return cmp.Compare(a[0], b[0])
	})

	// Convert indices to vertex values.
	i2v := g.v.Elements()
	return snm.Slice(len(poncs), func(i int) []T {
		return snm.Slice(len(poncs[i]), func(j int) T {
			return i2v[poncs[i][j]]
		})
	})
}

// Returns (without adding) an edge between a and b.
func (g *Graph[T]) toEdge(a, b T) [2]int {
	ia, ib := g.v.IndexOf(a), g.v.IndexOf(b)
	if ia > ib {
		return [2]int{ib, ia}
	}
	return [2]int{ia, ib}
}

// Returns a slice representation of this graph's edges.
func (g *Graph[T]) edgeSlices() [][]int {
	// Pre-allocate slices.
	counts := make([]int, g.NumVertices())
	for e := range g.e {
		counts[e[0]]++
		counts[e[1]]++
	}
	edges := snm.Slice(len(counts), func(i int) []int {
		return make([]int, 0, counts[i])
	})

	// Populate with values.
	for e := range g.e {
		edges[e[0]] = append(edges[e[0]], e[1])
		edges[e[1]] = append(edges[e[1]], e[0])
	}
	return edges
}
