// Package heaps provides generic heaps.
package heaps

import (
	"container/heap"

	"golang.org/x/exp/constraints"
)

// New returns a new heap that uses the given comparator.
func New[C Comparator[T], T any]() *Heap[C, T] {
	var c C
	return &Heap[C, T]{backheap[C, T]{c, nil}}
}

// NewMin returns a new min-heap of an ordered type by its natural order.
func NewMin[T constraints.Ordered]() *Heap[Min[T], T] {
	return New[Min[T], T]()
}

// NewMax returns a new max-heap of an ordered type by its natural order.
func NewMax[T constraints.Ordered]() *Heap[Max[T], T] {
	return New[Max[T], T]()
}

// Comparator provides the comparison function for a heap.
// Less should determine whether a should be popped before b.
type Comparator[T any] interface {
	Less(a, b T) bool
}

// Min is a comparator that places the minimal value at the top of the heap.
type Min[T constraints.Ordered] struct{}

func (_ Min[T]) Less(a, b T) bool {
	return a < b
}

// Max is a comparator that places the maximal value at the top of the heap.
type Max[T constraints.Ordered] struct{}

func (_ Max[T]) Less(a, b T) bool {
	return a > b
}

// backheap is used for communicating with the heap package.
type backheap[C Comparator[T], T any] struct {
	c C
	a []T
}

// Implement heap.Interface.

func (h *backheap[C, T]) Len() int {
	return len(h.a)
}

func (h *backheap[C, T]) Less(i, j int) bool {
	return h.c.Less(h.a[i], h.a[j])
}

func (h *backheap[C, T]) Swap(i, j int) {
	h.a[i], h.a[j] = h.a[j], h.a[i]
}

func (h *backheap[C, T]) Push(x any) {
	h.a = append(h.a, x.(T))
}

func (h *backheap[C, T]) Pop() any {
	x := h.a[len(h.a)-1]
	h.a = h.a[:len(h.a)-1]
	// Shrink if needed.
	if cap(h.a) >= 16 && len(h.a) <= cap(h.a)/4 {
		h.a = append(make([]T, 0, cap(h.a)/2), h.a...)
	}
	return x
}

// Heap is a generic heap.
type Heap[C Comparator[T], T any] struct {
	h backheap[C, T]
}

// Len returns the number of elements in h.
func (h *Heap[C, T]) Len() int {
	return h.h.Len()
}

// Push adds x to h while maintaining its heap invariants.
func (h *Heap[C, T]) Push(x T) {
	heap.Push(&h.h, x)
}

// Pop removes and returns the minimal element in h.
func (h *Heap[C, T]) Pop() T {
	return heap.Pop(&h.h).(T)
}

// Head returns the minimal element in h.
func (h *Heap[C, T]) Head() T {
	return h.h.a[0]
}

// View returns the underlying slice of h, containing all of its elements.
// Modifying the slice may invalidate the heap.
func (h *Heap[C, T]) View() []T {
	return h.h.a
}
