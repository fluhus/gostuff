// Package maps provides utility functions for handling maps, and transitioning between maps and
// slices.
package maps

import (
	"fmt"
	"reflect"
	"sort"
)

// TODO(amit): Consider adding a Flip function, to reverse a map (key<=>value).

// Keys returns the keys of the given map in a sorted slice.
//
// For a map of type map[k]v, result is of type []k.
func Keys(a interface{}) interface{} {
	keys := reflect.ValueOf(a).MapKeys()
	typ := reflect.TypeOf(a).Key()

	// Create result slice.
	resultValue := reflect.MakeSlice(reflect.SliceOf(typ),
		len(keys), len(keys))
	for i := range keys {
		resultValue.Index(i).Set(keys[i])
	}
	result := resultValue.Interface()
	sortSlice(result) // Sort result, for consistency.
	return result
}

// Values returns the values of the given map in a sorted slice. Duplicates are kept.
//
// For a map of type map[k]v, result is of type []v.
func Values(a interface{}) interface{} {
	val := reflect.ValueOf(a)
	keys := val.MapKeys()
	typ := reflect.TypeOf(a).Elem()

	// Create result slice.
	resultValue := reflect.MakeSlice(reflect.SliceOf(typ),
		len(keys), len(keys))
	for i := range keys {
		resultValue.Index(i).Set(val.MapIndex(keys[i]))
	}
	result := resultValue.Interface()
	sortSlice(result) // Sort result, for consistency.
	return result
}

// Of returns a map whose keys are the values of the given slice, and all have the same given
// value.
//
// For a slice of type []k and value of type v, result is of type map[k]v.
//
// Deprecated: use Map instead.
func Of(slice interface{}, value interface{}) interface{} {
	result := reflect.MakeMap(reflect.MapOf(
		reflect.TypeOf(slice).Elem(),
		reflect.TypeOf(value),
	))
	sliceVal := reflect.ValueOf(slice)
	for i := 0; i < sliceVal.Len(); i++ {
		result.SetMapIndex(sliceVal.Index(i), reflect.ValueOf(value))
	}
	return result.Interface()
}

// Dedup returns a copy of the given slice, sorted and with no duplicated values.
// For non-comparable types, removes duplicates without sorting. Input slice is unchanged.
func Dedup(a interface{}) interface{} {
	return Keys(Of(a, struct{}{}))
}

// Map creates a map using the elements in slice a and the values obtaned by calling
// f on those elements. If a is of type []K then f needs to be of type func(K) V and
// the result will be of type map[K]V. Panics if the types are not as expected.
func Map(a interface{}, f interface{}) interface{} {
	ftyp := reflect.TypeOf(f)
	if ftyp.NumIn() != 1 {
		panic(fmt.Sprintf("bad number of inputs to f: %v, want 1", ftyp.NumIn()))
	}
	if ftyp.NumOut() != 1 {
		panic(fmt.Sprintf("bad number of outputs of f: %v, want 1", ftyp.NumOut()))
	}
	atyp := reflect.TypeOf(a)
	if atyp.Kind() != reflect.Slice {
		panic(fmt.Sprintf("bad input type: %v, want slice", atyp.Kind()))
	}
	ktyp := atyp.Elem()
	if ftyp.In(0) != ktyp {
		panic(fmt.Sprintf("bad input type to f: %v, want %v",
			ftyp.In(0).Kind(), ktyp.Kind()))
	}
	vtyp := ftyp.Out(0)
	aval := reflect.ValueOf(a)
	fval := reflect.ValueOf(f)
	n := aval.Len()

	result := reflect.MakeMap(reflect.MapOf(ktyp, vtyp))
	for i := 0; i < n; i++ {
		k := aval.Index(i)
		result.SetMapIndex(k, fval.Call([]reflect.Value{k})[0])
	}
	return result.Interface()
}

// sortSlice sorts a slice, if its elements are comparable. Otherwise, leaves it as is.
func sortSlice(a interface{}) {
	switch a := a.(type) {
	case []int:
		sort.Ints(a)
	case []string:
		sort.Strings(a)
	case []float64:
		sort.Float64s(a)
	case []int8:
		sort.Slice(a, func(i, j int) bool { return a[i] < a[j] })
	case []int16:
		sort.Slice(a, func(i, j int) bool { return a[i] < a[j] })
	case []int32:
		sort.Slice(a, func(i, j int) bool { return a[i] < a[j] })
	case []int64:
		sort.Slice(a, func(i, j int) bool { return a[i] < a[j] })
	case []uint:
		sort.Slice(a, func(i, j int) bool { return a[i] < a[j] })
	case []uint8:
		sort.Slice(a, func(i, j int) bool { return a[i] < a[j] })
	case []uint16:
		sort.Slice(a, func(i, j int) bool { return a[i] < a[j] })
	case []uint32:
		sort.Slice(a, func(i, j int) bool { return a[i] < a[j] })
	case []uint64:
		sort.Slice(a, func(i, j int) bool { return a[i] < a[j] })
	case []float32:
		sort.Slice(a, func(i, j int) bool { return a[i] < a[j] })
	case []bool:
		sort.Slice(a, func(i, j int) bool { return !a[i] && a[j] })
	}
}
