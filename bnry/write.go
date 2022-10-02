package bnry

import (
	"encoding/binary"
	"fmt"
	"io"
	"math"
	"reflect"

	"golang.org/x/exp/constraints"
)

// Write writes the given values to the given writer.
// Values should be of any of the supported types.
// Returns an error if a value is of an unsupported type.
func Write(w io.Writer, vals ...any) error {
	buf, err := MarshalBinary(vals...)
	if err != nil {
		return err
	}
	_, err = w.Write(buf)
	return err
}

// MarshalBinary writes the given values to a byte slice.
// Values should be of any of the supported types.
// Returns an error if a value is of an unsupported type.
func MarshalBinary(vals ...any) ([]byte, error) {
	var buf []byte
	var err error
	for i, val := range vals {
		buf, err = appendSingle(buf, val)
		if err != nil {
			return nil, fmt.Errorf("argument #%d: %w", i+1, err)
		}
	}
	return buf, nil
}

// Appends a single value as binary to buf.
func appendSingle(buf []byte, val any) ([]byte, error) {
	switch val := val.(type) {
	case uint8:
		return append(buf, val), nil
	case uint16:
		return binary.AppendUvarint(buf, uint64(val)), nil
	case uint32:
		return binary.AppendUvarint(buf, uint64(val)), nil
	case uint64:
		return binary.AppendUvarint(buf, uint64(val)), nil
	case int8:
		return append(buf, byte(val)), nil
	case int16:
		return binary.AppendVarint(buf, int64(val)), nil
	case int32:
		return binary.AppendVarint(buf, int64(val)), nil
	case int64:
		return binary.AppendVarint(buf, int64(val)), nil
	case float32:
		return binary.AppendUvarint(buf, uint64(math.Float32bits(val))), nil
	case float64:
		return binary.AppendUvarint(buf, math.Float64bits(val)), nil
	case bool:
		return append(buf, boolToByte(val)), nil
	case string:
		return appendString(buf, val), nil
	case []uint8:
		return appendUint8Slice(buf, val), nil
	case []uint16:
		return appendUintSlice(buf, val), nil
	case []uint32:
		return appendUintSlice(buf, val), nil
	case []uint64:
		return appendUintSlice(buf, val), nil
	case []int8:
		return appendInt8Slice(buf, val), nil
	case []int16:
		return appendIntSlice(buf, val), nil
	case []int32:
		return appendIntSlice(buf, val), nil
	case []int64:
		return appendIntSlice(buf, val), nil
	case []float32:
		return appendFloat32Slice(buf, val), nil
	case []float64:
		return appendFloat64Slice(buf, val), nil
	case []bool:
		return appendBoolSlice(buf, val), nil
	case []string:
		return appendStringSlice(buf, val), nil
	default:
		return nil, fmt.Errorf("unsupported type: %v",
			reflect.TypeOf(val).Name())
	}
}

func appendUint8Slice(buf []byte, s []uint8) []byte {
	buf = binary.AppendUvarint(buf, uint64(len(s)))
	return append(buf, s...)
}

func appendString(buf []byte, s string) []byte {
	buf = binary.AppendUvarint(buf, uint64(len(s)))
	return append(buf, s...)
}

func appendInt8Slice(buf []byte, s []int8) []byte {
	buf = binary.AppendUvarint(buf, uint64(len(s)))
	for _, x := range s {
		buf = append(buf, byte(x))
	}
	return buf
}

func appendUintSlice[T constraints.Unsigned](buf []byte, s []T) []byte {
	buf = binary.AppendUvarint(buf, uint64(len(s)))
	for _, x := range s {
		buf = binary.AppendUvarint(buf, uint64(x))
	}
	return buf
}

func appendIntSlice[T constraints.Signed](buf []byte, s []T) []byte {
	buf = binary.AppendUvarint(buf, uint64(len(s)))
	for _, x := range s {
		buf = binary.AppendVarint(buf, int64(x))
	}
	return buf
}

func appendFloat32Slice(buf []byte, s []float32) []byte {
	buf = binary.AppendUvarint(buf, uint64(len(s)))
	for _, x := range s {
		buf = binary.AppendUvarint(buf, uint64(math.Float32bits(x)))
	}
	return buf
}

func appendFloat64Slice(buf []byte, s []float64) []byte {
	buf = binary.AppendUvarint(buf, uint64(len(s)))
	for _, x := range s {
		buf = binary.AppendUvarint(buf, math.Float64bits(x))
	}
	return buf
}

func appendBoolSlice(buf []byte, s []bool) []byte {
	buf = binary.AppendUvarint(buf, uint64(len(s)))
	for _, x := range s {
		buf = append(buf, boolToByte(x))
	}
	return buf
}

func appendStringSlice(buf []byte, s []string) []byte {
	buf = binary.AppendUvarint(buf, uint64(len(s)))
	for _, x := range s {
		buf = appendString(buf, x)
	}
	return buf
}

func boolToByte(b bool) byte {
	if b {
		return 1
	} else {
		return 0
	}
}
