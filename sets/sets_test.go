package sets

import (
	"encoding/json"
	"reflect"
	"testing"

	"golang.org/x/exp/maps"
)

func TestJSON(t *testing.T) {
	input := Set[int]{}.Add(1, 3, 6)
	want := Set[int]{}.Add(1, 3, 6)
	j, err := input.MarshalJSON()
	if err != nil {
		t.Fatalf("%v.MarshalJSON() failed: %v", input, err)
	}

	got := Set[int]{}
	err = got.UnmarshalJSON(j)
	if err != nil {
		t.Fatalf("%v.UnmarshalJSON(%q) failed: %v", input, j, err)
	}

	if !maps.Equal(got, want) {
		t.Fatalf("UnmarshalJSON(%q)=%v, want %v", j, got, want)
	}
}

func TestJSON_slice(t *testing.T) {
	input := []Set[int]{
		Set[int]{}.Add(1, 3, 6),
		Set[int]{}.Add(7, 9, 6),
	}
	want := []Set[int]{
		Set[int]{}.Add(1, 3, 6),
		Set[int]{}.Add(7, 9, 6),
	}
	j, err := json.Marshal(input)
	if err != nil {
		t.Fatalf("Marshal(%v) failed: %v", input, err)
	}

	got := []Set[int]{}
	err = json.Unmarshal(j, &got)
	if err != nil {
		t.Fatalf("Unmarshal(%q) failed: %v", j, err)
	}

	if !reflect.DeepEqual(input, got) {
		t.Fatalf("Unmarshal(%q)=%v, want %v", j, got, want)
	}
}

func TestJSON_map(t *testing.T) {
	input := map[string]Set[int]{
		"a": Set[int]{}.Add(1, 3, 6),
		"x": Set[int]{}.Add(7, 9, 6),
	}
	want := map[string]Set[int]{
		"a": Set[int]{}.Add(1, 3, 6),
		"x": Set[int]{}.Add(7, 9, 6),
	}
	j, err := json.Marshal(input)
	if err != nil {
		t.Fatalf("Marshal(%v) failed: %v", input, err)
	}

	got := map[string]Set[int]{}
	err = json.Unmarshal(j, &got)
	if err != nil {
		t.Fatalf("Unmarshal(%q) failed: %v", j, err)
	}

	if !reflect.DeepEqual(input, got) {
		t.Fatalf("Unmarshal(%q)=%v, want %v", j, got, want)
	}
}
