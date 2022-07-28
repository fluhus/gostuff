package ppln

import (
	"fmt"
	"sync"
)

// NonSerial starts a multi-goroutine transformation pipeline.
//
// Pusher should call push on every input value.
// Stop indicates if an error was returned and pushing should stop.
// Mapper receives an input (a), a push function for the results, 0-based
// goroutine number (g).
// It should call push zero or more times with the processing results of a.
// Puller acts on a single output. The order of the outputs is arbitrary, but
// correlated with the order of pusher's inputs.
//
// If one of the functions returns a non-nil error, the process stops and the
// error is returned. Otherwise returns nil.
func NonSerial[T1 any, T2 any](
	ngoroutines int,
	pusher func(push func(T1), stop func() bool) error,
	mapper func(a T1, push func(T2), g int) error,
	puller func(a T2) error) error {
	if ngoroutines < 1 {
		panic(fmt.Sprintf("bad number of goroutines: %d", ngoroutines))
	}

	var err error

	// An optimization for a single thread.
	if ngoroutines == 1 {
		perr := pusher(func(a T1) {
			if err != nil {
				return
			}
			merr := mapper(a, func(i T2) {
				if err != nil {
					return
				}
				perr := puller(i)
				if perr != nil && err == nil {
					err = perr
				}
			}, 0)
			if merr != nil && err == nil {
				err = merr
			}
		}, func() bool { return err != nil })
		if perr != nil && err == nil {
			err = perr
		}
		return err
	}

	push := make(chan T1, ngoroutines*chanLenCoef)
	pull := make(chan T2, ngoroutines*chanLenCoef)
	wait := &sync.WaitGroup{}
	wait.Add(ngoroutines)

	go func() {
		perr := pusher(func(a T1) {
			if err != nil {
				return
			}
			push <- a
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
				merr := mapper(item, func(a T2) {
					if err != nil {
						return
					}
					pull <- a
				}, i)
				if merr != nil && err == nil {
					err = merr
				}
			}
			wait.Done()
		}()
	}
	go func() {
		for item := range pull {
			if err != nil {
				continue // Drain channel.
			}
			perr := puller(item)
			if perr != nil && err == nil {
				err = perr
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

	return err
}
