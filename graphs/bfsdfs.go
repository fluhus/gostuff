package graphs

import (
	"iter"

	"github.com/fluhus/gostuff/sets"
	"github.com/fluhus/gostuff/snm"
)

// BFS iterates over this graph's nodes in a breadth-first ordering,
// including start.
func (g *Graph[T]) BFS(start T) iter.Seq[T] {
	return func(yield func(T) bool) {
		if _, ok := g.v[start]; !ok {
			return
		}

		elems := g.v.Elements()
		edges := g.edgeSlices()
		istart := g.v.IndexOf(start)
		done := sets.Set[int]{}.Add(istart)
		q := &snm.Queue[int]{}
		q.Enqueue(istart)

		for v := range q.Seq() {
			if !yield(elems[v]) {
				return
			}
			for _, e := range edges[v] {
				if done.Has(e) {
					continue
				}
				done.Add(e)
				q.Enqueue(e)
			}
		}
	}
}

// DFS iterates over this graph's nodes in a depth-first ordering,
// including start.
func (g *Graph[T]) DFS(start T) iter.Seq[T] {
	return func(yield func(T) bool) {
		if _, ok := g.v[start]; !ok {
			return
		}

		elems := g.v.Elements()
		edges := g.edgeSlices()
		istart := g.v.IndexOf(start)
		done := sets.Set[int]{}.Add(istart)
		q := []int{istart}

		for len(q) > 0 {
			v := q[len(q)-1]
			q = q[:len(q)-1]
			if !yield(elems[v]) {
				return
			}
			for _, e := range edges[v] {
				if done.Has(e) {
					continue
				}
				done.Add(e)
				q = append(q, e)
			}
		}
	}
}
