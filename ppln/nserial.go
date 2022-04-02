package ppln

import (
	"fmt"
	"sync"
)

// NonSerial starts a multi-goroutine transformation pipeline.
//
// Pusher should call push on every input value.
// Mapper receives an input (a), a push function for the results, 0-based
// goroutine number (g), and a Stopper. It should call push zero or more
// times with the processing results of a.
// Puller acts on a single output. The order of the outputs is arbitrary, but
// correlated with the order of pusher's inputs.
func NonSerial[T1 any, T2 any](
	ngoroutines int,
	pusher func(push func(T1), s Stopper),
	mapper func(a T1, push func(T2), g int, s Stopper),
	puller func(a T2, s Stopper)) {
	if ngoroutines < 1 {
		panic(fmt.Sprintf("bad number of goroutines: %d", ngoroutines))
	}

	stopper := make(Stopper)

	// An optimization for a single thread.
	if ngoroutines == 1 {
		pusher(func(a T1) {
			mapper(a, func(i T2) {
				puller(i, stopper)
			}, 0, stopper)
		}, stopper)
		return
	}

	push := make(chan T1, ngoroutines*chanLenCoef)
	pull := make(chan T2, ngoroutines*chanLenCoef)
	wait := &sync.WaitGroup{}
	wait.Add(ngoroutines)

	go func() {
		pusher(func(a T1) {
			push <- a
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
				mapper(item, func(a T2) {
					pull <- a
				}, i, stopper)
			}
			for range push { // Drain channel.
			}
			wait.Done()
		}()
	}
	go func() {
		for item := range pull {
			if stopper.Stopped() {
				break
			}
			puller(item, stopper)
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
