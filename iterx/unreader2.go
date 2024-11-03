package iterx

import (
	"fmt"
	"iter"
)

// An Unreader2 wraps an iterator and adds an unread function.
type Unreader2[T, S any] struct {
	read    func() (T, S, bool) // The underlying reader
	stop    func()
	headT   T    // The last read element
	headS   S    // The last read element
	hasHead bool // Unread was called
}

// Read returns the result of calling the underlying reader,
// or the last unread element.
func (r *Unreader2[T, S]) Read() (T, S, bool) {
	if r.hasHead {
		// Remove headX values to allow GC.
		var t T
		var s S
		t, r.headT = r.headT, t
		s, r.headS = r.headS, s
		r.hasHead = false
		return t, s, true
	}
	t, s, ok := r.read()
	r.headT = t
	r.headS = s
	return t, s, ok
}

// Unread makes the next call to Read return the last element and a nil error.
// Can be called up to once per call to Read.
func (r *Unreader2[T, S]) Unread() {
	if r.hasHead {
		panic(fmt.Sprintf("called Unread twice: first with (%v,%v)",
			r.headT, r.headS))
	}
	r.hasHead = true
}

// Until calls Read until stop returns true.
func (r *Unreader2[T, S]) Until(stop func(T, S) bool) iter.Seq2[T, S] {
	return func(yield func(T, S) bool) {
		for {
			t, s, ok := r.Read()
			if !ok {
				return
			}
			if stop(t, s) {
				r.Unread()
				return
			}
			if !yield(t, s) {
				r.stop()
				return
			}
		}
	}
}

// GroupBy returns an iterator over groups, where each group is an iterator that
// yields elements as long as the elements have sameGroup==true with the first
// element in the group.
func (r *Unreader2[T, S]) GroupBy(
	sameGroup func(oldT T, oldS S, newT T, newS S) bool,
) iter.Seq[iter.Seq2[T, S]] {
	return func(yield func(iter.Seq2[T, S]) bool) {
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
func (r *Unreader2[T, S]) nextGroup(sameGroup func(T, S, T, S) bool) (iter.Seq2[T, S], bool) {
	firstT, firstS, ok := r.Read()
	if !ok {
		return nil, false
	}
	return func(yield func(T, S) bool) {
		if !yield(firstT, firstS) {
			return
		}
		for {
			t, s, ok := r.Read()
			if !ok {
				return
			}
			if !sameGroup(firstT, firstS, t, s) {
				r.Unread()
				return
			}
			if !yield(t, s) {
				r.stop()
				return
			}
		}
	}, true
}

// NewUnreader2 returns an Unreader2 with seq as its underlying iterator.
func NewUnreader2[T, S any](seq iter.Seq2[T, S]) *Unreader2[T, S] {
	read, stop := iter.Pull2(seq)
	return &Unreader2[T, S]{read: read, stop: stop, hasHead: false}
}
