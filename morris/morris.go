// Package morris provides an implementation of Morris's algorithm
// for approximate counting with few bits.
//
// The original formula raises a counter i with probability 2^(-i).
// The restored value is 2^i - 1.
//
// This package introduces a parameter m, so that the first m increments are
// made with probability 1, then m increments with probability 1/2, then m with
// probability 1/4... Using m=1 is equivalent to the original formula.
// A large m increases accuracy but costs more bits.
// A single counter should use the same m for all calls to Raise and Restore.
//
// This package is experimental.
package morris

import (
	"fmt"
	"math/rand"

	"golang.org/x/exp/constraints"
)

// If true, panics when a counter is about to be raised beyond its maximal
// value.
const checkOverFlow = true

// Raise returns the new value of i after one increment.
// m controls the restoration accuracy.
// The approximate number of calls to Raise can be restored using Restore.
func Raise[T constraints.Unsigned](i T, m uint) T {
	if checkOverFlow {
		max := T(0) - 1
		if i == max {
			panic(fmt.Sprintf("counter reached maximal value: %d", i))
		}
	}
	r := 1 << (uint(i) / m)
	if rand.Intn(r) == 0 {
		return i + 1
	}
	return i
}

// Restore returns an approximation of the number of calls to Raise on i.
// m should have the same value that was used with Raise.
func Restore[T constraints.Unsigned](i T, m uint) uint {
	ui := uint(i)
	if ui <= m {
		return ui
	}
	return m*(1<<(ui/m)-1) + (ui%m)*(1<<(ui/m))
}
