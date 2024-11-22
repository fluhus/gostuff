package ppln

import (
	"fmt"
	"iter"
	"sync"
	"sync/atomic"
)

// NonSerial starts a multi-goroutine transformation pipeline.
//
// Input is an iterator over the input values to be transformed.
// It will be called in a thread-safe manner.
// Transform receives an input (a) and a 0-based goroutine number (g),
// and returns the result of processing a.
// Output acts on a single result, and will be called in a thread-safe manner.
// The order of outputs is arbitrary, but correlated with the order of
// inputs.
//
// If one of the functions returns a non-nil error, the process stops and the
// error is returned. Otherwise returns nil.
func NonSerial[T1 any, T2 any](
	ngoroutines int,
	input iter.Seq2[T1, error],
	transform func(a T1, g int) (T2, error),
	output func(a T2) error) error {
	if ngoroutines < 1 {
		panic(fmt.Sprintf("bad number of goroutines: %d", ngoroutines))
	}
	pull, pstop := iter.Pull2(input)
	defer pstop()

	// An optimization for a single thread.
	if ngoroutines == 1 {
		for {
			t1, err, ok := pull()

			if !ok {
				return nil
			}
			if err != nil {
				return err
			}

			t2, err := transform(t1, 0)
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

	for g := 0; g < ngoroutines; g++ {
		go func(g int) {
			for {
				if stop.Load() {
					errs <- nil
					return
				}

				ilock.Lock()
				t1, err, ok := pull()
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

				t2, err := transform(t1, g)
				if err != nil {
					stop.Store(true)
					errs <- err
					return
				}

				olock.Lock()
				err = output(t2)
				olock.Unlock()
				if err != nil {
					stop.Store(true)
					errs <- err
					return
				}
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
