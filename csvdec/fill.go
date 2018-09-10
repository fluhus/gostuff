package csvdec

import (
	"fmt"
	"reflect"
	"strconv"
)

// fillStruct populates a value's fields with the values in slice s. Value is
// assumed to be a struct.
func fillStruct(value reflect.Value, s []string) error {
	// Check number of fields.
	if len(s) < value.NumField() {
		return fmt.Errorf("not enough values to populate all fields (%d/%d)",
			len(s), value.NumField())
	}

	// Go over fields.
	for i := 0; i < value.NumField(); i++ {
		field := value.Field(i)
		kind := field.Kind()

		if !field.CanSet() {
			panic(fmt.Errorf("Field %d: Cannot be set. Is it unexported?", i))
		}

		// Assign value according to type.
		switch {
		case kind == reflect.String:
			field.SetString(s[i])

		case kind >= reflect.Int && kind <= reflect.Int64:
			v, err := strconv.ParseInt(s[i], 0, 64)
			if err != nil {
				return fmt.Errorf("Field %d: %v", i, err)
			}
			field.SetInt(v)

		case kind >= reflect.Uint && kind <= reflect.Uint64:
			v, err := strconv.ParseUint(s[i], 0, 64)
			if err != nil {
				return fmt.Errorf("Field %d: %v", i, err)
			}
			field.SetUint(v)

		case kind == reflect.Float64 || kind == reflect.Float32:
			v, err := strconv.ParseFloat(s[i], 64)
			if err != nil {
				return fmt.Errorf("Field %d: %v", i, err)
			}
			field.SetFloat(v)

		case kind == reflect.Bool:
			v, err := strconv.ParseBool(s[i])
			if err != nil {
				return fmt.Errorf("Field %d: %v", i, err)
			}
			field.SetBool(v)

		default:
			return fmt.Errorf("Field %d: Unsupported field type: %s", i,
				kind.String())
		}
	}

	return nil
}

// fillIntSlice populates the given int slice with values parsed from fields.
// Returns an error if parsing fails.
func fillIntSlice(value reflect.Value, fields []string) error {
	// Type of slice elements.
	typ := value.Type().Elem()

	size := int(typ.Size()) * 8
	slice := reflect.MakeSlice(reflect.SliceOf(typ), 0, len(fields))
	target := reflect.New(typ).Elem()

	// Parse fields.
	for _, field := range fields {
		val, err := strconv.ParseInt(field, 0, size)
		if err != nil {
			return err
		}
		target.SetInt(val)
		slice = reflect.Append(slice, target)
	}

	// Assign new slice.
	value.Set(slice)

	return nil
}

// fillUintSlice populates the given uint slice with values parsed from fields.
// Returns an error if parsing fails.
func fillUintSlice(value reflect.Value, fields []string) error {
	// Type of slice elements.
	typ := value.Type().Elem()

	size := int(typ.Size()) * 8
	slice := reflect.MakeSlice(reflect.SliceOf(typ), 0, len(fields))
	target := reflect.New(typ).Elem()

	// Parse fields.
	for _, field := range fields {
		val, err := strconv.ParseUint(field, 0, size)
		if err != nil {
			return err
		}
		target.SetUint(val)
		slice = reflect.Append(slice, target)
	}

	// Assign new slice.
	value.Set(slice)

	return nil
}

// fillFloatSlice populates the given float slice with values parsed from
// fields. Returns an error if parsing fails.
func fillFloatSlice(value reflect.Value, fields []string) error {
	// Type of slice elements.
	typ := value.Type().Elem()

	size := int(typ.Size()) * 8
	slice := reflect.MakeSlice(reflect.SliceOf(typ), 0, len(fields))
	target := reflect.New(typ).Elem()

	// Parse fields.
	for _, field := range fields {
		val, err := strconv.ParseFloat(field, size)
		if err != nil {
			return err
		}
		target.SetFloat(val)
		slice = reflect.Append(slice, target)
	}

	// Assign new slice.
	value.Set(slice)

	return nil
}

// fillStringSlice populates the given string slice with values parsed from
// fields. Returns an error if parsing fails.
func fillStringSlice(value reflect.Value, fields []string) error {
	// Type of slice elements.
	typ := value.Type().Elem()

	slice := reflect.MakeSlice(reflect.SliceOf(typ), 0, len(fields))
	target := reflect.New(typ).Elem()

	// Append fields.
	for _, field := range fields {
		target.SetString(field)
		slice = reflect.Append(slice, target)
	}

	// Assign new slice.
	value.Set(slice)

	return nil
}
