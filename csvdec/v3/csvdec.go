package csvdec

import (
	"encoding/csv"
	"fmt"
	"io"
	"iter"
	"reflect"
	"strconv"
	"strings"

	"github.com/fluhus/gostuff/iterx"
)

// TODO(amit): Change to 1-based column index
// TODO(amit): Struct pointers
// TODO(amit): Handle slices
// TODO(amit): Rename stuff
// TODO(amit): Write comments

func File[T any](file string, fn func(*csv.Reader)) iter.Seq2[T, error] {
	return yalla[T](iterx.CSVFile(file, fn), false)
}

func Reader[T any](r io.Reader, fn func(*csv.Reader)) iter.Seq2[T, error] {
	return yalla[T](iterx.CSVReader(r, fn), false)
}

func FileHeader[T any](file string, fn func(*csv.Reader)) iter.Seq2[T, error] {
	return yalla[T](iterx.CSVFile(file, fn), true)
}

func ReaderHeader[T any](r io.Reader, fn func(*csv.Reader)) iter.Seq2[T, error] {
	return yalla[T](iterx.CSVReader(r, fn), true)
}

func yalla[T any](r iter.Seq2[[]string, error], header bool) iter.Seq2[T, error] {
	return func(yield func(T, error) bool) {
		var zero T
		var m map[int][]setter
		first := true
		for line, err := range r {
			if err != nil {
				yield(zero, err)
				return
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
			var l T
			if err := populateStruct(&l, line, m); err != nil {
				yield(zero, err)
				return
			}
			if !yield(l, nil) {
				return
			}
		}
	}
}

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

func populateStruct(a any, vals []string, setters map[int][]setter) error {
	v := reflect.ValueOf(a).Elem()
	for i, ss := range setters {
		if i >= len(vals) {
			return fmt.Errorf("cannot read column #%v, input has %v columns",
				i, len(vals))
		}
		for _, s := range ss {
			if err := s(v, vals[i]); err != nil {
				return err
			}
		}
	}
	return nil
}

type setter func(dst reflect.Value, src string) error

func numeric(s string) bool {
	for _, c := range s {
		if c < '0' || c > '9' {
			return false
		}
	}
	return true
}

func isParseFunc(t, dst reflect.Type) bool {
	return t.Kind() == reflect.Func &&
		t.NumIn() == 2 && t.NumOut() == 2 &&
		t.In(1).Kind() == reflect.String &&
		t.Out(0).AssignableTo(dst) &&
		t.Out(1).Implements(reflect.TypeFor[error]())
}

func valueToError(a reflect.Value) error {
	if a.IsNil() {
		return nil
	}
	return a.Interface().(error)
}
