package ppln

import (
	"fmt"
	"sync"

	"github.com/fluhus/gostuff/heaps"
)

const (
	// Size of pipeline channel buffers per goroutine.
	chanLenCoef = 1000
)

// Serial starts a multi-goroutine transformation pipeline that maintains the
// order of the inputs.
//
// Pusher should call push on every input value. Stop indicates if an error was
// returned and pushing should stop.
// Mapper receives an input (a), 0-based input serial number (i), 0-based
// goroutine number (g), and returns the result of processing a.
// Puller acts on a single result, and will be called by the same
// order of pusher's input.
//
// If one of the functions returns a non-nil error, the process stops and the
// error is returned. Otherwise returns nil.
func Serial[T1 any, T2 any](
	ngoroutines int,
	pusher func(push func(T1), stop func() bool) error,
	mapper func(a T1, i int, g int) (T2, error),
	puller func(a T2) error) error {
	if ngoroutines < 1 {
		panic(fmt.Sprintf("bad number of goroutines: %d", ngoroutines))
	}

	var err error

	// An optimization for a single thread.
	if ngoroutines == 1 {
		i := 0
		pusher(func(a T1) {
			if err != nil {
				return
			}
			t2, merr := mapper(a, i, 0)
			if merr != nil {
				err = merr
				return
			}
			perr := puller(t2)
			if perr != nil {
				err = perr
				return
			}
			i++
		}, func() bool { return err != nil })
		return err
	}

	push := make(chan serialItem[T1], ngoroutines*chanLenCoef)
	pull := make(chan serialItem[T2], ngoroutines*chanLenCoef)
	wait := &sync.WaitGroup{}
	wait.Add(ngoroutines)

	go func() {
		i := 0
		perr := pusher(func(a T1) {
			if err != nil {
				return
			}
			push <- serialItem[T1]{i, a}
			i++
		}, func() bool { return err != nil })
		if perr != nil && err == nil {
			err = perr
		}
		close(push)
	}()
	for i := 0; i < ngoroutines; i++ {
		i := i
		go func() {
			for item := range push {
				if err != nil {
					continue // Drain channel.
				}
				t2, merr := mapper(item.data, item.i, i)
				if merr != nil {
					if err == nil {
						err = merr
					}
					continue
				}
				pull <- serialItem[T2]{item.i, t2}
			}
			wait.Done()
		}()
	}
	go func() {
		items := &serialHeap[T2]{
			data: heaps.New(func(a, b serialItem[T2]) bool {
				return a.i < b.i
			}),
		}
		for item := range pull {
			if err != nil {
				continue // Drain channel.
			}
			items.put(item)
			for items.ok() {
				perr := puller(items.pop())
				if perr != nil && err == nil {
					err = perr
				}
			}
		}
		wait.Done()
	}()

	wait.Wait() // Wait for workers.
	wait.Add(1)
	close(pull)
	wait.Wait() // Wait for pull.

	return err
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
