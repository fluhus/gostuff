// Package flagx provides additional [flag] functions.
package flagx

import (
	"flag"
	"fmt"
	"regexp"
	"strconv"

	"github.com/fluhus/gostuff/sets"
)

// RegexpFlagSet defines a regular expression flag with specified name,
// default value, and usage string.
// The return value is the address of a regular expression variable that
// stores the value of the flag.
func RegexpFlagSet(fs *flag.FlagSet, name string,
	value *regexp.Regexp, usage string) **regexp.Regexp {
	p := &value
	fs.Func(name, usage, func(s string) error {
		r, err := regexp.Compile(s)
		if err != nil {
			return err
		}
		*p = r
		return nil
	})
	return p
}

// Regexp defines a regular expression flag with specified name,
// default value, and usage string.
// The return value is the address of a regular expression variable that
// stores the value of the flag.
func Regexp(name string, value *regexp.Regexp, usage string) **regexp.Regexp {
	return RegexpFlagSet(flag.CommandLine, name, value, usage)
}

// IntBetweenFlagSet defines an int flag with specified name,
// default value, usage string and bounds.
// The return value is the address of an int variable that
// stores the value of the flag.
func IntBetweenFlagSet(fs *flag.FlagSet, name string,
	value int, usage string, minVal, maxVal int) *int {
	p := &value
	fs.Func(name, usage, func(s string) error {
		i, err := strconv.Atoi(s)
		if err != nil {
			return err
		}
		if i < minVal || i > maxVal {
			return fmt.Errorf("got %d, want %d-%d", i, minVal, maxVal)
		}
		*p = i
		return nil
	})
	return p
}

// IntBetween defines an int flag with specified name,
// default value, usage string and bounds.
// The return value is the address of an int variable that
// stores the value of the flag.
func IntBetween(name string, value int, usage string, minVal, maxVal int) *int {
	return IntBetweenFlagSet(
		flag.CommandLine, name, value, usage, minVal, maxVal)
}

// FloatBetweenFlagSet defines a float flag with specified name,
// default value, usage string and bounds.
// incMin and incMax define whether min and max are included in the
// allowed values.
// The return value is the address of a float variable that
// stores the value of the flag.
func FloatBetweenFlagSet(fs *flag.FlagSet, name string, value float64,
	usage string, minVal, maxVal float64, incMin, incMax bool) *float64 {
	p := &value
	fs.Func(name, usage, func(s string) error {
		f, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return err
		}
		if (incMin && f < minVal) || (!incMin && f <= minVal) ||
			(incMax && f > maxVal) || (!incMax && f >= maxVal) {
			smin, smax := "(", ")"
			if incMin {
				smin = "["
			}
			if incMax {
				smax = "]"
			}
			return fmt.Errorf("got %f, want %s%f,%f%s",
				f, smin, minVal, maxVal, smax)
		}
		*p = f
		return nil
	})
	return p
}

// FloatBetween defines a float flag with specified name,
// default value, usage string and bounds.
// incMin and incMax define whether min and max are included in the
// allowed values.
// The return value is the address of a float variable that
// stores the value of the flag.
func FloatBetween(name string, value float64, usage string,
	minVal, maxVal float64, incMin, incMax bool) *float64 {
	return FloatBetweenFlagSet(
		flag.CommandLine, name, value, usage,
		minVal, maxVal, incMin, incMax)
}

// StringFromFlagSet defines a string flag with specified name,
// default value, usage string and allowed values.
// The return value is the address of a string variable that
// stores the value of the flag.
func StringFromFlagSet(fs *flag.FlagSet, name string, value string,
	usage string, from ...string) *string {
	p := &value
	set := sets.Set[string]{}.Add(from...)
	fs.Func(name, value, func(s string) error {
		if !set.Has(s) {
			return fmt.Errorf("got %s, want one of: %v", s, from)
		}
		*p = s
		return nil
	})
	return p
}

// StringFrom defines a string flag with specified name,
// default value, usage string and allowed values.
// The return value is the address of a string variable that
// stores the value of the flag.
func StringFrom(name string, value string, usage string, from ...string) *string {
	return StringFromFlagSet(flag.CommandLine, name, value, usage, from...)
}
