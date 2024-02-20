// Package ppln provides generic parallel processing pipelines.
//
// # General usage
//
// This package provides two modes of operation: serial and non-serial.
// In [Serial] the outputs are ordered in the same order of the inputs.
// In [NonSerial] the order of outputs is arbitrary,
// but correlated with the order of inputs.
//
// Each of the functions blocks the calling function until either the processing
// is done (output was called on the last value) or until an error is returned.
//
// # Stopping
//
// Each user-function (input, transform, output) may return an error.
// Returning a non-nil error stops the pipeline prematurely, and that
// error is returned to the caller.
//
// # Experimental
//
// This package relies on the experimental [iter] package.
// In order to use it, go 1.22 is required with GOEXPERIMENT=rangefunc.
package ppln
