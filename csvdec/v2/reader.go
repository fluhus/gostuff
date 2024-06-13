// Package csvdec provides a generic CSV decoder. Wraps the encoding/csv
// package with a decoder that can populate structs and slices.
package csvdec

import (
	"encoding/csv"
	"fmt"
	"io"
	"reflect"
)

// A Reader reads CSV lines and converts them to data objects. Embeds a
// csv.Reader, so it can be used the same way.
type Reader struct {
	*csv.Reader
	SkipCols uint // How many columns to skip from the beginning of each line.
}

// New returns a new reader that reads from r. skipRows and skipCols
// indicate how many of the first rows and columns should be ignored.
func New(r io.Reader) *Reader {
	reader := csv.NewReader(r)
	return &Reader{reader, 0}
}

// SkipRow skips a row and returns an error if reading failed.
func (d *Reader) SkipRow() error {
	_, err := d.Read()
	return err
}

// ReadInto reads the next CSV line and populates the given object with parsed
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
func (d *Reader) ReadInto(a interface{}) error {
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
		panic("input must be a pointer. Got: " + value.Type().String())
	}
	elem := value.Elem()
	typ := elem.Type()

	switch typ.Kind() {
	case reflect.Struct:
		return fillStruct(elem, fields)
	case reflect.Slice:
		return fillSlice(elem, fields)
	default:
		panic("unsupported type: " + value.Type().String())
	}
}

// Iter iterates over the remaining entries, yielding instances of T.
func Iter[T any](r *Reader) func(yield func(T, error) bool) {
	return func(yield func(T, error) bool) {
		for {
			var t T
			err := r.ReadInto(&t)
			if err == io.EOF {
				return
			}
			if err != nil {
				yield(t, err)
				return
			}
			if !yield(t, nil) {
				return
			}
		}
	}
}
