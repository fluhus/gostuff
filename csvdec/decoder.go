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
	skipCols int
}

// NewDecoder returns a new decoder that reads from r. skipRows and skipCols
// indicate how many of the first rows and columns should be ignored.
func NewDecoder(r io.Reader, skipRows, skipCols int) *Decoder {
	if skipRows < 0 || skipCols < 0 {
		panic(fmt.Sprintf("skipRows and skipCols must be non-negative. "+
			"(skipRows=%d skipCols=%d)", skipRows, skipCols))
	}

	reader := csv.NewReader(r)
	for i := 0; i < skipRows; i++ {
		reader.Read()
	}

	return &Decoder{reader, skipCols}
}

// Decode reads the next CSV line and populates the given object with parsed
// values. Accepted input types are:
//
// Struct pointer: all fields must be exported and of type int*, uint* float*
// or string. Fields will be populated by order of appearance. Too few fields in
// the CSV line will result in an error. Excess fields in the CSV line will be
// ignored.
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

	// Skip given number of columns.
	if len(fields) < d.skipCols {
		fields = nil
	} else {
		fields = fields[d.skipCols:]
	}

	// Act according to type.
	value := reflect.ValueOf(a)
	if value.Kind() != reflect.Ptr {
		panic("Input must be a pointer.")
	}
	elem := value.Elem()
	typ := elem.Type()

	switch typ.Kind() {
	case reflect.Struct:
		return fillStruct(value.Elem(), fields)
	case reflect.Slice:
		t := typ.Elem().Kind()
		switch {
		case t >= reflect.Int && t <= reflect.Int64:
			return fillIntSlice(elem, fields)
		case t >= reflect.Uint && t <= reflect.Uint64:
			return fillUintSlice(elem, fields)
		case t == reflect.Float32 || t == reflect.Float64:
			return fillFloatSlice(elem, fields)
		case t == reflect.String:
			return fillStringSlice(elem, fields)
		}
	}

	panic("Unsupported type.")
}
