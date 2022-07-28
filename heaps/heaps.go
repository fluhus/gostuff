// Package heaps provides generic heaps.
package heaps

import (
	"container/heap"

	"golang.org/x/exp/constraints"
)

// New returns a new heap that uses the given comparator function.
func New[T any](less func(T, T) bool) *Heap[T] {
	return &Heap[T]{backheap[T]{less, nil}}
}

// Min returns a new min-heap of an ordered type by its natural order.
func Min[T constraints.Ordered]() *Heap[T] {
	return New(func(a, b T) bool {
		return a < b
	})
}

// Max returns a new max-heap of an ordered type by its natural order.
func Max[T constraints.Ordered]() *Heap[T] {
	return New(func(a, b T) bool {
		return a > b
	})
}

// backheap is used for communicating with the heap package.
type backheap[T any] struct {
	less func(T, T) bool
	a    []T
}

// Implement heap.Interface.

func (h *backheap[T]) Len() int {
	return len(h.a)
}

func (h *backheap[T]) Less(i, j int) bool {
	return h.less(h.a[i], h.a[j])
}

func (h *backheap[T]) Swap(i, j int) {
	h.a[i], h.a[j] = h.a[j], h.a[i]
}

func (h *backheap[T]) Push(x any) {
	h.a = append(h.a, x.(T))
}

func (h *backheap[T]) Pop() any {
	x := h.a[len(h.a)-1]
	h.a = h.a[:len(h.a)-1]
	// Shrink if needed.
	if cap(h.a) >= 16 && len(h.a) <= cap(h.a)/4 {
		h.a = append(make([]T, 0, cap(h.a)/2), h.a...)
	}
	return x
}

// Heap is a generic heap.
type Heap[T any] struct {
	h backheap[T]
}

// Len returns the number of elements in h.
func (h *Heap[T]) Len() int {
	return h.h.Len()
}

// Push adds x to h while maintaining its heap invariants.
func (h *Heap[T]) Push(x T) {
	heap.Push(&h.h, x)
}

// PushSlice adds the elements of s to h while maintaining its heap invariants.
// The complexity is O(new n), so it should be typically used to initialize a
// new heap.
func (h *Heap[T]) PushSlice(s []T) {
	h.h.a = append(h.h.a, s...)
	heap.Init(&h.h)
}

// Pop removes and returns the minimal element in h.
func (h *Heap[T]) Pop() T {
	return heap.Pop(&h.h).(T)
}

// Head returns the minimal element in h.
func (h *Heap[T]) Head() T {
	return h.h.a[0]
}

// View returns the underlying slice of h, containing all of its elements.
// Modifying the slice may invalidate the heap.
func (h *Heap[T]) View() []T {
	return h.h.a
}
