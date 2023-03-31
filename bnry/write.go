package bnry

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"math"
	"reflect"
	"strings"

	"golang.org/x/exp/constraints"
)

// Write writes the given values to the given writer.
// Values should be of any of the supported types.
// Panics if a value is of an unsupported type.
func Write(w io.Writer, vals ...any) error {
	return NewWriter(w).Write(vals...)
}

// MarshalBinary writes the given values to a byte slice.
// Values should be of any of the supported types.
// Panics if a value is of an unsupported type.
func MarshalBinary(vals ...any) []byte {
	buf := bytes.NewBuffer(nil)
	NewWriter(buf).Write(vals...)
	return buf.Bytes()
}

type Writer struct {
	w   io.Writer
	buf [binary.MaxVarintLen64]byte
}

func NewWriter(w io.Writer) *Writer {
	return &Writer{w: w}
}

// Write writes the given values.
// Values should be of any of the supported types.
// Panics if a value is of an unsupported type.
func (w *Writer) Write(vals ...any) error {
	for _, v := range vals {
		if err := w.writeSingle(v); err != nil {
			return err
		}
	}
	return nil
}

// Writes a single value as binary.
func (w *Writer) writeSingle(val any) error {
	switch val := val.(type) {
	case uint8:
		return w.writeByte(val)
	case uint16:
		return writeUint(w, val)
	case uint32:
		return writeUint(w, val)
	case uint64:
		return writeUint(w, val)
	case uint:
		return writeUint(w, val)
	case int8:
		return w.writeByte(byte(val))
	case int16:
		return writeInt(w, val)
	case int32:
		return writeInt(w, val)
	case int64:
		return writeInt(w, val)
	case int:
		return writeInt(w, val)
	case float32:
		return writeUint(w, math.Float32bits(val))
	case float64:
		return writeUint(w, math.Float64bits(val))
	case bool:
		return w.writeByte(boolToByte(val))
	case string:
		return w.writeString(val)
	case []uint8:
		return w.writeUint8Slice(val)
	case []uint16:
		return writeUintSlice(w, val)
	case []uint32:
		return writeUintSlice(w, val)
	case []uint64:
		return writeUintSlice(w, val)
	case []uint:
		return writeUintSlice(w, val)
	case []int8:
		return w.writeInt8Slice(val)
	case []int16:
		return writeIntSlice(w, val)
	case []int32:
		return writeIntSlice(w, val)
	case []int64:
		return writeIntSlice(w, val)
	case []int:
		return writeIntSlice(w, val)
	case []float32:
		return w.writeFloat32Slice(val)
	case []float64:
		return w.writeFloat64Slice(val)
	case []bool:
		return w.writeBoolSlice(val)
	case []string:
		return w.writeStringSlice(val)
	default:
		panic(fmt.Sprintf("unsupported type: %v",
			reflect.TypeOf(val).Name()))
	}
}

func (w *Writer) writeByte(b byte) error {
	w.buf[0] = b
	_, err := w.w.Write(w.buf[:1])
	return err
}

func writeUint[T constraints.Unsigned](w *Writer, i T) error {
	_, err := w.w.Write(binary.AppendUvarint(w.buf[:0], uint64(i)))
	return err
}

func writeInt[T constraints.Signed](w *Writer, i T) error {
	_, err := w.w.Write(binary.AppendVarint(w.buf[:0], int64(i)))
	return err
}

func (w *Writer) writeUint8Slice(s []uint8) error {
	if err := writeUint(w, uint(len(s))); err != nil {
		return err
	}
	_, err := w.w.Write(s)
	return err
}

func (w *Writer) writeString(s string) error {
	if err := writeUint(w, uint(len(s))); err != nil {
		return err
	}
	_, err := strings.NewReader(s).WriteTo(w.w)
	return err
}

func (w *Writer) writeInt8Slice(s []int8) error {
	if err := writeUint(w, uint(len(s))); err != nil {
		return err
	}
	for _, x := range s {
		if err := w.writeByte(byte(x)); err != nil {
			return err
		}
	}
	return nil
}

func writeUintSlice[T constraints.Unsigned](w *Writer, s []T) error {
	if err := writeUint(w, uint(len(s))); err != nil {
		return err
	}
	for _, x := range s {
		if err := writeUint(w, x); err != nil {
			return err
		}
	}
	return nil
}

func writeIntSlice[T constraints.Signed](w *Writer, s []T) error {
	if err := writeUint(w, uint(len(s))); err != nil {
		return err
	}
	for _, x := range s {
		if err := writeInt(w, x); err != nil {
			return err
		}
	}
	return nil
}

func (w *Writer) writeFloat32Slice(s []float32) error {
	if err := writeUint(w, uint(len(s))); err != nil {
		return err
	}
	for _, x := range s {
		if err := writeUint(w, math.Float32bits(x)); err != nil {
			return err
		}
	}
	return nil
}

func (w *Writer) writeFloat64Slice(s []float64) error {
	if err := writeUint(w, uint(len(s))); err != nil {
		return err
	}
	for _, x := range s {
		if err := writeUint(w, math.Float64bits(x)); err != nil {
			return err
		}
	}
	return nil
}

func (w *Writer) writeBoolSlice(s []bool) error {
	if err := writeUint(w, uint(len(s))); err != nil {
		return err
	}
	for _, x := range s {
		if err := w.writeByte(boolToByte(x)); err != nil {
			return err
		}
	}
	return nil
}

func (w *Writer) writeStringSlice(s []string) error {
	if err := writeUint(w, uint(len(s))); err != nil {
		return err
	}
	for _, x := range s {
		if err := w.writeString(x); err != nil {
			return err
		}
	}
	return nil
}

func boolToByte(b bool) byte {
	if b {
		return 1
	} else {
		return 0
	}
}
