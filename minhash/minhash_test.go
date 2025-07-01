package minhash

import (
	"fmt"
	"hash/crc64"
	"math"
	"math/rand"
	"reflect"
	"slices"
	"sort"
	"testing"
)

func TestCollection(t *testing.T) {
	tests := []struct {
		n     int
		input []uint64
		want  []uint64
	}{
		{
			3,
			[]uint64{1, 2, 2, 2, 2, 1, 1, 3, 3, 3, 1, 2, 3, 1, 3, 3, 2},
			[]uint64{1, 2, 3},
		},
		{
			3,
			[]uint64{1, 2, 3, 4, 5, 6, 7, 8, 9},
			[]uint64{1, 2, 3},
		},
		{
			3,
			[]uint64{9, 8, 7, 6, 5, 4, 3, 2, 1},
			[]uint64{1, 2, 3},
		},
		{
			5,
			[]uint64{40, 19, 55, 10, 32, 1, 100, 5, 99, 16, 16},
			[]uint64{1, 5, 10, 16, 19},
		},
	}
	for _, test := range tests {
		mh := New[uint64](test.n)
		for _, k := range test.input {
			mh.Push(k)
		}
		got := mh.View()
		sort.Slice(got, func(i, j int) bool {
			return got[i] < got[j]
		})
		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("New(%d).Push(%v)=%v, want %v",
				test.n, test.input, got, test.want)
		}
	}
}

func TestJSON(t *testing.T) {
	input := New[int](5)
	input.Push(1)
	input.Push(4)
	input.Push(9)
	input.Push(16)
	input.Push(25)
	input.Push(36)
	jsn, err := input.MarshalJSON()
	if err != nil {
		t.Fatalf("MinHash(1,4,9,16,25,36).MarshalJSON() failed: %v", err)
	}
	got := New[int](2)
	err = got.UnmarshalJSON(jsn)
	if err != nil {
		t.Fatalf("UnmarshalJSON(%q) failed: %v", jsn, err)
	}
	if !slices.Equal(got.View(), input.View()) {
		t.Fatalf("UnmarshalJSON(%q)=%v, want %v", jsn, got, input)
	}
}

func TestJaccard(t *testing.T) {
	tests := []struct {
		a, b []uint64
		k    int
		want float64
	}{
		{[]uint64{1, 2, 3}, []uint64{1, 2, 3}, 3, 1},
		{[]uint64{1, 2, 3}, []uint64{2, 3, 4}, 3, 2.0 / 3.0},
		{[]uint64{2, 3, 4}, []uint64{1, 2, 3}, 3, 2.0 / 3.0},
		{[]uint64{1, 2, 3, 4, 5}, []uint64{1, 3, 5}, 5, 0.6},
	}
	for _, test := range tests {
		a, b := New[uint64](test.k), New[uint64](test.k)
		for _, i := range test.a {
			a.Push(i)
		}
		for _, i := range test.b {
			b.Push(i)
		}
		a.Sort()
		b.Sort()
		if got := a.Jaccard(b); math.Abs(got-test.want) > 0.00001 {
			t.Errorf("Jaccard(%v,%v)=%f, want %f",
				test.a, test.b, got, test.want)
		}
	}
}

func TestCollection_largeInput(t *testing.T) {
	const k = 10000
	tests := []struct {
		from1, to1, from2, to2 int
	}{
		{1, 75000, 25000, 100000},
		{1, 60000, 40000, 60000},
		{1, 60000, 20000, 60000},
		{1, 40000, 40001, 60000},
	}
	for _, test := range tests {
		a, b := New[uint64](k), New[uint64](k)
		h := crc64.New(crc64.MakeTable(crc64.ECMA))
		for i := test.from1; i <= test.to1; i++ {
			h.Reset()
			fmt.Fprint(h, i)
			a.Push(h.Sum64())
		}
		for i := test.from2; i <= test.to2; i++ {
			h.Reset()
			fmt.Fprint(h, i)
			b.Push(h.Sum64())
		}
		a.Sort()
		b.Sort()
		want := float64(test.to1-test.from2+1) / float64(
			test.to2-test.from1+1)
		if got := a.Jaccard(b); math.Abs(got-want) > want/100 {
			t.Errorf("Jaccard(...)=%f, want %f", got, want)
		}
	}
}

func FuzzCollection(f *testing.F) {
	f.Add(1, 2, 3, 4, 5, 6)
	f.Fuzz(func(t *testing.T, a int, b int, c int, d int, e int, f int) {
		col := New[int](2)
		col.Push(a)
		col.Push(b)
		col.Push(c)
		col.Push(d)
		col.Push(e)
		col.Push(f)
		v := col.View()
		if len(v) != 2 {
			t.Errorf("len()=%d, want %d", len(v), 2)
		}
		if v[0] < v[1] {
			t.Errorf("v[0]<v[1]: %d<%d, want >=", v[0], v[1])
		}
	})
}

func TestFrozen(t *testing.T) {
	mh := New[int](3)
	mh.Push(27872)
	mh.Push(16978)
	mh.Push(28696)
	mh.Sort()

	fr := mh.Frozen()
	if !slices.Equal(mh.View(), fr.View()) {
		t.Fatalf("View()=%v, want %v", fr.View(), mh.View())
	}

	mh2 := New[int](3)
	mh.Push(27872)
	mh.Push(16978)
	mh.Push(28697)
	mh2.Sort()

	want := mh.Jaccard(mh2)
	got := fr.Jaccard(mh2.Frozen())
	if got != want {
		t.Fatalf("Jaccard=%v, want %v", got, want)
	}
}

func TestFrozen_modifySort(t *testing.T) {
	mh := New[int](1)
	mh.Push(27872)
	mh = mh.Frozen()
	defer func() {
		recover()
	}()
	mh.Sort()
	t.Fatalf(".Frozen().Sort() succeeded, want panic")
}

func TestFrozen_modifyPush(t *testing.T) {
	mh := New[int](1)
	mh.Push(27872)
	mh = mh.Frozen()
	defer func() {
		recover()
	}()
	mh.Push(123)
	t.Fatalf(".Frozen().Sort() succeeded, want panic")
}

func BenchmarkPush(b *testing.B) {
	nums := rand.Perm(b.N)
	mh := New[int](b.N)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		mh.Push(nums[i])
	}
}
