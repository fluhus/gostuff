// Package snm provides convenience functions for slices and maps.
package snm

import (
	"cmp"
	"slices"

	"golang.org/x/exp/constraints"
	"golang.org/x/exp/maps"
)

// Slice returns a new slice of size n whose values are the results
// of applying f on each index.
func Slice[T any](n int, f func(int) T) []T {
	s := make([]T, n)
	for i := range s {
		s[i] = f(i)
	}
	return s
}

// SliceToSlice returns a slice of the same length containing the results
// of applying f to the elements of s.
func SliceToSlice[A any, B any](a []A, f func(A) B) []B {
	b := make([]B, len(a))
	for i := range a {
		b[i] = f(a[i])
	}
	return b
}

// MapToMap returns a map containing the results of applying f to the key-value
// pairs of m.
// f should return a new key-value pair for the new map.
// Keys that appear more than once will override each other.
func MapToMap[K comparable, V any, K2 comparable, V2 any](
	m map[K]V, f func(K, V) (K2, V2)) map[K2]V2 {
	mm := make(map[K2]V2, len(m))
	for k, v := range m {
		k2, v2 := f(k, v)
		mm[k2] = v2
	}
	return mm
}

// FilterSlice returns a new slice containing only the elements
// for which keep returns true.
func FilterSlice[S any](s []S, keep func(S) bool) []S {
	var result []S
	for _, e := range s {
		if keep(e) {
			result = append(result, e)
		}
	}
	return result
}

// FilterMap returns a new map containing only the elements
// for which keep returns true.
func FilterMap[K comparable, V any](m map[K]V, keep func(k K, v V) bool) map[K]V {
	mm := map[K]V{}
	for k, v := range m {
		if keep(k, v) {
			mm[k] = v
		}
	}
	return mm
}

// Sorted sorts the input and returns it.
func Sorted[T constraints.Ordered](s []T) []T {
	slices.Sort(s)
	return s
}

// SortedFunc sorts the input and returns it.
func SortedFunc[T any](s []T, cmp func(T, T) int) []T {
	slices.SortFunc(s, cmp)
	return s
}

// At returns the elements of t at the indexes in at.
func At[T any, I constraints.Integer](t []T, at []I) []T {
	result := make([]T, 0, len(at))
	for _, i := range at {
		result = append(result, t[i])
	}
	return result
}

// DefaultMap wraps a map with a function that generates values for missing keys.
type DefaultMap[K comparable, V any] struct {
	M map[K]V   // Underlying map. Can be safely read from and written to.
	F func(K) V // Generator function.
}

// Get returns the value associated with key k.
// If k is missing from the map, the generator function is called with k and the
// result becomes k's value.
func (m DefaultMap[K, V]) Get(k K) V {
	if v, ok := m.M[k]; ok {
		return v
	}
	v := m.F(k)
	m.M[k] = v
	return v
}

// Set sets v as k's value.
func (m DefaultMap[K, V]) Set(k K, v V) {
	m.M[k] = v
}

// NewDefaultMap returns an empty map with the given function as the missing
// value generator.
func NewDefaultMap[K comparable, V any](f func(K) V) DefaultMap[K, V] {
	return DefaultMap[K, V]{map[K]V{}, f}
}

// Compare is a generic comparator function for ordered types.
//
// Deprecated: use [cmp.Compare] instead.
func Compare[T constraints.Ordered](a, b T) int {
	if a < b {
		return -1
	}
	if a > b {
		return 1
	}
	return 0
}

// CompareReverse orders values from big to small.
// Should be generally used as a parameter, not called.
func CompareReverse[T constraints.Ordered](a, b T) int {
	return -1 * cmp.Compare(a, b)
}

// SortedKeys sorts a map's keys according to their values' natural order.
func SortedKeys[K comparable, V constraints.Ordered](
	m map[K]V) []K {
	return SortedFunc(maps.Keys(m), func(a, b K) int {
		return Compare(m[a], m[b])
	})
}

// SortedKeysFunc sorts a map's keys by comparing their values.
func SortedKeysFunc[K comparable, V constraints.Ordered](
	m map[K]V, cmp func(V, V) int) []K {
	return SortedFunc(maps.Keys(m), func(a, b K) int {
		return cmp(m[a], m[b])
	})
}

// Number is an integer or a float.
type Number interface {
	constraints.Integer | constraints.Float
}

// Cast casts each element in the slice.
func Cast[TO Number, FROM Number](s []FROM) []TO {
	return SliceToSlice(s, func(x FROM) TO { return TO(x) })
}

// CapMap is a wrapper over a regular map,
// for reusing maps in order to reduce garbage generation.
// The allocated memory for the map is expanded and shrunk as needed.
type CapMap[K comparable, V any] struct {
	m map[K]V // Underlying map
	c int     // Capacity
}

func NewCapMap[K comparable, V any]() *CapMap[K, V] {
	return &CapMap[K, V]{map[K]V{}, 0}
}

// Map returns the underlying map for regular use.
func (s *CapMap[K, V]) Map() map[K]V {
	return s.m
}

// Clear clears the contents of this map, reducing its
// capacity if needed.
// May change which object is returned by Map.
func (s *CapMap[K, V]) Clear() {
	s.c = max(s.c, len(s.m))
	if s.c > 64 && len(s.m) <= s.c/3 {
		newCap := s.c / 2
		s.m = make(map[K]V, newCap)
		s.c = newCap
	} else {
		clear(s.m)
	}
}

// Enumerator enumerates values by their order of appearance.
type Enumerator[T comparable] map[T]int

// IndexOf returns the index of t, possibly allocating a new one.
//
// Equal values always have the same index, while different values
// always have different indexes.
// Indexes are sequential.
func (e Enumerator[T]) IndexOf(t T) int {
	if i, ok := e[t]; ok {
		return i
	}
	i := len(e)
	e[t] = i
	return i
}

// Elements returns the enumerated elements, by order of appearance.
func (e Enumerator[T]) Elements() []T {
	result := make([]T, len(e))
	for t, i := range e {
		result[i] = t
	}
	return result
}
