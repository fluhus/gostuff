// Package bnry provides simple functions for encoding and decoding values as
// binary.
//
// # Supported data types
//
// The types that can be encoded and decoded are
// int*, uint* (excluding int and uint), float*, bool, string and
// slices of these types.
// [Read] and [UnmarshalBinary] expect pointers to these types,
// while [Write] and [MarshalBinary] expect non-pointers.
package bnry
