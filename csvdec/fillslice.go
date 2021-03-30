// ***** DO NOT EDIT THIS FILE MANUALLY. *****
//
// This file was auto-generated from 'fillslice.got' using Gent.
//
// Gent: https://www.github.com/fluhus/gent

package csvdec

import (
	"reflect"
	"strconv"
)

// Populates any slice value.
func fillSlice(value reflect.Value, fields []string) error {
	kind := value.Type().Elem().Kind()
	switch kind {
	case reflect.String:
		return fillStringSlice(value, fields)
	case reflect.Int:
		return fillIntSlice(value, fields)
	case reflect.Int8:
		return fillInt8Slice(value, fields)
	case reflect.Int16:
		return fillInt16Slice(value, fields)
	case reflect.Int32:
		return fillInt32Slice(value, fields)
	case reflect.Int64:
		return fillInt64Slice(value, fields)
	case reflect.Uint:
		return fillUintSlice(value, fields)
	case reflect.Uint8:
		return fillUint8Slice(value, fields)
	case reflect.Uint16:
		return fillUint16Slice(value, fields)
	case reflect.Uint32:
		return fillUint32Slice(value, fields)
	case reflect.Uint64:
		return fillUint64Slice(value, fields)
	case reflect.Float32:
		return fillFloat32Slice(value, fields)
	case reflect.Float64:
		return fillFloat64Slice(value, fields)
	}
	panic("Unsupported type: " + value.Type().String())
}

// Populates the given int slice with values parsed from fields.
// Returns an error if parsing fails.
func fillIntSlice(value reflect.Value, fields []string) error {
	parsed := make([]int, len(fields))
	for i, field := range fields {
		n, err := strconv.ParseInt(field, 0, 0)
		if err != nil {
			return err
		}
		parsed[i] = int(n)
	}
	value.Set(reflect.ValueOf(parsed))
	return nil
}

// Populates the given int8 slice with values parsed from fields.
// Returns an error if parsing fails.
func fillInt8Slice(value reflect.Value, fields []string) error {
	parsed := make([]int8, len(fields))
	for i, field := range fields {
		n, err := strconv.ParseInt(field, 0, 8)
		if err != nil {
			return err
		}
		parsed[i] = int8(n)
	}
	value.Set(reflect.ValueOf(parsed))
	return nil
}

// Populates the given int16 slice with values parsed from fields.
// Returns an error if parsing fails.
func fillInt16Slice(value reflect.Value, fields []string) error {
	parsed := make([]int16, len(fields))
	for i, field := range fields {
		n, err := strconv.ParseInt(field, 0, 16)
		if err != nil {
			return err
		}
		parsed[i] = int16(n)
	}
	value.Set(reflect.ValueOf(parsed))
	return nil
}

// Populates the given int32 slice with values parsed from fields.
// Returns an error if parsing fails.
func fillInt32Slice(value reflect.Value, fields []string) error {
	parsed := make([]int32, len(fields))
	for i, field := range fields {
		n, err := strconv.ParseInt(field, 0, 32)
		if err != nil {
			return err
		}
		parsed[i] = int32(n)
	}
	value.Set(reflect.ValueOf(parsed))
	return nil
}

// Populates the given int64 slice with values parsed from fields.
// Returns an error if parsing fails.
func fillInt64Slice(value reflect.Value, fields []string) error {
	parsed := make([]int64, len(fields))
	for i, field := range fields {
		n, err := strconv.ParseInt(field, 0, 64)
		if err != nil {
			return err
		}
		parsed[i] = int64(n)
	}
	value.Set(reflect.ValueOf(parsed))
	return nil
}

// Populates the given uint slice with values parsed from fields.
// Returns an error if parsing fails.
func fillUintSlice(value reflect.Value, fields []string) error {
	parsed := make([]uint, len(fields))
	for i, field := range fields {
		n, err := strconv.ParseUint(field, 0, 0)
		if err != nil {
			return err
		}
		parsed[i] = uint(n)
	}
	value.Set(reflect.ValueOf(parsed))
	return nil
}

// Populates the given uint8 slice with values parsed from fields.
// Returns an error if parsing fails.
func fillUint8Slice(value reflect.Value, fields []string) error {
	parsed := make([]uint8, len(fields))
	for i, field := range fields {
		n, err := strconv.ParseUint(field, 0, 8)
		if err != nil {
			return err
		}
		parsed[i] = uint8(n)
	}
	value.Set(reflect.ValueOf(parsed))
	return nil
}

// Populates the given uint16 slice with values parsed from fields.
// Returns an error if parsing fails.
func fillUint16Slice(value reflect.Value, fields []string) error {
	parsed := make([]uint16, len(fields))
	for i, field := range fields {
		n, err := strconv.ParseUint(field, 0, 16)
		if err != nil {
			return err
		}
		parsed[i] = uint16(n)
	}
	value.Set(reflect.ValueOf(parsed))
	return nil
}

// Populates the given uint32 slice with values parsed from fields.
// Returns an error if parsing fails.
func fillUint32Slice(value reflect.Value, fields []string) error {
	parsed := make([]uint32, len(fields))
	for i, field := range fields {
		n, err := strconv.ParseUint(field, 0, 32)
		if err != nil {
			return err
		}
		parsed[i] = uint32(n)
	}
	value.Set(reflect.ValueOf(parsed))
	return nil
}

// Populates the given uint64 slice with values parsed from fields.
// Returns an error if parsing fails.
func fillUint64Slice(value reflect.Value, fields []string) error {
	parsed := make([]uint64, len(fields))
	for i, field := range fields {
		n, err := strconv.ParseUint(field, 0, 64)
		if err != nil {
			return err
		}
		parsed[i] = uint64(n)
	}
	value.Set(reflect.ValueOf(parsed))
	return nil
}

// Populates the given float32 slice with values parsed from fields.
// Returns an error if parsing fails.
func fillFloat32Slice(value reflect.Value, fields []string) error {
	parsed := make([]float32, len(fields))
	for i, field := range fields {
		n, err := strconv.ParseFloat(field, 32)
		if err != nil {
			return err
		}
		parsed[i] = float32(n)
	}
	value.Set(reflect.ValueOf(parsed))
	return nil
}

// Populates the given float64 slice with values parsed from fields.
// Returns an error if parsing fails.
func fillFloat64Slice(value reflect.Value, fields []string) error {
	parsed := make([]float64, len(fields))
	for i, field := range fields {
		n, err := strconv.ParseFloat(field, 64)
		if err != nil {
			return err
		}
		parsed[i] = float64(n)
	}
	value.Set(reflect.ValueOf(parsed))
	return nil
}

// Populates the given string slice with values parsed from fields.
// Returns an error if parsing fails.
func fillStringSlice(value reflect.Value, fields []string) error {
	// Fields may be a part of a bigger slice, so copying to allow the big
	// slice to get CG'ed.
	slice := make([]string, len(fields))
	copy(slice, fields)
	value.Set(reflect.ValueOf(slice))
	return nil
}
