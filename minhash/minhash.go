// Package minhash provides a min-hash collection for approximating Jaccard
// similarity.
package minhash

import (
	"encoding/json"
	"fmt"

	"github.com/fluhus/gostuff/gnum"
	"github.com/fluhus/gostuff/heaps"
	"github.com/fluhus/gostuff/sets"
	"golang.org/x/exp/constraints"
	"golang.org/x/exp/slices"
)

// A MinHash is a min-hash collection. Retains the k lowest unique values out of all
// the values that were added to it.
type MinHash[T constraints.Integer] struct {
	h *heaps.Heap[T] // Min-hash heap
	s sets.Set[T]    // Keeps elements unique
	k int            // Max size of the collection
	n int            // Number of calls to Push
}

// New returns an empty collection that stores k values.
func New[T constraints.Integer](k int) *MinHash[T] {
	if k < 1 {
		panic(fmt.Sprintf("invalid n: %d, should be positive", k))
	}
	return &MinHash[T]{
		heaps.Max[T](),
		sets.Set[T]{},
		k, 0,
	}
}

// Push tries to add a hash to the collection. x is added only if it does not
// already exist, and there are less than k elements lesser than x.
// Returns true if x was added and false if not.
func (mh *MinHash[T]) Push(x T) bool {
	mh.n++
	if mh.h.Len() == mh.k && x >= mh.h.Head() {
		// x is too large.
		return false
	}
	if mh.s.Has(x) {
		return false
	}
	if mh.h.Len() == mh.k {
		mh.s.Remove(mh.h.Pop())
	}
	mh.h.Push(x)
	mh.s.Add(x)
	return true
}

// K returns the maximal number of elements in mh.
func (mh *MinHash[T]) K() int {
	return mh.k
}

// N returns the number of calls that were made to Push.
// Represents the size of the original set.
func (mh *MinHash[T]) N() int {
	return mh.n
}

// View returns the underlying slice of values.
func (mh *MinHash[T]) View() []T {
	return mh.h.View()
}

// MarshalJSON implements the json.Marshaler interface.
func (mh *MinHash[T]) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		K int `json:"k"`
		N int `json:"n"`
		H []T `json:"h"`
	}{mh.k, mh.n, mh.View()})
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (mh *MinHash[T]) UnmarshalJSON(b []byte) error {
	var raw struct {
		K int
		N int
		H []T
	}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	ss := New[T](raw.K)
	for _, v := range raw.H {
		ss.Push(v)
	}
	ss.n = raw.N
	*mh = *ss
	return nil
}

// Returns the intersection and union sizes of mh and other,
// in min-hash terms.
func (mh *MinHash[T]) intersect(other *MinHash[T]) (int, int) {
	a, b := mh.View(), other.View()
	if !slices.IsSortedFunc(a, func(x, y T) bool {
		return x > y
	}) {
		panic("receiver is not sorted")
	}
	if !slices.IsSortedFunc(b, func(x, y T) bool {
		return x > y
	}) {
		panic("other is not sorted")
	}
	intersection := 0
	i, j, m := len(a)-1, len(b)-1, 0
	for ; i >= 0 && j >= 0 && m < mh.k; m++ {
		if a[i] > b[j] {
			j--
		} else if a[i] < b[j] {
			i--
		} else { // a[i] == b[j]
			intersection++
			i--
			j--
		}
	}
	union := gnum.Min2(mh.k, m+len(a)-i+len(b)-j)
	return intersection, union
}

// Jaccard returns the approximated Jaccard similarity between mh and other.
//
// Sort needs to be called before calling this function.
func (mh *MinHash[T]) Jaccard(other *MinHash[T]) float64 {
	i, u := mh.intersect(other)
	return float64(i) / float64(u)
}

// SoftJaccard returns the Jaccard similarity between mh and other,
// adding one agreed upon element and one disagreed upon element to
// the calculation.
//
// Sort needs to be called before calling this function.
func (mh *MinHash[T]) SoftJaccard(other *MinHash[T]) float64 {
	r := mh.Jaccard(other)
	sum := float64(mh.N() + other.N())
	ri, ru := r*sum/(r+1), sum/(r+1)
	return (ri + 1) / (ru + 2)
}

// Sort sorts the collection, making it ready for Jaccard calculation.
// The collection is still valid after calling Sort.
func (mh *MinHash[T]) Sort() {
	slices.SortFunc(mh.View(), func(x, y T) bool {
		return x > y
	})
}
