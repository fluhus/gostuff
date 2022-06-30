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
// Pusher should call push on every input value.
// Mapper receives an input (a), 0-based input serial number (i), 0-based
// goroutine number (g), and a Stopper, and returns the result of processing a.
// Puller acts on a single result, and will be called by the same
// order of pusher's input.
func Serial[T1 any, T2 any](
	ngoroutines int,
	pusher func(push func(T1), s Stopper),
	mapper func(a T1, i int, g int, s Stopper) T2,
	puller func(a T2, s Stopper)) {
	if ngoroutines < 1 {
		panic(fmt.Sprintf("bad number of goroutines: %d", ngoroutines))
	}

	stopper := make(Stopper)

	// An optimization for a single thread.
	if ngoroutines == 1 {
		i := 0
		pusher(func(a T1) {
			puller(mapper(a, i, 0, stopper), stopper)
			i++
		}, stopper)
		return
	}

	push := make(chan serialItem[T1], ngoroutines*chanLenCoef)
	pull := make(chan serialItem[T2], ngoroutines*chanLenCoef)
	wait := &sync.WaitGroup{}
	wait.Add(ngoroutines)

	go func() {
		i := 0
		pusher(func(a T1) {
			push <- serialItem[T1]{i, a}
			i++
		}, stopper)
		close(push)
	}()
	for i := 0; i < ngoroutines; i++ {
		i := i
		go func() {
			for item := range push {
				if stopper.Stopped() {
					break
				}
				pull <- serialItem[T2]{
					item.i,
					mapper(item.data, item.i, i, stopper),
				}
			}
			for range push { // Drain channel.
			}
			wait.Done()
		}()
	}
	go func() {
		items := &serialHeap[T2]{
			data: heaps.New[serialItemComparator[T2], serialItem[T2]](),
		}
		for item := range pull {
			if stopper.Stopped() {
				break
			}
			items.put(item)
			for items.ok() {
				puller(items.pop(), stopper)
			}
		}
		for range pull { // Drain channel.
		}
		wait.Done()
	}()

	wait.Wait() // Wait for workers.
	wait.Add(1)
	close(pull)
	wait.Wait() // Wait for pull.
}

// General data with a serial number.
type serialItem[T any] struct {
	i    int
	data T
}

type serialItemComparator[T any] struct{}

func (_ serialItemComparator[T]) Less(a, b serialItem[T]) bool {
	return a.i < b.i
}

// A heap of serial items. Sorts by serial number.
type serialHeap[T any] struct {
	next int
	data *heaps.Heap[serialItemComparator[T], serialItem[T]]
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
