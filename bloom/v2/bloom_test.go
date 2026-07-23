package bloom

import (
	"fmt"
	"math/rand/v2"
	"testing"

	"github.com/fluhus/gostuff/gnum"
)

func TestLen(t *testing.T) {
	tests := []struct {
		bits int
		want int
	}{
		{1, 8},
		{2, 8},
		{3, 8},
		{4, 8},
		{5, 8},
		{6, 8},
		{7, 8},
		{8, 8},
		{9, 16},
		{16, 16},
		{17, 24},
	}

	for _, test := range tests {
		f := New[int](test.bits, 1, nil)
		if l := f.NBits(); l != test.want {
			t.Errorf("New(%v,1).Len()=%v, want %v",
				test.bits, l, test.want)
		}
		if f.NHash() != 1 {
			t.Errorf("New(%v,1).K()=%v, want 1",
				test.bits, f.NHash())
		}
	}
}

func TestFilter_int(t *testing.T) {
	testFilter(t, []int{1, 2, 3}, []int{4, 5, 6}, Int[int]())
}

func TestFilter_float(t *testing.T) {
	testFilter(t, []float64{1.1, 22.22, 333.333},
		[]float64{444.444, 55.55, 6.6}, Float[float64]())
}

func TestFilter_string(t *testing.T) {
	testFilter(t, []string{"1", "2", "3"}, []string{"4", "5", "6"}, String)
}

func testFilter[T any](t *testing.T, add, noAdd []T, e func(T) []byte) {
	f := New(80, 4, e)

	for _, x := range add {
		if f.Has(x) {
			t.Fatalf("Has(%v)=true, want false", x)
		}
		if f.Add(x) {
			t.Fatalf("Add(%v)=true, want false", x)
		}
		if !f.Has(x) {
			t.Fatalf("Has(%v)=false, want true", x)
		}
		if !f.Add(x) {
			t.Fatalf("Add(%v)=false, want true", x)
		}
	}

	for _, x := range noAdd {
		if f.Has(x) {
			t.Fatalf("Has(%v)=true, want false", x)
		}
	}
}

func TestNewOptimal(t *testing.T) {
	n := 1000000
	p := 0.01
	m, k := OptimalParams(n, p)
	f := New(m, k, Int[int]())
	t.Logf("bits=%v, k=%v", f.NBits(), f.NHash())
	fp := 0
	for i := range n {
		if f.Add(i) {
			fp++
		}
	}
	if fpr := float64(fp) / float64(n); fpr > p {
		t.Fatalf("fp=%v, want <%v", fpr, p)
	}
}

func TestAddFilter(t *testing.T) {
	f1 := New(100, 3, Int[int]())
	f2 := New(100, 3, Int[int]())

	a, b, c := 10, 20, 30

	f1.Add(a)
	f2.Add(b)
	f1.AddFilter(f2)

	if !f1.Has(a) {
		t.Fatalf("Has(%v)=false, want true", a)
	}
	if !f1.Has(b) {
		t.Fatalf("Has(%v)=false, want true", b)
	}
	if f1.Has(c) {
		t.Fatalf("Has(%v)=true, want false", c)
	}
}

func TestNElements(t *testing.T) {
	bf := New(2000, 4, Int[int]())
	n := 100
	errs := 0
	for i := range n {
		bf.Add(i)
		if gnum.Diff(bf.NElements(), i) > 1 {
			errs++
		}
	}
	if errs > n/10 {
		t.Fatalf("Too many errors: %v, want %v", errs, n/10)
	}
}

func BenchmarkHas(b *testing.B) {
	for k := 1; k <= 4; k++ {
		b.Run(fmt.Sprintf("%v", k), func(b *testing.B) {
			f := New(1000000, k, Int[int]())
			// Make some bits 1 to avoid short-circuiting Has.
			for range 100 {
				f.Add(rand.Int())
			}
			i := 0
			for b.Loop() {
				f.Has(i)
				i++
			}
		})
	}
}

func BenchmarkAdd(b *testing.B) {
	for k := 1; k <= 4; k++ {
		b.Run(fmt.Sprintf("%v", k), func(b *testing.B) {
			f := New(1000000, k, Int[int]())
			i := 0
			for b.Loop() {
				f.Add(i)
				i++
			}
		})
	}
}

func ExampleInt() {
	// A Bloom-filter for ints, with optimal parameters for 10 elements
	// and 0.1% false-positive rate.
	m, k := OptimalParams(10, 0.001)
	f := New(m, k, Int[int]())

	f.Add(1)
	f.Add(3)
	f.Add(5)

	fmt.Println(f.Has(1))
	fmt.Println(f.Has(2))
	fmt.Println(f.Has(3))
	fmt.Println(f.Has(4))
	fmt.Println(f.Has(5))

	//Output:
	// true
	// false
	// true
	// false
	// true
}
