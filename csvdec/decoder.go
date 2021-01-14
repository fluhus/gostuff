// Package csvdec provides a generic CSV decoder. Wraps the encoding/csv
// package with a decoder that can populate structs and slices.
package csvdec

import (
	"encoding/csv"
	"fmt"
	"io"
	"reflect"
)

// A Decoder reads CSV lines and converts them to data objects. Embeds a
// csv.Reader, so it can be used the same way.
type Decoder struct {
	*csv.Reader
	SkipCols uint // How many columns to skip from the beginning of each line.
}

// NewDecoder returns a new decoder that reads from r. skipRows and skipCols
// indicate how many of the first rows and columns should be ignored.
func NewDecoder(r io.Reader) *Decoder {
	reader := csv.NewReader(r)
	return &Decoder{reader, 0}
}

// SkipRow skips a row and returns an error if reading failed.
func (d *Decoder) SkipRow() error {
	_, err := d.Read()
	return err
}

// Decode reads the next CSV line and populates the given object with parsed
// values. Accepted input types are struct pointers and slice pointers, as
// explained below.
//
// Struct pointer: all fields must be exported and of type int*, uint* float*,
// string or bool. Fields will be populated by order of appearance. Too few
// values in the CSV line will result in an error. Excess values in the CSV
// line will be ignored. The struct's last field may be a slice, in which case
// all the remaining values will be parsed for that slice's type, according to
// the restrictions below.
//
// Slice pointer of type int*, uint*, float*, string: the pointer will be
// populated with a slice of parsed values, according to the length of the CSV
// line.
//
// Any other type will cause a panic.
func (d *Decoder) Decode(a interface{}) error {
	fields, err := d.Read()
	if err != nil {
		return err
	}

	// Skip columns.
	if uint(len(fields)) < d.SkipCols {
		return fmt.Errorf("cannot skip %v columns, found only %v columns",
			d.SkipCols, len(fields))
	}
	fields = fields[d.SkipCols:]

	// Act according to type.
	value := reflect.ValueOf(a)
	if value.Kind() != reflect.Ptr {
		panic("Input must be a pointer. Got: " + value.Type().String())
	}
	elem := value.Elem()
	typ := elem.Type()

	switch typ.Kind() {
	case reflect.Struct:
		return fillStruct(elem, fields)
	case reflect.Slice:
		return fillSlice(elem, fields)
	default:
		panic("Unsupported type: " + value.Type().String())
	}
}
