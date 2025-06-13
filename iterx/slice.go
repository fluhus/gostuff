package iterx

import (
	"iter"
	"slices"
)

// Slice returns an iterator over the slice values.
//
// Deprecated: use [slices.Values] instead.
func Slice[T any](s []T) iter.Seq[T] {
	return slices.Values(s)
}

// ISlice returns an iterator over the slice values and their indices,
// like in a range expression.
//
// Deprecated: use [slices.All] instead.
func ISlice[T any](s []T) iter.Seq2[int, T] {
	return slices.All(s)
}

// Limit returns an iterator that stops after n elements,
// if the underlying iterator does not stop before.
func Limit[T any](it iter.Seq[T], n int) iter.Seq[T] {
	return func(yield func(T) bool) {
		i := 0
		for x := range it {
			i++
			if i > n {
				return
			}
			if !yield(x) {
				return
			}
		}
	}
}

// Limit2 returns an iterator that stops after n elements,
// if the underlying iterator does not stop before.
func Limit2[T any, S any](it iter.Seq2[T, S], n int) iter.Seq2[T, S] {
	return func(yield func(T, S) bool) {
		i := 0
		for x, y := range it {
			i++
			if i > n {
				return
			}
			if !yield(x, y) {
				return
			}
		}
	}
}

// Skip returns an iterator without the first n elements.
func Skip[T any](it iter.Seq[T], n int) iter.Seq[T] {
	return func(yield func(T) bool) {
		i := 0
		for x := range it {
			i++
			if i <= n {
				continue
			}
			if !yield(x) {
				return
			}
		}
	}
}

// Skip2 returns an iterator without the first n elements.
func Skip2[T any, S any](it iter.Seq2[T, S], n int) iter.Seq2[T, S] {
	return func(yield func(T, S) bool) {
		i := 0
		for x, y := range it {
			i++
			if i <= n {
				continue
			}
			if !yield(x, y) {
				return
			}
		}
	}
}
