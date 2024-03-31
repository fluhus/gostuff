// Prints error rates of Morris counter using different m's.
package main

import (
	"fmt"

	"github.com/fluhus/gostuff/gnum"
	"github.com/fluhus/gostuff/morris"
	"github.com/fluhus/gostuff/ptimer"
)

const (
	upto = 10000000
	reps = 100
)

func main() {
	ms := []uint{1, 3, 10, 30, 100, 300, 1000, 3000, 10000}
	var errs []float64
	for _, m := range ms {
		pt := ptimer.NewMessage(fmt.Sprint("{} (", m, ")"))
		err := 0.0
		for rep := 1; rep <= reps; rep++ {
			c := uint(0)
			for i := 1; i <= upto; i++ {
				c = morris.Raise(c, m)
				r := morris.Restore(c, m)
				err += float64(gnum.Abs(int(i)-int(r))) / float64(i)
				pt.Inc()
			}
		}
		pt.Done()
		errs = append(errs, err/upto/reps)
	}
	for i := range ms {
		fmt.Printf("// % 10d: %.1f%%\n", ms[i], errs[i]*100)
	}
}
