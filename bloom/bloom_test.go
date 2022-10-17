package bloom

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/fluhus/gostuff/bnry"
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
		f := New(test.bits, 1)
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

func TestFilter(t *testing.T) {
	f := New(80, 4)
	data := []byte{1, 2, 3, 4}
	if f.Has(data) {
		t.Fatalf("Has(%v)=true, want false", data)
	}
	if f.Add(data) {
		t.Fatalf("Add(%v)=true, want false", data)
	}
	if !f.Has(data) {
		t.Fatalf("Has(%v)=false, want true", data)
	}
	if !f.Add(data) {
		t.Fatalf("Add(%v)=false, want true", data)
	}

	data2 := []byte{4, 3, 2, 1}
	if f.Has(data2) {
		t.Fatalf("Has(%v)=true, want false", data2)
	}
}

func TestNewOptimal(t *testing.T) {
	n := 1000000
	p := 0.01
	f := NewOptimal(n, p)
	t.Logf("bits=%v, k=%v", f.NBits(), f.NHash())
	fp := 0
	for i := 0; i < n; i++ {
		buf, _ := bnry.MarshalBinary(uint64(i))
		if f.Add(buf) {
			fp++
		}
	}
	if fpr := float64(fp) / float64(n); fpr > p {
		t.Fatalf("fp=%v, want <%v", fpr, p)
	}
}

func TestEncode(t *testing.T) {
	data1 := []byte{1, 2, 3, 4}
	data2 := []byte{4, 3, 2, 1}
	f1 := New(80, 4)
	f1.SetSeed(5678)
	f1.Add(data1)

	if !f1.Has(data1) {
		t.Fatalf("Has(%v)=false, want true", data1)
	}
	if f1.Has(data2) {
		t.Fatalf("Has(%v)=true, want false", data2)
	}

	buf := bytes.NewBuffer(nil)
	if err := f1.Encode(buf); err != nil {
		t.Fatalf("Encode(...) failed: %v", err)
	}
	f2 := &Filter{}
	if err := f2.Decode(buf); err != nil {
		t.Fatalf("Decode(...) failed: %v", err)
	}

	if !bytes.Equal(f1.b, f2.b) {
		t.Fatalf("Decode(...) bytes=%v, want %v", f2.b, f1.b)
	}
	if f1.seed != f2.seed {
		t.Fatalf("Decode(...) seed=%v, want %v", f2.seed, f1.seed)
	}

	if !f2.Has(data1) {
		t.Fatalf("Decode(...).Has(%v)=false, want true", data1)
	}
	if f2.Has(data2) {
		t.Fatalf("Decode(...).Has(%v)=true, want false", data2)
	}
}

func TestAddFilter(t *testing.T) {
	data := [][]byte{
		{1, 2, 3, 4},
		{5, 6, 7, 8},
		{9, 10},
	}
	f1 := New(100, 3)
	f2 := New(100, 3)
	f2.SetSeed(f1.Seed())

	f1.Add(data[0])
	f2.Add(data[1])
	f1.AddFilter(f2)

	if !f1.Has(data[0]) {
		t.Fatalf("Has(%v)=false, want true", data[0])
	}
	if !f1.Has(data[1]) {
		t.Fatalf("Has(%v)=false, want true", data[1])
	}
	if f1.Has(data[2]) {
		t.Fatalf("Has(%v)=true, want false", data[2])
	}
}

func TestNElements(t *testing.T) {
	bf := New(2000, 4)
	bf.SetSeed(0)
	n := 100
	errs := 0
	for i := 1; i <= n; i++ {
		bf.Add([]byte{byte(i)})
		if gnum.Diff(bf.NElements(), i) > 1 {
			errs++
		}
	}
	if errs > n/10 {
		t.Fatalf("Too many errors: %v, want %v", errs, n/10)
	}
}

func BenchmarkHas(b *testing.B) {
	for _, n := range []int{10, 30, 100} {
		for k := 1; k <= 3; k++ {
			b.Run(fmt.Sprintf("n=%v,k=%v", n, k), func(b *testing.B) {
				f := New(1000000, k)
				buf := make([]byte, n)
				f.Add(buf)
				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					f.Has(buf)
				}
			})
		}
	}
}

func BenchmarkAdd(b *testing.B) {
	for _, n := range []int{10, 30, 100} {
		for k := 1; k <= 3; k++ {
			b.Run(fmt.Sprintf("n=%v,k=%v", n, k), func(b *testing.B) {
				f := New(1000000, k)
				buf := make([]byte, n)
				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					f.Add(buf)
				}
			})
		}
	}
}
