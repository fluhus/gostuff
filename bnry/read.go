package bnry

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"math"
	"reflect"
	"slices"
	"strings"

	"golang.org/x/exp/constraints"
)

// UnmarshalBinary decodes binary data into the given values.
// Values should be pointers to any of the supported types.
// Panics if a value is of an unsupported type.
func UnmarshalBinary(data []byte, vals ...any) error {
	return Read(bytes.NewBuffer(data), vals...)
}

// Read reads and decodes binary data into the given values.
// Values should be pointers to any of the supported types.
// Panics if a value is of an unsupported type.
func Read(r io.ByteReader, vals ...any) error {
	for i, val := range vals {
		if err := readSingle(r, val); err != nil {
			if i > 0 {
				err = notExpectingEOF(err)
			}
			if err != io.EOF {
				err = fmt.Errorf("reading value #%d: %w", i+1, err)
			}
			return err
		}
	}
	return nil
}

// Decodes a single value from r.
func readSingle(r io.ByteReader, val any) error {
	switch val := val.(type) {
	case *uint8:
		return readUint8(r, val)
	case *uint16:
		return readUvarint(r, val)
	case *uint32:
		return readUvarint(r, val)
	case *uint64:
		return readUvarint(r, val)
	case *uint:
		return readUvarint(r, val)
	case *int8:
		return readInt8(r, val)
	case *int16:
		return readVarint(r, val)
	case *int32:
		return readVarint(r, val)
	case *int64:
		return readVarint(r, val)
	case *int:
		return readVarint(r, val)
	case *float32:
		return readFloat32(r, val)
	case *float64:
		return readFloat64(r, val)
	case *bool:
		return readBool(r, val)
	case *string:
		return readString(r, val)
	case *[]uint8:
		return readUint8Slice(r, val)
	case *[]uint16:
		return readUintSlice(r, val)
	case *[]uint32:
		return readUintSlice(r, val)
	case *[]uint64:
		return readUintSlice(r, val)
	case *[]uint:
		return readUintSlice(r, val)
	case *[]int8:
		return readInt8Slice(r, val)
	case *[]int16:
		return readIntSlice(r, val)
	case *[]int32:
		return readIntSlice(r, val)
	case *[]int64:
		return readIntSlice(r, val)
	case *[]int:
		return readIntSlice(r, val)
	case *[]float32:
		return readFloat32Slice(r, val)
	case *[]float64:
		return readFloat64Slice(r, val)
	case *[]bool:
		return readBoolSlice(r, val)
	case *[]string:
		return readStringSlice(r, val)
	default:
		panic(fmt.Sprintf("unsupported type: %v", reflect.TypeOf(val).Name()))
	}
}

func readUint8(r io.ByteReader, val *uint8) error {
	x, err := r.ReadByte()
	*val = x
	return err
}

func readInt8(r io.ByteReader, val *int8) error {
	x, err := r.ReadByte()
	*val = int8(x)
	return err
}

func readUvarint[T constraints.Unsigned](r io.ByteReader, val *T) error {
	x, err := binary.ReadUvarint(r)
	*val = T(x)
	return err
}

func readVarint[T constraints.Signed](r io.ByteReader, val *T) error {
	x, err := binary.ReadVarint(r)
	*val = T(x)
	return err
}

func readFloat32(r io.ByteReader, val *float32) error {
	x, err := binary.ReadUvarint(r)
	*val = math.Float32frombits(uint32(x))
	return err
}

func readFloat64(r io.ByteReader, val *float64) error {
	x, err := binary.ReadUvarint(r)
	*val = math.Float64frombits(x)
	return err
}

func readBool(r io.ByteReader, val *bool) error {
	b, err := r.ReadByte()
	if err != nil {
		return err
	}
	switch b {
	case 0:
		*val = false
	case 1:
		*val = true
	default:
		return fmt.Errorf("unexpected value for bool: %v, want 0 or 1", b)
	}
	return nil
}

func readString(r io.ByteReader, s *string) error {
	n, err := binary.ReadUvarint(r)
	if err != nil {
		return err
	}
	bld := &strings.Builder{}
	bld.Grow(int(n))
	for range n {
		b, err := r.ReadByte()
		if err != nil {
			return notExpectingEOF(err)
		}
		bld.WriteByte(b)
	}
	*s = bld.String()
	return nil
}

func readUint8Slice(r io.ByteReader, val *[]uint8) error {
	n, err := binary.ReadUvarint(r)
	if err != nil {
		return err
	}
	buf := slices.Grow(*val, int(n))[:0]
	for range n {
		b, err := r.ReadByte()
		if err != nil {
			return notExpectingEOF(err)
		}
		buf = append(buf, b)
	}
	*val = buf
	return nil
}

func readInt8Slice(r io.ByteReader, val *[]int8) error {
	n, err := binary.ReadUvarint(r)
	if err != nil {
		return err
	}
	buf := slices.Grow(*val, int(n))[:0]
	for range n {
		b, err := r.ReadByte()
		if err != nil {
			return notExpectingEOF(err)
		}
		buf = append(buf, int8(b))
	}
	*val = buf
	return nil
}

func readUintSlice[T constraints.Unsigned](r io.ByteReader, val *[]T) error {
	n, err := binary.ReadUvarint(r)
	if err != nil {
		return err
	}
	buf := slices.Grow(*val, int(n))[:0]
	for range n {
		x, err := binary.ReadUvarint(r)
		if err != nil {
			return notExpectingEOF(err)
		}
		buf = append(buf, T(x))
	}
	*val = buf
	return nil
}

func readIntSlice[T constraints.Signed](r io.ByteReader, val *[]T) error {
	n, err := binary.ReadUvarint(r)
	if err != nil {
		return err
	}
	buf := slices.Grow(*val, int(n))[:0]
	for range n {
		x, err := binary.ReadVarint(r)
		if err != nil {
			return notExpectingEOF(err)
		}
		buf = append(buf, T(x))
	}
	*val = buf
	return nil
}

func readFloat32Slice(r io.ByteReader, val *[]float32) error {
	n, err := binary.ReadUvarint(r)
	if err != nil {
		return err
	}
	buf := slices.Grow(*val, int(n))[:0]
	for range n {
		var x float32
		if err := readFloat32(r, &x); err != nil {
			return notExpectingEOF(err)
		}
		buf = append(buf, x)
	}
	*val = buf
	return nil
}

func readFloat64Slice(r io.ByteReader, val *[]float64) error {
	n, err := binary.ReadUvarint(r)
	if err != nil {
		return err
	}
	buf := slices.Grow(*val, int(n))[:0]
	for range n {
		var x float64
		if err := readFloat64(r, &x); err != nil {
			return notExpectingEOF(err)
		}
		buf = append(buf, x)
	}
	*val = buf
	return nil
}

func readBoolSlice(r io.ByteReader, val *[]bool) error {
	n, err := binary.ReadUvarint(r)
	if err != nil {
		return err
	}
	buf := slices.Grow(*val, int(n))[:0]
	for range n {
		var x bool
		if err := readBool(r, &x); err != nil {
			return notExpectingEOF(err)
		}
		buf = append(buf, x)
	}
	*val = buf
	return nil
}

func readStringSlice(r io.ByteReader, val *[]string) error {
	n, err := binary.ReadUvarint(r)
	if err != nil {
		return err
	}
	buf := slices.Grow(*val, int(n))[:0]
	for range n {
		var x string
		if err := readString(r, &x); err != nil {
			return notExpectingEOF(err)
		}
		buf = append(buf, x)
	}
	*val = buf
	return nil
}

func notExpectingEOF(err error) error {
	if err == io.EOF {
		return io.ErrUnexpectedEOF
	}
	return err
}
