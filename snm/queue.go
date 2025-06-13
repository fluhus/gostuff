package snm

import "iter"

// Queue is a memory-efficient FIFO container.
type Queue[T any] struct {
	q    []T
	i, n int
}

// Enqueue inserts an element to the queue.
func (q *Queue[T]) Enqueue(x T) {
	if q.n == len(q.q) {
		nq := make([]T, len(q.q)*2+1)
		_ = append(append(append(nq[:0], q.q[q.i:]...), q.q[:q.i]...), x)
		q.q = nq
		q.n++
		q.i = 0
		return
	}
	i := (q.i + q.n) % len(q.q)
	q.q[i] = x
	q.n++
}

// Dequeue removes and returns the next element in the queue.
// Panics if the queue is empty.
func (q *Queue[T]) Dequeue() T {
	if q.n == 0 {
		panic("pull with 0 elements")
	}
	x := q.q[q.i]
	var zero T
	q.q[q.i] = zero // Remove element to allow GC.
	q.n--
	q.i = (q.i + 1) % len(q.q)
	return x
}

// Peek returns the next element in the queue,
// without modifying its contents.
// Panics if the queue is empty.
func (q *Queue[T]) Peek() T {
	if q.n == 0 {
		panic("pull with 0 elements")
	}
	return q.q[q.i]
}

// Len return the current number of elements in the queue.
func (q *Queue[T]) Len() int {
	return q.n
}

// Seq returns an iterator over the queue's elements,
// dequeueing each one.
//
// It is okay to enqueue elements while iterating,
// from within the same goroutine.
// The new elements will be included in the same loop.
func (q *Queue[T]) Seq() iter.Seq[T] {
	return func(yield func(T) bool) {
		for q.Len() > 0 {
			if !yield(q.Dequeue()) {
				break
			}
		}
	}
}
