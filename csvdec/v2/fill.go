package csvdec

import (
	"fmt"
	"reflect"
	"strconv"
)

// TODO(amit): Support bool slices.

// Populates a value's fields with the values in slice s.
// Value is assumed to be a struct.
func fillStruct(value reflect.Value, s []string) error {
	// Check number of fields.
	expectedLength := value.NumField()
	if value.Field(value.NumField()-1).Kind() == reflect.Slice {
		expectedLength--
	}
	if len(s) < expectedLength {
		return fmt.Errorf("not enough values to populate all fields (%d/%d) values: %q",
			len(s), value.NumField(), s)
	}

	// Go over fields.
	for i := 0; i < value.NumField(); i++ {
		field := value.Field(i)
		kind := field.Kind()

		if !field.CanSet() {
			panic(fmt.Errorf("field %d cannot be set. Is it unexported?", i))
		}

		// Assign value according to type.
		switch {
		case kind == reflect.String:
			field.SetString(s[i])

		case kind >= reflect.Int && kind <= reflect.Int64:
			v, err := strconv.ParseInt(s[i], 0, 64)
			if err != nil {
				return fmt.Errorf("field %d: %v", i, err)
			}
			field.SetInt(v)

		case kind >= reflect.Uint && kind <= reflect.Uint64:
			v, err := strconv.ParseUint(s[i], 0, 64)
			if err != nil {
				return fmt.Errorf("field %d: %v", i, err)
			}
			field.SetUint(v)

		case kind == reflect.Float64 || kind == reflect.Float32:
			v, err := strconv.ParseFloat(s[i], 64)
			if err != nil {
				return fmt.Errorf("field %d: %v", i, err)
			}
			field.SetFloat(v)

		case kind == reflect.Bool:
			v, err := strconv.ParseBool(s[i])
			if err != nil {
				return fmt.Errorf("field %d: %v", i, err)
			}
			field.SetBool(v)

		case kind == reflect.Slice:
			if i != value.NumField()-1 {
				panic(fmt.Sprintf("field %v is a slice. A slice may only be"+
					" the last field.", i))
			}
			if err := fillSlice(field, s[i:]); err != nil {
				return fmt.Errorf("field %v: %v", i, err)
			}

		default:
			panic(fmt.Sprintf("field %d is of an unsupported type: %v",
				i, kind))
		}
	}

	return nil
}
