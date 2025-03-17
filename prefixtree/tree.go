// Package prefixtree provides a basic prefix tree implementation.
//
// Add, Has, HasPrefix, Delete and DeletePrefix are linear in query length,
// regardless of the size of the tree.
// All operations are non-recursive.
package prefixtree

import (
	"iter"
	"slices"

	"golang.org/x/exp/maps"
)

// A Tree is a prefix tree.
//
// A zero value tree is invalid; use New to create a new instance.
type Tree struct {
	isElem bool // Is this node an element in the tree.
	m      map[byte]*Tree
}

// New returns an empty tree.
func New() *Tree {
	return &Tree{m: map[byte]*Tree{}}
}

// Add inserts x to the tree.
// If a was already added, the tree is unchanged.
func (t *Tree) Add(x []byte) {
	cur := t
	for _, b := range x {
		next := cur.m[b]
		if next == nil {
			next = New()
			cur.m[b] = next
		}
		cur = next
	}
	cur.isElem = true
}

// Has returns whether x was added to the tree.
func (t *Tree) Has(x []byte) bool {
	cur := t
	for _, b := range x {
		next := cur.m[b]
		if next == nil {
			return false
		}
		cur = next
	}
	return cur.isElem
}

// HasPrefix returns whether x is a prefix of an element in the tree.
func (t *Tree) HasPrefix(x []byte) bool {
	cur := t
	for _, b := range x {
		next := cur.m[b]
		if next == nil {
			return false
		}
		cur = next
	}
	return true
}

// FindPrefixes returns all the elements in the tree
// that are prefixes of x, ordered by length.
func (t *Tree) FindPrefixes(x []byte) [][]byte {
	cur := t
	var result [][]byte
	if cur.isElem {
		result = append(result, x[:0])
	}
	for i, b := range x {
		cur = cur.m[b]
		if cur == nil {
			break
		}
		if cur.isElem {
			result = append(result, x[:i+1])
		}
	}
	return result
}

// Delete removes x from the tree, if possible.
// Returns the result of Has(a) before deletion.
func (t *Tree) Delete(x []byte) bool {
	// Delve in and create a stack.
	stack := make([]*Tree, len(x))
	cur := t
	for i := range x {
		stack[i] = cur
		cur = cur.m[x[i]]
		if cur == nil {
			return false
		}
	}
	if !cur.isElem {
		return false
	}
	cur.isElem = false
	if len(cur.m) > 0 {
		return true
	}

	// Go back and delete nodes.
	for i := len(stack) - 1; i >= 0; i-- {
		delete(stack[i].m, x[i])
		if len(stack[i].m) > 0 {
			// Stop deleting if node has other children.
			break
		}
	}
	return true
}

// DeletePrefix removes prefix x from the tree.
// All sequences that have x as their prefix are removed.
// Other sequences are unchanged.
func (t *Tree) DeletePrefix(x []byte) {
	// Length 0 mean all strings.
	if len(x) == 0 {
		clear(t.m)
		return
	}

	// Delve in and create a stack.
	stack := make([]*Tree, len(x))
	cur := t
	for i := range x {
		stack[i] = cur
		cur = cur.m[x[i]]
		if cur == nil {
			return
		}
	}

	// Go back and delete nodes.
	for i := len(stack) - 1; i >= 0; i-- {
		delete(stack[i].m, x[i])
		if len(stack[i].m) > 0 {
			// Stop deleting if node has other children.
			break
		}
	}
}

// Iter iterates over the elements of t,
// in no particular order.
// The tree should not be modified during iteration.
func (t *Tree) Iter() iter.Seq[[]byte] {
	return func(yield func([]byte) bool) {
		stack := []*iterStep{{t, maps.Keys(t.m)}}
		var cur []byte
		if t.isElem && !yield(cur) {
			return
		}
		for {
			step := stack[len(stack)-1]
			if len(step.k) == 0 { // Finished with this branch.
				stack = stack[:len(stack)-1]
				if len(stack) == 0 { // Done.
					break
				}
				cur = cur[:len(cur)-1]
				continue
			}
			// Handle next child.
			key := step.k[0]
			child := step.t.m[key]
			stack = append(stack, &iterStep{child, maps.Keys(child.m)})
			step.k = step.k[1:]
			cur = append(cur, key)
			if child.isElem && !yield(slices.Clone(cur)) {
				return
			}
		}
	}
}

// IterPrefix iterates the elements in the tree that have x as their prefix.
// Calling IterPrefix with a zero-length slice is equivalent to Iter().
func (t *Tree) IterPrefix(x []byte) iter.Seq[[]byte] {
	return func(yield func([]byte) bool) {
		cur := t
		for i := range x {
			cur = cur.m[x[i]]
			if cur == nil {
				return
			}
		}
		for y := range cur.Iter() {
			if !yield(slices.Concat(x, y)) {
				return
			}
		}
	}
}

// A step in the iteration stack.
type iterStep struct {
	t *Tree  // Tree to process
	k []byte // Tree's remaining children
}
