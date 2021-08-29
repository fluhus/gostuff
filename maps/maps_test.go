package maps

import (
	"reflect"
	"testing"
)

func TestKeys(t *testing.T) {
	tests := []struct {
		input interface{}
		want  interface{}
	}{
		{
			map[string]int{"b": 0, "c": 0, "a": 0},
			[]string{"a", "b", "c"},
		},
		{
			map[string]int{},
			[]string{},
		},
		{
			map[int]int{1: 1, -5: 1, 2: 1},
			[]int{-5, 1, 2},
		},
		{
			map[int16]int{1: 1, -5: 1, 2: 1},
			[]int16{-5, 1, 2},
		},
		{
			map[uint]int{1: 1, 5: 1, 2: 1},
			[]uint{1, 2, 5},
		},
		{
			map[uint16]int{1: 1, 5: 1, 2: 1},
			[]uint16{1, 2, 5},
		},
		{
			map[float32]int{1.1: 1, 1.3: 1, 1.2: 1},
			[]float32{1.1, 1.2, 1.3},
		},
		{
			map[float64]int{1.1: 1, 1.3: 1, 1.2: 1},
			[]float64{1.1, 1.2, 1.3},
		},
	}

	for i, test := range tests {
		got := Keys(test.input)
		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("#%v Keys(%v) = %v, want %v", i+1, test.input, got, test.want)
		}
	}
}

func TestValues(t *testing.T) {
	tests := []struct {
		input interface{}
		want  interface{}
	}{
		{
			map[string]int{"b": 2, "c": 1, "a": 3},
			[]int{1, 2, 3},
		},
		{
			map[string]int{},
			[]int{},
		},
		{
			map[int]int{1: -1, -5: 5, 2: -2},
			[]int{-2, -1, 5},
		},
		{
			map[int16]int{1: 1, -5: 1, 2: 1},
			[]int{1, 1, 1},
		},
		{
			map[int]float32{1: 1.1, 3: 1.3, 2: 1.2},
			[]float32{1.1, 1.2, 1.3},
		},
		{
			map[int]float64{1: 1.1, 3: 1.3, 2: 1.2},
			[]float64{1.1, 1.2, 1.3},
		},
	}

	for i, test := range tests {
		got := Values(test.input)
		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("#%v Values(%v) = %v, want %v", i+1, test.input, got, test.want)
		}
	}
}

func TestOf(t *testing.T) {
	tests := []struct {
		slice interface{}
		value interface{}
		want  interface{}
	}{
		{
			[]int{1, 2, 3},
			true,
			map[int]bool{1: true, 2: true, 3: true},
		},
		{
			[]string{"c", "a", "b"},
			1,
			map[string]int{"a": 1, "b": 1, "c": 1},
		},
		{
			[]string{},
			1,
			map[string]int{},
		},
		{
			[]int{2},
			struct{}{},
			map[int]struct{}{2: {}},
		},
	}

	for i, test := range tests {
		got := Of(test.slice, test.value)
		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("#%v Of(%v,%v) = %v, want %v", i+1, test.slice, test.value, got, test.want)
		}
	}
}

func TestDedup(t *testing.T) {
	tests := []struct {
		input interface{}
		want  interface{}
	}{
		{
			[]int{4, 3, 1, 2, 3, 1, 2, 3, 1, 2, 3, 4},
			[]int{1, 2, 3, 4},
		},
		{
			[]string{"b", "b", "A", "A", "a", "a", "b", "b", "A", "A", "a", "a"},
			[]string{"A", "a", "b"},
		},
	}

	for i, test := range tests {
		got := Dedup(test.input)
		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("#%v Dedup(%v) = %v, want %v", i+1, test.input, got, test.want)
		}
	}
}

func TestMap(t *testing.T) {
	tests := []struct {
		a    interface{}
		f    interface{}
		want interface{}
	}{
		{
			[]int{1, 3, 5},
			func(a int) float64 { return float64(a) + 0.5 },
			map[int]float64{1: 1.5, 3: 3.5, 5: 5.5},
		},
		{
			[]string{"a", "bb", "ccc"},
			func(s string) string { return s[:1] },
			map[string]string{"a": "a", "bb": "b", "ccc": "c"},
		},
	}
	for _, test := range tests {
		got := Map(test.a, test.f)
		if !reflect.DeepEqual(got, test.want) {
			t.Fatalf("Map(%v, %v)=%v, want %v", test.a, test.f, got, test.want)
		}
	}
}
