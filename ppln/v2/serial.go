package ppln

import (
	"fmt"
	"iter"
	"sync"
	"sync/atomic"

	"github.com/fluhus/gostuff/heaps"
)

// Serial starts a multi-goroutine transformation pipeline that maintains the
// order of the inputs.
//
// Input is an iterator over the input values to be transformed.
// It will be called in a thread-safe manner.
// Transform receives an input (a), 0-based input serial number (i), 0-based
// goroutine number (g), and returns the result of processing a.
// Output acts on a single result, and will be called by the same
// order of the input, in a thread-safe manner.
//
// If one of the functions returns a non-nil error, the process stops and the
// error is returned. Otherwise returns nil.
func Serial[T1 any, T2 any](
	ngoroutines int,
	input iter.Seq2[T1, error],
	transform func(a T1, i int, g int) (T2, error),
	output func(a T2) error) error {
	if ngoroutines < 1 {
		panic(fmt.Sprintf("bad number of goroutines: %d", ngoroutines))
	}
	pull, pstop := iter.Pull2(input)
	defer pstop()

	// An optimization for a single thread.
	if ngoroutines == 1 {
		i := 0
		for {
			t1, err, ok := pull()
			ii := i
			i++

			if !ok {
				return nil
			}
			if err != nil {
				return err
			}

			t2, err := transform(t1, ii, 0)
			if err != nil {
				return err
			}
			if err := output(t2); err != nil {
				return err
			}
		}
	}

	ilock := &sync.Mutex{}
	olock := &sync.Mutex{}
	errs := make(chan error, ngoroutines)
	stop := &atomic.Bool{}
	items := &serialHeap[T2]{
		data: heaps.New(func(a, b serialItem[T2]) bool {
			return a.i < b.i
		}),
	}

	i := 0
	for g := 0; g < ngoroutines; g++ {
		go func(g int) {
			for {
				if stop.Load() {
					errs <- nil
					return
				}

				ilock.Lock()
				t1, err, ok := pull()
				ii := i
				i++
				ilock.Unlock()

				if !ok {
					errs <- nil
					return
				}
				if err != nil {
					stop.Store(true)
					errs <- err
					return
				}

				t2, err := transform(t1, ii, g)
				if err != nil {
					stop.Store(true)
					errs <- err
					return
				}

				olock.Lock()
				items.put(serialItem[T2]{ii, t2})
				for items.ok() {
					err = output(items.pop())
					if err != nil {
						olock.Unlock()
						stop.Store(true)
						errs <- err
						return
					}
				}
				olock.Unlock()
			}
		}(g)
	}

	for g := 0; g < ngoroutines; g++ {
		if err := <-errs; err != nil {
			return err
		}
	}
	return nil
}

// General data with a serial number.
type serialItem[T any] struct {
	i    int
	data T
}

// A heap of serial items. Sorts by serial number.
type serialHeap[T any] struct {
	next int
	data *heaps.Heap[serialItem[T]]
}

// Checks whether the minimal element in the heap is the next in the series.
func (s *serialHeap[T]) ok() bool {
	return s.data.Len() > 0 && s.data.Head().i == s.next
}

// Removes and returns the minimal element in the heap. Panics if the element
// is not the next in the series.
func (s *serialHeap[T]) pop() T {
	if !s.ok() {
		panic("get when not ok")
	}
	s.next++
	a := s.data.Pop()
	return a.data
}

// Adds an item to the heap.
func (s *serialHeap[T]) put(item serialItem[T]) {
	if item.i < s.next {
		panic(fmt.Sprintf("put(%d) when next is %d", item.i, s.next))
	}
	s.data.Push(item)
}
