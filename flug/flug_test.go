package flug

import (
	"flag"
	"fmt"
	"testing"
)

func TestBasic(t *testing.T) {
	a := struct {
		A int     `flug:"a"`
		B float64 `flug:"b"`
		C string  `flug:"c"`
		D string  `flug:"d"`
		E float64 `flug:"e"`
		F int     `flug:"f"`
	}{1, 2, "3", "4", 5, 6}
	f := flag.NewFlagSet("", flag.ContinueOnError)
	err := RegisterFlagSet(&a, f)
	if err != nil {
		t.Fatalf("RegisterFlagSet(%v) = error: %v", a, err)
	}
	err = f.Parse([]string{"-d", "40", "-c", "30", "-e", "50", "-b", "20", "-f", "60", "-a", "10"})
	if err != nil {
		t.Fatalf("flag.Parse = error: %v", err)
	}
	want := "{10 20 30 40 50 60}"
	if got := fmt.Sprint(a); got != want {
		t.Errorf("a = %v, want %v", got, want)
	}
}

func TestAdvanced(t *testing.T) {
	a := struct {
		A int `flug:"a"`
		b int `flug:"b"`
		C int `FLUG:"c"`
		D int
	}{1, 2, 3, 4}
	f := flag.NewFlagSet("", flag.ContinueOnError)
	err := RegisterFlagSet(&a, f)
	if err != nil {
		t.Fatalf("Failed to register: %v", err)
	}
	err = f.Parse([]string{"-a", "10"})
	if err != nil {
		t.Errorf("Failed to parse: %v", err)
	}
	badFlags := [][]string{{"-b", "20"}, {"-c", "30"}, {"-d", "40"}}
	for _, b := range badFlags {
		err = f.Parse(b)
		if err == nil {
			t.Errorf("flag.Parse(%v) succeded, expected error.", b)
		}
	}
	want := "{10 2 3 4}"
	if got := fmt.Sprint(a); got != want {
		t.Errorf("a = %v, want %v", got, want)
	}
}

func TestBadType(t *testing.T) {
	f := flag.NewFlagSet("", flag.ContinueOnError)

	// A good case, for control.
	err := RegisterFlagSet(&struct {
		A int `flug:"a"`
	}{}, f)
	if err != nil {
		t.Errorf("RegisterFlagSet() = error: %v", err)
	}

	// Unsupported field type.
	err = RegisterFlagSet(&struct {
		A []int `flug:"a"`
	}{}, f)
	if err == nil {
		t.Errorf("RegisterFlagSet() succeded, expected error.")
	}

	// Input struct is not a pointer.
	err = RegisterFlagSet(struct {
		A int `flug:"a"`
	}{}, f)
	if err == nil {
		t.Errorf("RegisterFlagSet() succeded, expected error.")
	}
}
