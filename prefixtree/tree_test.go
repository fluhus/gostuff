package prefixtree

import (
	"fmt"
	"math/rand/v2"
	"slices"
	"testing"

	"github.com/fluhus/gostuff/sets"
	"github.com/fluhus/gostuff/snm"
	"golang.org/x/exp/maps"
)

func TestHas(t *testing.T) {
	tr := treeOf("amit", "amut", "lavon", "ami")

	has := []string{"amit", "amut", "amit", "lavon"}
	hasnt := []string{"", "a", "am",
		"A", "aa", "amm", "amitt", "amat", "-", "amu",
		"l", "la", "lav", "lavo"}

	for _, s := range has {
		if !tr.Has([]byte(s)) {
			t.Fatalf("Has(%q)=false, want true", s)
		}
	}
	for _, s := range hasnt {
		if tr.Has([]byte(s)) {
			t.Fatalf("Has(%q)=true, want false", s)
		}
	}
}

func TestHasPrefix(t *testing.T) {
	tr := treeOf("amit", "amut", "lavon", "ami")

	has := []string{
		"", "a", "am", "ami", "amu",
		"amit", "amut", "amit",
		"l", "la", "lav", "lavo", "lavon"}
	hasnt := []string{
		"A", "aa", "amm", "amitt", "amat", "-",
		"lavonn"}

	for _, s := range has {
		if !tr.HasPrefix([]byte(s)) {
			t.Fatalf("IsPrefix(%q)=false, want true", s)
		}
	}
	for _, s := range hasnt {
		if tr.HasPrefix([]byte(s)) {
			t.Fatalf("IsPrefix(%q)=true, want false", s)
		}
	}
}

func TestFindPrefixes(t *testing.T) {
	tr := treeOf("amit", "amut", "lavon", "ami")
	tests := []struct {
		input string
		want  []string
	}{
		{"am", nil},
		{"ami", []string{"ami"}},
		{"amit", []string{"ami", "amit"}},
		{"amitt", []string{"ami", "amit"}},
		{"amutt", []string{"amut"}},
		{"", nil},
	}
	for _, test := range tests {
		var got []string
		for _, p := range tr.FindPrefixes([]byte(test.input)) {
			got = append(got, string(p))
		}
		if !slices.Equal(got, test.want) {
			t.Fatalf("FindPrefixes(%q)=%q, want %q", test.input, got, test.want)
		}
	}
}

func TestDelete(t *testing.T) {
	tr := treeOf("amit", "amut", "lavon")

	tests := []struct {
		del     string
		wantDel bool
		want    *Tree
	}{
		{"amam", false, treeOf("amit", "amut", "lavon")},
		{"lavon", true, treeOf("amit", "amut")},
		{"am", false, treeOf("amit", "amut")},
	}

	for _, test := range tests {
		if tr.Delete([]byte(test.del)) != test.wantDel {
			t.Fatalf("Delete(%q)=%v, want %v", test.del, !test.wantDel, test.wantDel)
		}
		if tr.Has([]byte(test.del)) {
			t.Fatalf("Delete(%q)=true, want false", test.del)
		}
		if !treeEqual(tr, test.want) {
			t.Fatalf("tr=%v, want %v", tr, test.want)
		}
	}
}

func TestDeletePrefix(t *testing.T) {
	input := []string{"amit", "amut", "lavon"}

	tests := []struct {
		del     string
		wantDel bool
		want    *Tree
	}{
		{"a", false, treeOf("lavon")},
		{"am", false, treeOf("lavon")},
		{"ami", false, treeOf("lavon", "amut")},
		{"amit", false, treeOf("lavon", "amut")},
		{"amu", false, treeOf("lavon", "amit")},
		{"amut", false, treeOf("lavon", "amit")},
		{"amam", false, treeOf("amit", "amut", "lavon")},
		{"l", true, treeOf("amit", "amut")},
		{"la", true, treeOf("amit", "amut")},
		{"lav", true, treeOf("amit", "amut")},
		{"lavo", true, treeOf("amit", "amut")},
		{"lavon", true, treeOf("amit", "amut")},
		{"", true, treeOf()},
	}

	for _, test := range tests {
		tr := treeOf(input...)
		tr.DeletePrefix([]byte(test.del))
		if !treeEqual(tr, test.want) {
			t.Fatalf("DeletePrefix(%q)=%v, want %v", test.del, tr, test.want)
		}
	}
}

func TestIter(t *testing.T) {
	tr := treeOf("alice", "alicey", "alicie", "bob", "boris", "charles")
	want := sets.Of("alice", "alicey", "alicie", "bob", "boris", "charles")
	got := sets.Set[string]{}
	for x := range tr.Iter() {
		if got.Has(string(x)) {
			t.Fatalf("Iter() visited %q twice", x)
		}
		got.Add(string(x))
	}
	if !maps.Equal(got, want) {
		t.Fatalf("Iter()=%v, want %v", got, want)
	}
}

func TestIterPrefix(t *testing.T) {
	tr := treeOf("alice", "alicey", "alicie", "bob", "boris", "charles")
	tests := []struct {
		p    string
		want []string
	}{
		{"", []string{"alice", "alicey", "alicie", "bob", "boris", "charles"}},
		{"a", []string{"alice", "alicey", "alicie"}},
		{"al", []string{"alice", "alicey", "alicie"}},
		{"ali", []string{"alice", "alicey", "alicie"}},
		{"alice", []string{"alice", "alicey"}},
		{"b", []string{"bob", "boris"}},
		{"bo", []string{"bob", "boris"}},
		{"bob", []string{"bob"}},
		{"bobb", nil},
		{"x", nil},
	}

	for _, test := range tests {
		var got []string
		for x := range tr.IterPrefix([]byte(test.p)) {
			got = append(got, string(x))
		}
		slices.Sort(got)
		if !slices.Equal(got, test.want) {
			t.Fatalf("IterPrefix(%q)=%q, want %q", test.p, got, test.want)
		}
	}
}

func TestIter_empty(t *testing.T) {
	tr := New()
	for x := range tr.Iter() {
		t.Errorf("Iter() yielded %q, want nothing", x)
	}
}

func BenchmarkAdd(b *testing.B) {
	for _, k := range []int{20, 50, 100} {
		b.Run(fmt.Sprint(k), func(b *testing.B) {
			tr := New()
			rawData := snm.Slice(1<<20, func(i int) byte {
				return byte(rand.Uint64())
			})
			data := rawData
			for b.Loop() {
				tr.Add(data[:k])
				data = data[k:]
				if len(data) < k {
					data = rawData
				}
			}
		})
	}
}

func BenchmarkHas(b *testing.B) {
	tr := New()
	text := []byte("aaaaaaaaaaaaaaaaaaa")
	tr.Add(text)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tr.Has(text)
	}
}

func FuzzAdd(f *testing.F) {
	f.Add([]byte{}, []byte{}, []byte{})
	f.Fuzz(func(t *testing.T, a, b, c []byte) {
		tr := New()
		tr.Add(a)
		tr.Add(b)
		tr.Add(c)
		if !tr.Has(a) {
			t.Errorf("Has(%q)=false, want true", a)
		}
		if !tr.Has(b) {
			t.Errorf("Has(%q)=false, want true", b)
		}
		if !tr.Has(c) {
			t.Errorf("Has(%q)=false, want true", c)
		}
	})
}

func treeOf(s ...string) *Tree {
	t := New()
	for _, x := range s {
		t.Add([]byte(x))
	}
	return t
}

func treeEqual(a, b *Tree) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	if a.isElem != b.isElem {
		return false
	}
	if !maps.EqualFunc(a.m, b.m, treeEqual) {
		return false
	}
	return true
}
