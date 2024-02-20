package ppln

import (
	"iter"
)

// SliceInput returns a function that iterates over a slice,
// to be used as the input function in [Serial] and [NonSerial].
func SliceInput[T any](s []T) iter.Seq2[T, error] {
	return func(yield func(T, error) bool) {
		for _, t := range s {
			if !yield(t, nil) {
				return
			}
		}
	}
}

// RangeInput returns a function that iterates over a range of integers,
// starting at start and ending at (and excluding) stop,
// to be used as the input function in [Serial] and [NonSerial].
func RangeInput(start, stop int) iter.Seq2[int, error] {
	return func(yield func(int, error) bool) {
		for i := start; i < stop; i++ {
			if !yield(i, nil) {
				return
			}
		}
	}
}
