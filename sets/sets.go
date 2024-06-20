// Package sets provides generic sets.
package sets

import (
	"encoding/json"

	"golang.org/x/exp/maps"
)

// Set is a convenience wrapper around map[T]struct{}.
type Set[T comparable] map[T]struct{}

// Add inserts the given elements to s and returns s.
func (s Set[T]) Add(t ...T) Set[T] {
	for _, v := range t {
		s[v] = struct{}{}
	}
	return s
}

// AddSet inserts the elements of t to s and returns s.
func (s Set[T]) AddSet(t Set[T]) Set[T] {
	for v := range t {
		s[v] = struct{}{}
	}
	return s
}

// Remove deletes the given elements from s and returns s.
func (s Set[T]) Remove(t ...T) Set[T] {
	for _, v := range t {
		delete(s, v)
	}
	return s
}

// RemoveSet deletes the elements of t from s and returns s.
func (s Set[T]) RemoveSet(t Set[T]) Set[T] {
	for v := range t {
		delete(s, v)
	}
	return s
}

// Has returns whether t is a member of s.
func (s Set[T]) Has(t T) bool {
	_, ok := s[t]
	return ok
}

// Intersect returns a new set holding the elements that are common
// to s and t.
func (s Set[T]) Intersect(t Set[T]) Set[T] {
	if len(s) > len(t) { // Iterate over the smaller one.
		s, t = t, s
	}
	result := Set[T]{}
	for v := range s {
		if t.Has(v) {
			result.Add(v)
		}
	}
	return result
}

// MarshalJSON implements the json.Marshaler interface.
func (s Set[T]) MarshalJSON() ([]byte, error) {
	return json.Marshal(maps.Keys(s))
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (s *Set[T]) UnmarshalJSON(b []byte) error {
	var slice []T
	if err := json.Unmarshal(b, &slice); err != nil {
		return err
	}
	if *s == nil {
		*s = Set[T]{}
	}
	s.Add(slice...)
	return nil
}

// AddKeys adds the keys of a map to a set.
func AddKeys[K comparable, V any](s Set[K], m map[K]V) Set[K] {
	for k := range m {
		s.Add(k)
	}
	return s
}

// AddValues adds the values of a map to a set.
func AddValues[K comparable, V comparable](s Set[V], m map[K]V) Set[V] {
	for _, v := range m {
		s.Add(v)
	}
	return s
}

// Of returns a new set containing the given elements.
func Of[T comparable](t ...T) Set[T] {
	return Set[T]{}.Add(t...)
}

// FromKeys returns a new set containing the keys of the given map.
func FromKeys[K comparable, V any](m map[K]V) Set[K] {
	return AddKeys(make(Set[K], len(m)), m)
}

// FromValues returns a new set containing the values of the given map.
func FromValues[K comparable, V comparable](m map[K]V) Set[V] {
	return AddValues(Set[V]{}, m)
}
