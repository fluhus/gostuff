// Convinience wrapper for the runtime/pprof package.
//
// Sometimes you just want to profile a piece of code, without the mess of
// opening files and checking errors. This package will help you profile your
// code while keeping it clean.
package ezpprof

import (
	"runtime/pprof"
	"os"
	"bufio"
)

var fout *os.File
var bout *bufio.Writer

// Starts CPU profiling and writes to the given file. Panics if an error occurs.
func Start(file string) {
	if fout != nil {
		panic("Already profiling.")
	}

	var err error
	fout, err = os.Create(file)
	if err != nil {
		fout, bout = nil, nil
		panic(err)
	}
	
	bout = bufio.NewWriter(fout)
	pprof.StartCPUProfile(bout)
}

// Stops CPU profiling and closes the output file. Panics if called without
// calling Start.
func Stop() {
	if fout == nil {
		panic("Stop called without calling Start.")
	}

	pprof.StopCPUProfile()
	bout.Flush()
	fout.Close()
	fout, bout = nil, nil
}

// Writes heap profile to the given file. Panics if an error occurs.
func Heap(file string) {
	f, err := os.Create(file)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	
	err = pprof.WriteHeapProfile(f)
	if err != nil {
		panic(err)
	}
}
