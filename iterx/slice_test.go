package iterx

import (
	"slices"
	"testing"
)

func TestSlice(t *testing.T) {
	input := []string{"hello", "world", "hi"}
	want := slices.Clone(input)
	var got []string
	for x := range Slice(input) {
		got = append(got, x)
	}
	if !slices.Equal(input, want) {
		t.Fatalf("Slice(%v) changed input to %v", want, input)
	}
	if !slices.Equal(got, want) {
		t.Fatalf("Slice(%v)=%v, want %v", input, got, want)
	}
}

func TestISlice(t *testing.T) {
	input := []string{"hello", "world", "hi"}
	want := slices.Clone(input)
	var got []string
	for i, x := range ISlice(input) {
		if i != len(got) {
			t.Fatalf("ISlice(%v) i=%v, want %v", input, i, len(got))
		}
		got = append(got, x)
	}
	if !slices.Equal(input, want) {
		t.Fatalf("ISlice(%v) changed input to %v", want, input)
	}
	if !slices.Equal(got, want) {
		t.Fatalf("ISlice(%v)=%v, want %v", input, got, want)
	}
}

func TestLimit(t *testing.T) {
	input := []string{"bla", "blu", "bli", "ble"}
	tests := []struct {
		n    int
		want []string
	}{
		{-1, nil}, {0, nil}, {1, input[:1]}, {2, input[:2]},
		{3, input[:3]}, {4, input}, {5, input}, {6, input},
	}
	for _, test := range tests {
		var got []string
		for x := range Limit(Slice(input), test.n) {
			got = append(got, x)
		}
		if !slices.Equal(got, test.want) {
			t.Errorf("Limit(%v,%v)=%v, want %v", input, test.n, got, test.want)
		}
	}
}

func TestSkip(t *testing.T) {
	input := []string{"bla", "blu", "bli", "ble"}
	tests := []struct {
		n    int
		want []string
	}{
		{-1, input}, {0, input}, {1, input[1:]}, {2, input[2:]},
		{3, input[3:]}, {4, nil}, {5, nil}, {6, nil},
	}
	for _, test := range tests {
		var got []string
		for x := range Skip(Slice(input), test.n) {
			got = append(got, x)
		}
		if !slices.Equal(got, test.want) {
			t.Errorf("Skip(%v,%v)=%v, want %v", input, test.n, got, test.want)
		}
	}
}
