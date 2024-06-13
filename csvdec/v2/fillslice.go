package csvdec

import (
	"reflect"
	"strconv"
	"unsafe"

	"golang.org/x/exp/constraints"
)

// Populates any slice value.
func fillSlice(value reflect.Value, fields []string) error {
	kind := value.Type().Elem().Kind()
	switch kind {
	case reflect.String:
		return fillStringSlice(value, fields)
	case reflect.Int:
		return fillIntSlice[int](value, fields)
	case reflect.Int8:
		return fillIntSlice[int8](value, fields)
	case reflect.Int16:
		return fillIntSlice[int16](value, fields)
	case reflect.Int32:
		return fillIntSlice[int32](value, fields)
	case reflect.Int64:
		return fillIntSlice[int64](value, fields)
	case reflect.Uint:
		return fillUintSlice[uint](value, fields)
	case reflect.Uint8:
		return fillUintSlice[uint8](value, fields)
	case reflect.Uint16:
		return fillUintSlice[uint16](value, fields)
	case reflect.Uint32:
		return fillUintSlice[uint32](value, fields)
	case reflect.Uint64:
		return fillUintSlice[uint64](value, fields)
	case reflect.Float32:
		return fillFloatSlice[float32](value, fields)
	case reflect.Float64:
		return fillFloatSlice[float64](value, fields)
	}
	panic("unsupported type: " + value.Type().String())
}

// Populates the given int slice with values parsed from fields.
// Returns an error if parsing fails.
func fillIntSlice[T constraints.Signed](value reflect.Value, fields []string) error {
	parsed := make([]T, len(fields))
	size := int(unsafe.Sizeof(T(0))) * 8
	for i, field := range fields {
		n, err := strconv.ParseInt(field, 0, size)
		if err != nil {
			return err
		}
		parsed[i] = T(n)
	}
	value.Set(reflect.ValueOf(parsed))
	return nil
}

// Populates the given int slice with values parsed from fields.
// Returns an error if parsing fails.
func fillUintSlice[T constraints.Unsigned](value reflect.Value, fields []string) error {
	parsed := make([]T, len(fields))
	size := int(unsafe.Sizeof(T(0))) * 8
	for i, field := range fields {
		n, err := strconv.ParseUint(field, 0, size)
		if err != nil {
			return err
		}
		parsed[i] = T(n)
	}
	value.Set(reflect.ValueOf(parsed))
	return nil
}

// Populates the given float slice with values parsed from fields.
// Returns an error if parsing fails.
func fillFloatSlice[T constraints.Float](value reflect.Value, fields []string) error {
	parsed := make([]T, len(fields))
	size := int(unsafe.Sizeof(T(0))) * 8
	for i, field := range fields {
		n, err := strconv.ParseFloat(field, size)
		if err != nil {
			return err
		}
		parsed[i] = T(n)
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
