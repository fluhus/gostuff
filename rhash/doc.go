// Package rhash provides implementations of rolling-hash functions.
//
// A rolling-hash is a hash function that "remembers" only the last n bytes it
// received, where n is a parameter. Meaning, the hash of a byte sequence
// always equals the hash of its last n bytes.
package rhash
