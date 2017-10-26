// Package flug allows definition of flags directly from struct fields.
//
// The idea is to define flags using struct metadata instead of repeating calls to flag.X.
//
// The `flug` Tag
//
// For a struct's field to be recognized by flug, it has to be exported and have the `flug` tag:
//
//  `flug:"name"`
//  or
//  `flug:"name,usage message"`
//
// The name and usage message will be passed as-are to the flag library, so they have the same
// restrictions.
//
// Usage
//
// Using the package involves creating a struct with flags, registering it, and then continuing
// with the flag library as usual. After calling flag.Parse, the struct's fields will be populated
// with parsed values.
//
//  myFlags := struct {
//    N string `flug:"name,the person's name"`
//    A int    `flug:"age,how many years the person had lived"`
//  }{
//    "joey", // default value for N
//    20,     // default value for A
//  }
//  flug.Register(myFlags)
//
//  flag.Parse() // back to flag
//  fmt.Println(myFlags.N)
package flug

import (
	"flag"
	"fmt"
	"reflect"
	"strings"
)

// Register calls RegisterFlagSet with flag's default flag set.
func Register(a interface{}) error {
	return RegisterFlagSet(a, flag.CommandLine)
}

// RegisterFlagSet adds a's fields as flags to the given flag set.
// Acts only on exported fields that have the `flug` tag.
// Returns an error if a field's type is unsupported or if a field's tag is malformed.
func RegisterFlagSet(a interface{}, f *flag.FlagSet) error {
	fields, err := fieldsOf(a)
	if err != nil {
		return err
	}
	for _, field := range fields {
		err := field.register(f)
		if err != nil {
			return err
		}
	}
	return nil
}

// A flagField contains information about a struct's field, so that it can be made a flag.
type flagField struct {
	name string // Flag name.
	desc string // Flag usage message.
	kind reflect.Kind
	pval interface{} // Pointer to the value, to pass to the flag library's TypeVar function.
}

// For debugging.
func (f *flagField) String() string {
	return fmt.Sprintf("{%v,%v,%v}", f.name, f.kind, f.desc)
}

// fieldsOf converts a structs fields to flag-field objects.
// Acts only on exported fields that have the `flug` tag.
func fieldsOf(a interface{}) ([]*flagField, error) {
	v := reflect.ValueOf(a)

	// Check that input is a struct pointer.
	if v.Kind() != reflect.Ptr || v.Elem().Kind() != reflect.Struct {
		var typeName string
		for v.Kind() == reflect.Ptr {
			typeName += "*"
			v = v.Elem()
		}
		typeName += v.Kind().String()
		return nil, fmt.Errorf("bad input type: %s, expected *struct", typeName)
	}

	v = v.Elem()
	t := v.Type()
	var result []*flagField

	// Go over struct fields.
	for i := 0; i < t.NumField(); i++ {
		// Get tag.
		field := t.Field(i)
		if field.PkgPath != "" {
			continue // Field is unexported.
		}
		tag := field.Tag.Get("flug")
		if tag == "" {
			continue // Field has no flug tag.
		}

		// Extract information.
		f := &flagField{}
		parts := strings.SplitN(tag, ",", 2)
		f.name = parts[0]
		if f.name == "" {
			return nil, fmt.Errorf("field %q has no value for name", field.Name)
		}
		if len(parts) > 1 {
			f.desc = parts[1]
		}
		f.kind = field.Type.Kind()
		f.pval = v.Field(i).Addr().Interface()

		result = append(result, f)
	}

	return result, nil
}

// register adds the flag to the given flag set.
func (f *flagField) register(s *flag.FlagSet) error {
	// TODO(amit): Add the rest of the types that flag library supports.
	switch f.kind {
	case reflect.Int:
		s.IntVar(f.pval.(*int), f.name, *f.pval.(*int), f.desc)
	case reflect.Float64:
		s.Float64Var(f.pval.(*float64), f.name, *f.pval.(*float64), f.desc)
	case reflect.Bool:
		s.BoolVar(f.pval.(*bool), f.name, *f.pval.(*bool), f.desc)
	case reflect.String:
		s.StringVar(f.pval.(*string), f.name, *f.pval.(*string), f.desc)
	default:
		return fmt.Errorf("unsupported flag type: %v", f.kind)
	}
	return nil
}
