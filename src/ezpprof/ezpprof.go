// Convinience wrapper for the runtime/pprof package.
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

// Stops CPU profiling and closes the output file. If called without calling
// Start, does nothing.
func Stop() {
	if fout == nil {
		return
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
