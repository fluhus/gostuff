// Package binio provides high-performance funcitonality for binary IO.
//
// The package is currently thread-unsafe, even when writing to different
// writers.
package binio

// TODO(amit): Make MaxBytesLen a parameter to string functions.
// TODO(amit): Make the package thread-safe.

import (
	"fmt"
	"io"
)

var (
	// TODO(amit): This renders the package thread-unsafe.
	buf = make([]byte, 10)
)

// Uint64ToBytes populates the first 8 bytes in the given slice with the given
// number.
func Uint64ToBytes(x uint64, b []byte) {
	for i := 0; i < 8; i++ {
		b[i] = byte(x & 255)
		x >>= 8
	}
}

// Uint64FromBytes returns the number represented by the given slice.
func Uint64FromBytes(b []byte) uint64 {
	var result uint64
	for i := 0; i < 8; i++ {
		result += uint64(b[i]) << (8 * i)
	}
	return result
}

// WriteUint64 writes a uint64 to the given writer.
func WriteUint64(w io.Writer, x uint64) error {
	Uint64ToBytes(x, buf)
	_, err := w.Write(buf[:8])
	return err
}

// ReadUint64 reads a uint64 from the given reader.
func ReadUint64(r io.Reader) (uint64, error) {
	_, err := io.ReadFull(r, buf[:8])
	if err != nil {
		return 0, err
	}
	return Uint64FromBytes(buf), nil
}

// WriteUvarint writes a varint-encoded uint64 to the given writer.
func WriteUvarint(w io.Writer, x uint64) error {
	if x == 0 {
		return WriteByte(w, 0)
	}
	buf[0] = byte(x & 127)
	x >>= 7
	i := 1
	for x > 0 {
		buf[i-1] |= 128
		buf[i] = byte(x & 127)
		x >>= 7
		i++
	}
	_, err := w.Write(buf[:i])
	return err
}

// ReadUvarint reads a varint-encoded uint64 from the given reader.
func ReadUvarint(r io.Reader) (uint64, error) {
	var result uint64
	for i := 0; ; i++ {
		b, err := ReadByte(r)
		if err != nil {
			return 0, err
		}
		result += uint64(b&127) << (7 * i)
		if b&128 == 0 {
			break
		}
	}
	return result, nil
}

// WriteByte writes a byte to the given writer.
// Saves the need to create and keep a buffer for writing.
func WriteByte(w io.Writer, b byte) error {
	buf[0] = b
	_, err := w.Write(buf[:1])
	return err
}

// ReadByte reads a byte from the given writer.
// Saves the need to create and keep a buffer for reading.
func ReadByte(r io.Reader) (byte, error) {
	n, err := io.ReadFull(r, buf[:1])
	if n != 1 {
		return 0, err
	}
	return buf[0], nil
}

// WriteBytes writes a slice of bytes to the given writer.
// Returns an error if the string is longer than MaxBytesLen.
func WriteBytes(w io.Writer, b []byte) error {
	if err := WriteUvarint(w, uint64(len(b))); err != nil {
		return err
	}
	_, err := w.Write(b)
	return err
}

// ReadBytes reads a slice of bytes from the given writer.
// Returns an error if the string is longer than MaxBytesLen.
func ReadBytes(r io.Reader) ([]byte, error) {
	n, err := ReadUvarint(r)
	if err != nil {
		return nil, err
	}
	b := make([]byte, n)
	_, err = io.ReadFull(r, b)
	if err == io.EOF && n > 0 {
		return nil, io.ErrUnexpectedEOF
	}
	if err != nil {
		return nil, err
	}
	return b, nil
}

// WriteString writes a string to the given writer.
// Returns an error if the string is longer than MaxBytesLen.
func WriteString(w io.Writer, s string) error {
	return WriteBytes(w, []byte(s))
}

// ReadString reads a string from the given writer.
// Returns an error if the string is longer than MaxBytesLen.
func ReadString(r io.Reader) (string, error) {
	b, err := ReadBytes(r)
	return string(b), err
}

// GetBit returns the value of the n'th bit in a byte slice.
func GetBit(b []byte, n int) int {
	return int(b[n/8] >> (n % 8) & 1)
}

// SetBit sets the value of the n'th bit in a byte slice.
func SetBit(b []byte, n, v int) {
	if v == 0 {
		b[n/8] &= ^(byte(1) << (n % 8))
	} else if v == 1 {
		b[n/8] |= byte(1) << (n % 8)
	} else {
		panic(fmt.Sprintf("Bad value: %v, expected 0 or 1", v))
	}
}

// GetHalfByte returns the value of the n'th half byte (4 bits) in a byte slice.
func GetHalfByte(b []byte, n int) byte {
	return (b[n/2] & (0b00001111 << (n % 2 * 4))) >> (n % 2 * 4)
}

// SetHalfByte sets the value of the n'th half byte (4 bits) in a byte slice.
func SetHalfByte(b []byte, n int, v byte) {
	if v > 0b00001111 {
		panic(fmt.Sprintf("Invalid value: %v, should be less than %v",
			v, 0b00001111))
	}
	v ^= GetHalfByte(b, n)
	b[n/2] ^= v << (n % 2 * 4)
}
