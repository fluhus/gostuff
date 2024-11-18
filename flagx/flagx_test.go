package flagx

import (
	"flag"
	"fmt"
	"testing"
)

func TestRegexp(t *testing.T) {
	fs := flag.NewFlagSet("", flag.PanicOnError)

	re := RegexpFlagSet(fs, "a", nil, "")
	fs.Parse([]string{"-a", "a..b"})
	if (*re).String() != "a..b" {
		t.Errorf("RegexpFlagSet(...)=%q, want %q",
			(*re).String(), "a..b")
	}

	fs = flag.NewFlagSet("", flag.PanicOnError)
	re = RegexpFlagSet(fs, "a", nil, "")
	fs.Parse(nil)
	if (*re) != nil {
		t.Errorf("RegexpFlagSet(...)=%q, want nil", (*re).String())
	}
}

func TestIntBetween(t *testing.T) {
	fs := flag.NewFlagSet("", flag.PanicOnError)

	ii := IntBetweenFlagSet(fs, "i", 3, "", 1, 5)
	if *ii != 3 {
		t.Errorf("IntBetweenFlagSet(...)=%v, want %v", ii, 3)
	}

	// Valid values.
	for i := 1; i <= 5; i++ {
		args := []string{"-i", fmt.Sprint(i)}
		fs.Parse(args)
		if *ii != i {
			t.Errorf("Parse(%v)=%v, want %v", args, ii, i)
		}
	}

	// Invalid values.
	for _, i := range []int{-1, 0, 6, 7, 10} {
		func() {
			args := []string{"-i", fmt.Sprint(i)}
			defer func() {
				recover()
			}()
			fs.Parse(args)
			t.Errorf("Parse(%v)=%v, want error", args, ii)
		}()
	}
}

func TestOneOf_string(t *testing.T) {
	fs := flag.NewFlagSet("", flag.PanicOnError)

	vals := []string{"blue", "yellow", "red"}
	ss := OneOfFlagSet(fs, "s", vals[0], "", vals...)
	if *ss != vals[0] {
		t.Errorf("StringFromFlagSet(...)=%v, want %v", ss, vals[0])
	}

	// Valid values.
	for _, s := range vals {
		args := []string{"-s", s}
		fs.Parse(args)
		if *ss != s {
			t.Errorf("Parse(%v)=%v, want %v", args, ss, s)
		}
	}

	// Invalid values.
	for _, s := range vals {
		func() {
			args := []string{"-s", s + "."}
			defer func() {
				recover()
			}()
			fs.Parse(args)
			t.Errorf("Parse(%v)=%v, want error", args, ss)
		}()
	}
}

func TestOneOf_int(t *testing.T) {
	fs := flag.NewFlagSet("", flag.PanicOnError)

	vals := []int{3, 55, 888}
	oct := []string{"0o3", "0o67", "0o1570"}
	ss := OneOfFlagSet(fs, "s", vals[0], "", vals...)
	if *ss != vals[0] {
		t.Errorf("StringFromFlagSet(...)=%v, want %v", ss, vals[0])
	}

	// Valid values.
	for _, i := range vals {
		args := []string{"-s", fmt.Sprint(i)}
		fs.Parse(args)
		if *ss != i {
			t.Errorf("Parse(%v)=%v, want %v", args, ss, i)
		}
	}

	// Octal representation.
	for i, s := range oct {
		args := []string{"-s", s}
		fs.Parse(args)
		want := vals[i]
		if *ss != want {
			t.Errorf("Parse(%v)=%v, want %v", args, ss, want)
		}
	}

	// Invalid values.
	for _, i := range vals {
		func() {
			args := []string{"-s", fmt.Sprint(i) + "aaa"}
			defer func() {
				recover()
			}()
			fs.Parse(args)
			t.Errorf("Parse(%v)=%v, want error", args, *ss)
		}()
	}
}
