// Package bits provides operations on bit arrays.
package bits

import (
	"iter"
	mbits "math/bits"

	"golang.org/x/exp/constraints"
)

// Set sets the n'th bit to 1 or 0 for values of true or false
// respectively.
func Set[I constraints.Integer](data []byte, n I, value bool) {
	if value {
		Set1(data, n)
	} else {
		Set0(data, n)
	}
}

// Set1 sets the n'th bit in data to 1.
func Set1[I constraints.Integer](data []byte, n I) {
	data[n/8] |= 1 << (n % 8)
}

// Set0 sets the n'th bit in data to 0.
func Set0[I constraints.Integer](data []byte, n I) {
	data[n/8] &= ^(1 << (n % 8))
}

// Get returns the value of the n'th bit (0 or 1).
func Get[I constraints.Integer](data []byte, n I) int {
	return int((data[n/8] >> (n % 8)) & 1)
}

// Sum returns the number of bits that have a value of 1.
func Sum(data []byte) int {
	a := 0
	for _, b := range data {
		a += mbits.OnesCount8(b)
	}
	return a
}

// Ones iterates over the indexes of bits whose values are 1.
func Ones(data []byte) iter.Seq[int] {
	return func(yield func(int) bool) {
		for i, x := range data {
			for _, b := range byteOnes[x] {
				if !yield(i*8 + b) {
					return
				}
			}
		}
	}
}

// Zeros iterates over the indexes of bits whose values are 0.
func Zeros(data []byte) iter.Seq[int] {
	return func(yield func(int) bool) {
		for i, x := range data {
			for _, b := range byteZeros[x] {
				if !yield(i*8 + b) {
					return
				}
			}
		}
	}
}

// Indexes of ones for each byte value.
var byteOnes [][]int

// Indexes of zeros for each byte value.
var byteZeros [][]int

// Calculates indexes of ones and zeros for each byte value.
func init() {
	byteOnes = make([][]int, 256)
	byteZeros = make([][]int, 256)
	for i := range 256 {
		var ones, zeros []int
		for b := range 8 {
			if (i>>b)&1 == 1 {
				ones = append(ones, b)
			} else {
				zeros = append(zeros, b)
			}
		}
		byteOnes[i] = ones
		byteZeros[i] = zeros
	}
}
