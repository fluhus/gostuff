// Package csvx provides convenience wrappers around [encoding/csv].
//
// It provides a set of Decode functions, that populate struct
// fields with delimiter-separated values.
//
// # Decode Accepted Types
//
// The T type parameter for Decode functions accepts structs.
// Fields may be of types bool, int*, uint*, float* or string
// for automatic parsing.
// For manual parsing with a method, any type is allowed.
// Unexported fields are ignored.
//
// # Decode Default Behavior
//
// [DecodeFile] and [DecodeReader] match column to field according
// to their order. For example:
//
//	type a struct {
//	  Name   string  // Matches first column
//	  Age    int     // Matches second column
//	  Height float64 // Matches third column
//	}
//
// [DecodeFileHeader] and [DecodeReaderHeader] match column to field
// according to the first line and the field's name, case insensitively.
// For example:
//
//	type a struct {
//	  Name   string  // Matches a column titled name, Name, nAmE, etc.
//	  Age    int     // Matches a column titled age, Age, aGe, etc.
//	  Height float64 // Matches a column titled height, Height, hEiGhT, etc.
//	}
//
// Note that the Header functions use the first line as metadata,
// while the other two functions expect data starting from the first line.
// In all functions yielding continues upon parsing errors,
// so that a caller may choose to skip lines.
//
// # Decode Field Tags
//
// Field tags can be used to change the default behavior.
// The format for a field tag is as follows:
//
//	Field int `csvdec:"column,modifier1,modifier2..."`
//
// The column part may be:
//   - empty: use the default behavior
//   - column name or index: associate this field with the column with this
//     name or at this 0-based index, case sensitively
//   - "-": a single hyphen, ignore this field entirely
//
// Modifiers may be:
//   - "allowempty": the input value may be empty, in which case no parsing
//     will be attempted
//   - "optional": don't err if the column for this field is missing
//   - exported method name: use T's method with this name to parse the
//     input value
//
// A custom parsing method will replace the default parsing.
// The method's signature must take a string as input,
// and return the field's type and an error.
package csvx

import (
	"fmt"
	"io"
	"iter"
	"reflect"
	"strconv"
	"strings"
)

// TODO(amit): Struct pointers?
// TODO(amit): Handle slices?

// DecodeFile returns an iterator over parsed instances of T,
// using column numbers for matching columns to fields.
//
// fn is an optional function for modifying the CSV parser,
// for example for changing the delimiter.
func DecodeFile[T any](file string, mods ...ReaderModifier) iter.Seq2[T, error] {
	return read[T](File(file, mods...), false)
}

// DecodeReader returns an iterator over parsed instances of T,
// using column numbers for matching columns to fields.
//
// fn is an optional function for modifying the CSV parser,
// for example for changing the delimiter.
func DecodeReader[T any](r io.Reader, mods ...ReaderModifier) iter.Seq2[T, error] {
	return read[T](Reader(r, mods...), false)
}

// DecodeFileHeader returns an iterator over parsed instances of T,
// using the first line for matching columns to fields.
//
// fn is an optional function for modifying the CSV parser,
// for example for changing the delimiter.
func DecodeFileHeader[T any](file string, mods ...ReaderModifier) iter.Seq2[T, error] {
	return read[T](File(file, mods...), true)
}

// DecodeReaderHeader returns an iterator over parsed instances of T,
// using the first line for matching columns to fields.
//
// fn is an optional function for modifying the CSV parser,
// for example for changing the delimiter.
func DecodeReaderHeader[T any](r io.Reader, mods ...ReaderModifier) iter.Seq2[T, error] {
	return read[T](Reader(r, mods...), true)
}

// Turns an iterator over string slices into an iterator over T.
func read[T any](r iter.Seq2[[]string, error], header bool) iter.Seq2[T, error] {
	return func(yield func(T, error) bool) {
		var zero T
		var m map[int][]setter
		first := true
		for line, err := range r {
			if err != nil {
				if !yield(zero, err) {
					return
				}
				continue
			}
			if first {
				first = false
				if header {
					m, err = matchColToField(reflect.TypeFor[T](), line)
				} else {
					m, err = matchColToFieldNoHeader(reflect.TypeFor[T]())
				}
				if err != nil {
					yield(zero, err)
					return
				}
				if header {
					continue
				}
			}
			var t T
			if err := populateStruct(&t, line, m); err != nil {
				if !yield(zero, err) {
					return
				}
			}
			if !yield(t, nil) {
				return
			}
		}
	}
}

// Creates a map from column number to setter functions that
// should run on that column's value, based on the type's metadata.
func matchColToField(t reflect.Type, cols []string) (map[int][]setter, error) {
	if t.Kind() != reflect.Struct {
		return nil, fmt.Errorf("expected struct, got %v", t)
	}
	cs := map[string][]setter{}
	ci := map[string][]setter{}
	ii := map[string][]setter{}

	var found []bool
	csf := map[string][]int{}
	cif := map[string][]int{}
	iif := map[string][]int{}

	for i := range t.NumField() {
		f := t.Field(i)
		if !f.IsExported() {
			found = append(found, true)
			continue
		}
		name, m, mi := strings.ToLower(f.Name), ci, cif
		allowEmpty, optional, parseFunc := false, false, false
		parts := strings.Split(f.Tag.Get("csvdec"), ",")
		if tag := parts[0]; tag != "" {
			if tag == "-" {
				continue
			}
			name, m, mi = tag, cs, csf
			if numeric(tag) {
				m, mi = ii, iif
			}
		}
		mi[name] = append(mi[name], i)
		for _, p := range parts[1:] {
			if p == "allowempty" {
				allowEmpty = true
				continue
			}
			if p == "optional" {
				optional = true
				continue
			}
			method, ok := t.MethodByName(p)
			if !ok {
				return nil, fmt.Errorf("method not found: %v", p)
			}
			if !isParseFunc(method.Type, f.Type) {
				return nil, fmt.Errorf("not a valid parse function: %v %v", p, method.Type)
			}
			v := method.Func
			m[name] = append(m[name], func(dst reflect.Value, src string) error {
				if allowEmpty && src == "" {
					return nil
				}
				out := v.Call([]reflect.Value{dst, reflect.ValueOf(src)})
				if err := valueToError(out[1]); err != nil {
					return err
				}
				dst.Field(i).Set(out[0])
				return nil
			})
			parseFunc = true
		}
		found = append(found, optional)
		if parseFunc {
			continue
		}
		switch f.Type.Kind() {
		case reflect.String:
			m[name] = append(m[name], func(dst reflect.Value, src string) error {
				dst.Field(i).SetString(src)
				return nil
			})
		case reflect.Float32, reflect.Float64:
			m[name] = append(m[name], func(dst reflect.Value, src string) error {
				if allowEmpty && src == "" {
					return nil
				}
				x, err := strconv.ParseFloat(src, f.Type.Bits())
				if err != nil {
					return err
				}
				dst.Field(i).SetFloat(x)
				return nil
			})
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			m[name] = append(m[name], func(dst reflect.Value, src string) error {
				if allowEmpty && src == "" {
					return nil
				}
				x, err := strconv.ParseInt(src, 0, f.Type.Bits())
				if err != nil {
					return err
				}
				dst.Field(i).SetInt(x)
				return nil
			})
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			m[name] = append(m[name], func(dst reflect.Value, src string) error {
				if allowEmpty && src == "" {
					return nil
				}
				x, err := strconv.ParseUint(src, 0, f.Type.Bits())
				if err != nil {
					return err
				}
				dst.Field(i).SetUint(x)
				return nil
			})
		case reflect.Bool:
			m[name] = append(m[name], func(dst reflect.Value, src string) error {
				if allowEmpty && src == "" {
					return nil
				}
				x, err := strconv.ParseBool(src)
				if err != nil {
					return err
				}
				dst.Field(i).SetBool(x)
				return nil
			})
		}
	}
	m := map[int][]setter{}
	for i, col := range cols {
		m[i] = append(m[i], cs[col]...)
		for _, x := range csf[col] {
			found[x] = true
		}
		lower := strings.ToLower(col)
		m[i] = append(m[i], ci[lower]...)
		for _, x := range cif[lower] {
			found[x] = true
		}
	}
	for i, s := range ii {
		n, _ := strconv.Atoi(i) // Not expecting error.
		m[n] = append(m[n], s...)
		for _, x := range iif[i] {
			found[x] = true
		}
	}
	for i := range found {
		if !found[i] {
			return nil, fmt.Errorf("field not matched in input: %v",
				t.Field(i).Name)
		}
	}
	return m, nil
}

// Creates a map from column number to setter functions that
// should run on that column's value, based on the type's metadata.
func matchColToFieldNoHeader(t reflect.Type) (map[int][]setter, error) {
	if t.Kind() != reflect.Struct {
		return nil, fmt.Errorf("expected struct, got %v", t)
	}

	ii := map[string][]setter{} // TODO(amit): Change to int keys.
	cur := 0

	for i := range t.NumField() {
		f := t.Field(i)
		if !f.IsExported() {
			continue
		}
		name := fmt.Sprint(cur)
		allowEmpty, parseFunc := false, false
		parts := strings.Split(f.Tag.Get("csvdec"), ",")
		if tag := parts[0]; tag != "" {
			if tag == "-" {
				continue
			}
			if !numeric(tag) {
				return nil, fmt.Errorf(
					"field %q has a string tag name, which is not allowed in no-header mode",
					f.Name)
			}
			name = tag
		} else {
			cur++
		}
		for _, p := range parts[1:] {
			if p == "allowempty" {
				allowEmpty = true
				continue
			}
			method, ok := t.MethodByName(p)
			if !ok {
				return nil, fmt.Errorf("method not found: %v", p)
			}
			if !isParseFunc(method.Type, f.Type) {
				return nil, fmt.Errorf("not a valid parse function: %v %v", p, method.Type)
			}
			v := method.Func
			ii[name] = append(ii[name], func(dst reflect.Value, src string) error {
				if allowEmpty && src == "" {
					return nil
				}
				out := v.Call([]reflect.Value{dst, reflect.ValueOf(src)})
				if err := valueToError(out[1]); err != nil {
					return err
				}
				dst.Field(i).Set(out[0])
				return nil
			})
			parseFunc = true
		}
		if parseFunc {
			continue
		}
		switch f.Type.Kind() {
		case reflect.String:
			ii[name] = append(ii[name], func(dst reflect.Value, src string) error {
				dst.Field(i).SetString(src)
				return nil
			})
		case reflect.Float32, reflect.Float64:
			ii[name] = append(ii[name], func(dst reflect.Value, src string) error {
				if allowEmpty && src == "" {
					return nil
				}
				x, err := strconv.ParseFloat(src, f.Type.Bits())
				if err != nil {
					return err
				}
				dst.Field(i).SetFloat(x)
				return nil
			})
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			ii[name] = append(ii[name], func(dst reflect.Value, src string) error {
				if allowEmpty && src == "" {
					return nil
				}
				x, err := strconv.ParseInt(src, 0, f.Type.Bits())
				if err != nil {
					return err
				}
				dst.Field(i).SetInt(x)
				return nil
			})
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			ii[name] = append(ii[name], func(dst reflect.Value, src string) error {
				if allowEmpty && src == "" {
					return nil
				}
				x, err := strconv.ParseUint(src, 0, f.Type.Bits())
				if err != nil {
					return err
				}
				dst.Field(i).SetUint(x)
				return nil
			})
		case reflect.Bool:
			ii[name] = append(ii[name], func(dst reflect.Value, src string) error {
				if allowEmpty && src == "" {
					return nil
				}
				x, err := strconv.ParseBool(src)
				if err != nil {
					return err
				}
				dst.Field(i).SetBool(x)
				return nil
			})
		}
	}
	m := map[int][]setter{}
	for i, s := range ii {
		n, _ := strconv.Atoi(i) // Not expecting error.

		m[n] = append(m[n], s...)
	}
	return m, nil
}

// Populates a's fields given the input values and setter-map.
func populateStruct(a any, vals []string, setters map[int][]setter) error {
	v := reflect.ValueOf(a).Elem()
	for i, ss := range setters {
		if i >= len(vals) {
			return fmt.Errorf("cannot read column #%v, input has %v columns: %q",
				i, len(vals), vals)
		}
		for _, s := range ss {
			if err := s(v, vals[i]); err != nil {
				return err
			}
		}
	}
	return nil
}

// A function that parses a string and sets the given value accordingly.
type setter func(dst reflect.Value, src string) error

// Returns true if the given string contains only digits.
func numeric(s string) bool {
	if s == "" {
		return false
	}
	for _, c := range s {
		if c < '0' || c > '9' {
			return false
		}
	}
	return true
}

// Checks that t's type matches the requirements for populating
// the type of dst.
func isParseFunc(t, dst reflect.Type) bool {
	return t.Kind() == reflect.Func &&
		t.NumIn() == 2 && t.NumOut() == 2 &&
		t.In(1).Kind() == reflect.String &&
		t.Out(0).AssignableTo(dst) &&
		t.Out(1).Implements(reflect.TypeFor[error]())
}

// Returns a as a possibly nil error. Panics if a is not of error type.
func valueToError(a reflect.Value) error {
	if a.IsNil() {
		return nil
	}
	return a.Interface().(error)
}
