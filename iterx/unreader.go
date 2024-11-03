package iterx

import (
	"fmt"
	"iter"
)

// An Unreader wraps an iterator and adds an unread function.
type Unreader[T any] struct {
	read    func() (T, bool) // The underlying reader
	stop    func()
	head    T    // The last read element
	hasHead bool // Unread was called
}

// Read returns the result of calling the underlying reader,
// or the last unread element.
func (r *Unreader[T]) Read() (T, bool) {
	if r.hasHead {
		// Remove headX values to allow GC.
		var t T
		t, r.head = r.head, t
		r.hasHead = false
		return t, true
	}
	t, ok := r.read()
	r.head = t
	return t, ok
}

// Unread makes the next call to Read return the last element and a nil error.
// Can be called up to once per call to Read.
func (r *Unreader[T]) Unread() {
	if r.hasHead {
		panic(fmt.Sprintf("called Unread twice: first with %v", r.head))
	}
	r.hasHead = true
}

// Until calls Read until stop returns true.
func (r *Unreader[T]) Until(stop func(T) bool) iter.Seq[T] {
	return func(yield func(T) bool) {
		for {
			t, ok := r.Read()
			if !ok {
				return
			}
			if stop(t) {
				r.Unread()
				return
			}
			if !yield(t) {
				r.stop()
				return
			}
		}
	}
}

// GroupBy returns an iterator over groups, where each group is an iterator that
// yields elements as long as the elements have sameGroup==true with the first
// element in the group.
func (r *Unreader[T]) GroupBy(sameGroup func(old T, nu T) bool,
) iter.Seq[iter.Seq[T]] {
	return func(yield func(iter.Seq[T]) bool) {
		for {
			group, ok := r.nextGroup(sameGroup)
			if !ok {
				return
			}
			if !yield(group) {
				r.stop()
				return
			}
		}
	}
}

// Returns an iterator that yields elements while sameGroup==true.
func (r *Unreader[T]) nextGroup(sameGroup func(T, T) bool) (iter.Seq[T], bool) {
	firstT, ok := r.Read()
	if !ok {
		return nil, false
	}
	return func(yield func(T) bool) {
		if !yield(firstT) {
			return
		}
		for {
			t, ok := r.Read()
			if !ok {
				return
			}
			if !sameGroup(firstT, t) {
				r.Unread()
				return
			}
			if !yield(t) {
				r.stop()
				return
			}
		}
	}, true
}

// NewUnreader returns an Unreader with seq as its underlying iterator.
func NewUnreader[T any](seq iter.Seq[T]) *Unreader[T] {
	read, stop := iter.Pull(seq)
	return &Unreader[T]{read: read, stop: stop, hasHead: false}
}
