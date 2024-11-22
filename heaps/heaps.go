// Package heaps provides generic heaps.
//
// This package provides better run speeds than the standard [heap] package.
package heaps

import (
	"golang.org/x/exp/constraints"
)

// Heap is a generic heap.
type Heap[T any] struct {
	a    []T
	less func(T, T) bool
}

// New returns a new heap that uses the given comparator function.
func New[T any](less func(T, T) bool) *Heap[T] {
	return &Heap[T]{nil, less}
}

// Min returns a new min-heap of an ordered type by its natural order.
func Min[T constraints.Ordered]() *Heap[T] {
	return New(func(t1, t2 T) bool {
		return t1 < t2
	})
}

// Max returns a new max-heap of an ordered type by its natural order.
func Max[T constraints.Ordered]() *Heap[T] {
	return New(func(t1, t2 T) bool {
		return t1 > t2
	})
}

// Push adds x to h while maintaining its heap invariants.
func (h *Heap[T]) Push(x T) {
	h.a = append(h.a, x)
	i := len(h.a) - 1
	for i != -1 {
		i = h.bubbleUp(i)
	}
}

// PushSlice adds the elements of s to h while maintaining its heap invariants.
// The complexity is O(new n), so it should be typically used to initialize a
// new heap.
func (h *Heap[T]) PushSlice(s []T) {
	h.a = append(h.a, s...)
	for i := len(h.a) - 1; i >= 0; i-- {
		j := i
		for j != -1 {
			j = h.bubbleDown(j)
		}
	}
}

// Pop removes and returns the minimal element in h.
func (h *Heap[T]) Pop() T {
	if len(h.a) == 0 {
		panic("called Pop() on an empty heap")
	}
	x := h.a[0]
	h.a[0] = h.a[len(h.a)-1]
	h.a = h.a[:len(h.a)-1]
	i := 0
	for i != -1 {
		i = h.bubbleDown(i)
	}
	// Shrink if needed.
	if cap(h.a) >= 16 && len(h.a) <= cap(h.a)/4 {
		h.a = append(make([]T, 0, cap(h.a)/2), h.a...)
	}
	return x
}

// Len returns the number of elements in h.
func (h *Heap[T]) Len() int {
	return len(h.a)
}

// Moves the i'th element down and returns its new index.
// Returns -1 when no more bubble-downs are needed.
func (h *Heap[T]) bubbleDown(i int) int {
	ia, ib := i*2+1, i*2+2
	if len(h.a) < ib { // No children
		return -1
	}
	if len(h.a) == ib { // Only one child
		if h.less(h.a[ia], h.a[i]) {
			h.a[i], h.a[ia] = h.a[ia], h.a[i]
		}
		return -1
	}
	if h.less(h.a[ib], h.a[ia]) {
		ia = ib
	}
	if h.less(h.a[ia], h.a[i]) {
		h.a[i], h.a[ia] = h.a[ia], h.a[i]
	}
	return ia
}

// Moves the i'th element up and returns its new index.
// Returns -1 when no more bubble-ups are needed.
func (h *Heap[T]) bubbleUp(i int) int {
	if i == 0 {
		return -1
	}
	ia := (i - 1) / 2
	if !h.less(h.a[i], h.a[ia]) {
		return -1
	}
	h.a[i], h.a[ia] = h.a[ia], h.a[i]
	return ia
}

// View returns the underlying slice of h, containing all of its elements.
// Modifying the slice may invalidate the heap.
func (h *Heap[T]) View() []T {
	return h.a
}

// Head returns the minimal element in h.
func (h *Heap[T]) Head() T {
	return h.a[0]
}

// Fix fixes the heap after a single value had been modified.
// i is the index of the modified value.
func (h *Heap[T]) Fix(i int) {
	wentUp := false
	for {
		j := h.bubbleUp(i)
		if j == -1 {
			break
		}
		i = j
		wentUp = true
	}
	if !wentUp {
		for {
			j := h.bubbleDown(i)
			if j == -1 {
				break
			}
			i = j
		}
	}
}

// Clip removes unused capacity from the heap.
func (h *Heap[T]) Clip() {
	h.a = append(make([]T, 0, len(h.a)), h.a...)
}
