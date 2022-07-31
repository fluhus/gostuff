// Package ppln provides generic parallel processing pipelines.
//
// NOTE: this API is currently experimental and may change in future releases.
//
// General usage
//
// This package provides two modes of operation: serial and non-serial.
// Serial transforms each value of type T1 to a value of type T2. The outputs
// are ordered in the same order of the inputs. Non-serial transforms each value
// of type T1 to zero or more values of type T2. The order of the outputs is
// arbitrary, but correlated with the order of inputs.
//
// Each of the functions blocks the calling function until either the processing
// is done (puller was called on the last value) or until stopped.
//
// Stopping
//
// Each user-function (pusher, mapper, puller) may return an error.
// Returning a non-nil error stops the pipeline prematurely, and that
// error will be propagated to the caller.
//
// Number of goroutines
//
// Each pipeline creates ngoroutines+2 new goroutines and blocks the calling
// one. There is one pusher goroutine, one puller goroutine, and ngoroutines
// mapper goroutines.
//
// A special case is when ngoroutines==1, in which case no new goroutines are
// created. Processing is done serially using the calling goroutine.
package ppln
