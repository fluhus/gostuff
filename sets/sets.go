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
	if len(s) > len(t) {
		s, t = t, s
	}
	result := Set[T]{}
	for v := range s {
		if t.Has(v) {
			result[v] = struct{}{}
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
