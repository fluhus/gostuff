// Package maps provides utility functions for handling maps, and transitioning between maps and
// slices.
package maps

import (
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
