package snm

import (
	"cmp"
	"encoding/json"
	"fmt"
	"math"
	"math/rand/v2"
	"slices"
	"testing"

	"golang.org/x/exp/maps"
)

func TestSlice(t *testing.T) {
	want := []int{1, 4, 9, 16}
	got := Slice(4, func(i int) int { return (i + 1) * (i + 1) })
	if !slices.Equal(got, want) {
		t.Fatalf("Slice((i+1)*(i+1))=%v, want %v", got, want)
	}
}

func TestSliceToSlice(t *testing.T) {
	input := []int{1, 4, 9, 16}
	want := []float64{1.5, 4.5, 9.5, 16.5}
	got := SliceToSlice(input, func(i int) float64 {
		return float64(i) + 0.5
	})
	if !slices.Equal(got, want) {
		t.Fatalf("SliceToSlice(%v)=%v, want %v", input, got, want)
	}
}

func TestMapToMap(t *testing.T) {
	input := map[string]string{"a": "bbb", "cccc": "ddddddd"}
	want := map[int]int{1: 3, 4: 7}
	got := MapToMap(input, func(k, v string) (int, int) {
		return len(k), len(v)
	})
	if !maps.Equal(got, want) {
		t.Fatalf("MapToMap(%v)=%v, want %v", input, got, want)
	}
}

func TestMapToMap_equalKeys(t *testing.T) {
	input := map[string]string{"a": "bbb", "cccc": "ddddddd", "e": "ff"}
	want1 := map[int]int{1: 3, 4: 7}
	want2 := map[int]int{1: 2, 4: 7}
	got := MapToMap(input, func(k, v string) (int, int) {
		return len(k), len(v)
	})
	if !maps.Equal(got, want1) && !maps.Equal(got, want2) {
		t.Fatalf("MapToMap(%v)=%v, want %v or %v", input, got, want1, want2)
	}
}

func TestDefaultMap(t *testing.T) {
	m := NewDefaultMap(func(i int) string {
		return fmt.Sprint(i + 1)
	})
	if got, want := m.Get(2), "3"; got != want {
		t.Fatalf("Get(%d)=%s, want %s", 2, got, want)
	}
	if got, want := m.Get(6), "7"; got != want {
		t.Fatalf("Get(%d)=%s, want %s", 6, got, want)
	}
	m.Set(2, "a")
	if got, want := m.Get(2), "a"; got != want {
		t.Fatalf("Get(%d)=%s, want %s", 2, got, want)
	}
	if got, want := m.Get(6), "7"; got != want {
		t.Fatalf("Get(%d)=%s, want %s", 6, got, want)
	}
	if got, want := len(m.M), 2; got != want {
		t.Fatalf("Len=%d, want %d", got, want)
	}
}

func TestCompareReverse(t *testing.T) {
	input := []int{3, 4, 2, 1, 5}
	want := []int{5, 4, 3, 2, 1}

	cp := slices.Clone(input)
	slices.SortFunc(cp, CompareReverse)
	if !slices.Equal(cp, want) {
		t.Errorf("SortFunc(%v, Compare)=%v, want %v",
			input, cp, want)
	}
}

func ExampleSortedKeys() {
	ages := map[string]int{
		"Alice":   30,
		"Bob":     20,
		"Charlie": 25,
	}
	for _, name := range SortedKeys(ages) {
		fmt.Printf("%s: %d\n", name, ages[name])
	}
	// Output:
	// Bob: 20
	// Charlie: 25
	// Alice: 30
}

func ExampleSortedKeysFunc_reverse() {
	ages := map[string]int{
		"Alice":   30,
		"Bob":     20,
		"Charlie": 25,
	}
	// Sort by reverse natural order.
	for _, name := range SortedKeysFunc(ages, CompareReverse) {
		fmt.Printf("%s: %d\n", name, ages[name])
	}
	// Output:
	// Alice: 30
	// Charlie: 25
	// Bob: 20
}

func TestEnumerator(t *testing.T) {
	tests := []struct {
		i, want int
	}{
		{6, 0}, {3, 1}, {6, 0}, {2, 2}, {3, 1}, {10, 3}, {10, 3}, {2, 2},
		{6, 0}, {3, 1},
	}
	e := Enumerator[int]{}
	for _, test := range tests {
		if got := e.IndexOf(test.i); got != test.want {
			t.Fatalf("%v.IndexOf(%v)=%v, want %v", e, test.i, got, test.want)
		}
	}

	wantElem := []int{6, 3, 2, 10}
	if got := e.Elements(); !slices.Equal(got, wantElem) {
		t.Fatalf("%v.Elements()=%v, want %v", e, got, wantElem)
	}
}

func ExampleCapMap() {
	data := [][]string{
		{"a", "b", "c", "a", "b", "b"},
		// ...
	}
	counter := NewCapMap[string, int]()
	for _, x := range data {
		m := counter.Map()
		countValues(x, m)

		// Do something with m.
		j, _ := json.Marshal(m)
		fmt.Println(string(j))
		counter.Clear()
	}
	//Output:
	//{"a":2,"b":3,"c":1}
}

func countValues(vals []string, out map[string]int) {
	for _, v := range vals {
		out[v]++
	}
}

func TestShuffle(t *testing.T) {
	nums := Slice(10, func(i int) int { return i })
	found := make([]bool, len(nums))
	counts := Slice(len(nums), func(i int) []int { return make([]int, len(nums)) })
	for range 1000 {
		Shuffle(nums)
		clear(found)
		for i, x := range nums {
			found[x] = true
			counts[i][x]++
		}
		for i, f := range found {
			if !f {
				t.Fatalf("did not find %v: %v", i, nums)
			}
		}
	}
	for i, c := range counts {
		for j, x := range c {
			if x < 70 {
				t.Errorf("count of %v at position %v: %v, want >%v",
					j, i, x, 70)
			}
		}
	}
}

func BenchmarkShuffle(b *testing.B) {
	a := Slice(10000, func(i int) int { return rand.Int() })
	b.Run("snm.Shuffle", func(b *testing.B) {
		for b.Loop() {
			Shuffle(a)
		}
	})
	b.Run("rand.Shuffle", func(b *testing.B) {
		for b.Loop() {
			rand.Shuffle(len(a), func(i, j int) { a[i], a[j] = a[j], a[i] })
		}
	})
}

func TestSortByKey_int(t *testing.T) {
	input := []int{5, 3, 8, 6, 3, 2}
	want := []int{8, 6, 5, 3, 3, 2}
	got := slices.Clone(input)
	SortByKey(got, func(i int) int { return -i })
	if !slices.Equal(got, want) {
		t.Fatalf("sortByKey(%v)=%v, want %v",
			input, got, want)
	}
}

func TestSortByKey_string(t *testing.T) {
	input := []string{"hello", "oi", "bonjour", "shalom", "salam"}
	want := []string{"salam", "hello", "shalom", "oi", "bonjour"}
	got := slices.Clone(input)
	SortByKey(got, func(s string) string { return s[1:] })
	if !slices.Equal(got, want) {
		t.Fatalf("sortByKey(%v)=%v, want %v",
			input, got, want)
	}
}

func BenchmarkSortByKey(b *testing.B) {
	a := Slice(1000, func(i int) int {
		return rand.Int()
	})
	b.Run("SortFunc", func(b *testing.B) {
		aa := slices.Clone(a)
		for b.Loop() {
			copy(aa, a)
			slices.SortFunc(aa, func(a, b int) int {
				return cmp.Compare(math.Log(float64(a)), math.Log(float64(b)))
			})
			if !slices.IsSorted(aa) {
				b.Fatalf("slice not sorted")
			}
		}
	})
	b.Run("SortByKey", func(b *testing.B) {
		aa := slices.Clone(a)
		for b.Loop() {
			copy(aa, a)
			SortByKey(aa, func(i int) float64 {
				return math.Log(float64(i))
			})
			if !slices.IsSorted(aa) {
				b.Fatalf("slice not sorted")
			}
		}
	})
}
