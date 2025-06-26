package csvx

import (
	"slices"
	"strings"
	"testing"

	"github.com/fluhus/gostuff/iterx"
)

func TestReader(t *testing.T) {
	input := "a,bb,ccc\nddd,ee,f"
	want := [][]string{{"a", "bb", "ccc"}, {"ddd", "ee", "f"}}
	got, err := iterx.CollectErr(Reader(strings.NewReader(input)))
	if err != nil {
		t.Fatalf("Reader(%q) failed: %v", input, err)
	}
	if !slices.EqualFunc(got, want, slices.Equal) {
		t.Fatalf("Reader(%q)=%q, want %q", input, got, want)
	}
}

func TestReaderTSV(t *testing.T) {
	input := "a\tbb\tccc\nddd\tee\tf"
	want := [][]string{{"a", "bb", "ccc"}, {"ddd", "ee", "f"}}
	got, err := iterx.CollectErr(Reader(strings.NewReader(input), TSV))
	if err != nil {
		t.Fatalf("Reader(%q) failed: %v", input, err)
	}
	if !slices.EqualFunc(got, want, slices.Equal) {
		t.Fatalf("Reader(%q)=%q, want %q", input, got, want)
	}
}
