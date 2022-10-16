// Package ezpprof is a convenience wrapper over the runtime/pprof package.
//
// This package helps to quickly introduce profiling to a piece of code without
// the mess of opening files and checking errors.
//
// A typical use of this package looks like:
//
//	ezpprof.Start("myfile.pprof")
//	{... some complicated code ...}
//	ezpprof.Stop()
//
// Or alternatively:
//
//	const profile = true
//
//	if profile {
//	  ezpprof.Start("myfile.pprof")
//	  defer ezpprof.Stop()
//	}
package ezpprof

import (
	"io"
	"runtime/pprof"

	"github.com/fluhus/gostuff/aio"
)

var fout io.WriteCloser

// Start starts CPU profiling and writes to the given file.
// Panics if an error occurs.
func Start(file string) {
	if fout != nil {
		panic("already profiling")
	}
	f, err := aio.CreateRaw(file)
	if err != nil {
		panic(err)
	}
	fout = f
	pprof.StartCPUProfile(fout)
}

// Stop stops CPU profiling and closes the output file.
// Panics if called without calling Start.
func Stop() {
	if fout == nil {
		panic("Stop called without calling Start")
	}
	pprof.StopCPUProfile()
	if err := fout.Close(); err != nil {
		panic(err)
	}
	fout = nil
}

// Heap writes heap profile to the given file. Panics if an error occurs.
func Heap(file string) {
	f, err := aio.Create(file)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	err = pprof.WriteHeapProfile(f)
	if err != nil {
		panic(err)
	}
}
