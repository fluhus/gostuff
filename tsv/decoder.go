// A generic TSV parser.
package tsv

import (
	"bufio"
	"fmt"
	"io"
	"reflect"
	"strings"
)

// Reads TSV lines and converts them to data objects.
type Decoder struct {
	scanner  *bufio.Scanner
	skipCols int
}

// Returns a new decoder. skipRows and skipCols indicate how many of the first
// rows and columns should be ignored.
func NewDecoder(r io.Reader, skipRows, skipCols int) *Decoder {
	if skipRows < 0 || skipCols < 0 {
		panic(fmt.Sprintf("skipRows and skipCols must be non-negative. "+
			"(skipRows=%d skipCols=%d)", skipRows, skipCols))
	}

	scanner := bufio.NewScanner(r)
	for i := 0; i < skipRows; i++ {
		scanner.Scan()
	}

	return &Decoder{scanner, skipCols}
}

// Reads the next TSV line and populates the given object with parsed values.
// Accepted input types are:
//
// Struct pointer: all fields must be exported and of type int*, uint* float*
// or string. Fields will be populated by order of appearance. Too few fields in
// the TSV line will result in an error. Excess fields in the TSV line will be
// ignored.
//
// Slice pointer of type int*, uint*, float*, string: the pointer will be
// populated with a slice of parsed values, according to the length of the TSV
// line.
//
// Any other type will cause a panic.
func (d *Decoder) Decode(a interface{}) error {
	if !d.scanner.Scan() {
		return io.EOF
	}

	// TODO(amit): Split according to TSV rules - take quoted fields in account.
	fields := strings.Split(d.scanner.Text(), "\t")

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
