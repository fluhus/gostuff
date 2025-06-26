package csvx

import (
	"bytes"
	"fmt"
	"reflect"
	"strconv"
	"testing"
)

func TestDecodeReader_basic(t *testing.T) {
	tests := []testCase[BasicItem]{
		{"stringy,inty,floaty\nbla,5,6", BasicItem{Stringy: "bla", Inty: 5, Floaty: 6}, false},
		{"inty,floaty,stringy\n5,6,bla", BasicItem{Stringy: "bla", Inty: 5, Floaty: 6}, false},
		{"floaty,stringy,inty\n6,bla,5", BasicItem{Stringy: "bla", Inty: 5, Floaty: 6}, false},
		{"InTy,STRINGY,floaty\n5,bla,6", BasicItem{Stringy: "bla", Inty: 5, Floaty: 6}, false},

		{"stringy,inty,floaty\nbla,a,6", BasicItem{}, true},
		{"stringy,inty,floaty\nbla,5,a", BasicItem{}, true},
		{"stringy,inty\nbla,5", BasicItem{}, true},
		{"stringyy,inty,floaty\nbla,5,6", BasicItem{}, true},
		{"stringy,intyy,floaty\nbla,5,6", BasicItem{}, true},
		{"stringy,intyy,floatyy\nbla,5,6", BasicItem{}, true},
	}
	testGeneric(t, true, tests)
}

func TestDecodeReader_basicNoHeader(t *testing.T) {
	tests := []testCase[BasicItem]{
		{"bla,5,6", BasicItem{Stringy: "bla", Inty: 5, Floaty: 6}, false},
		{"bla", BasicItem{}, true},
	}
	testGeneric(t, false, tests)
}

func TestDecodeReader_tagged(t *testing.T) {
	tests := []testCase[TaggedItem]{
		{"intu,stringy,floaty\n5,bla,6", TaggedItem{Stringy: "bla", Inty: 5, Floaty: 6}, false},
		{"intu,blabla,floaty\n5,bla,6", TaggedItem{Stringy: "bla", Inty: 5, Floaty: 6}, false},
		{"floaty,stringy,intu\n6,bla,5", TaggedItem{Stringy: "bla", Inty: 5, Floaty: 6}, false},

		{"inty,stringy,floaty\n5,bla,6", TaggedItem{}, true},
		{"intU,stringy,floaty\n5,bla,6", TaggedItem{}, true},
		{"inty,floaty,stringy\n5,6,bla", TaggedItem{}, true},
	}
	testGeneric(t, true, tests)
}

func TestDecodeReader_taggedNoHeader(t *testing.T) {
	tests := []testCase[NumberedItem]{
		{"5,bla", NumberedItem{Stringy: "bla", Inty: 5}, false},
	}
	testGeneric(t, false, tests)
}

func TestDecodeReader_missing(t *testing.T) {
	tests := []testCase[MissingItem]{
		{"stringy\nbla", MissingItem{Stringy: "bla"}, false},
		{"stringy,nothing\nbla,blu", MissingItem{Stringy: "bla"}, false},
		{"stringy,inty\nbla,blu", MissingItem{Stringy: "bla"}, false},
		{"stringy,inty\nbla,5", MissingItem{Stringy: "bla"}, false},
		{"stringy,intu\nbla,5", MissingItem{Stringy: "bla", Inty: 5}, false},
		{"stringy,floaty\nbla,6", MissingItem{Stringy: "bla", Floaty: 6}, false},
		{"stringy,intu,floaty\nbla,5,6", MissingItem{Stringy: "bla", Inty: 5, Floaty: 6}, false},
		{"intu,floaty,stringy\n5,6,bla", MissingItem{Stringy: "bla", Inty: 5, Floaty: 6}, false},

		{"intu,floaty\n5,6", MissingItem{}, true},
	}
	testGeneric(t, true, tests)
}

func TestDecodeReader_parse(t *testing.T) {
	tests := []testCase[ParseItem]{
		{"inty,stringy,floaty\n123,bla,6", ParseItem{Stringy: "la", Inty: 1, Floaty: 6}, false},
		{"inty,stringy,floaty\n,bla,6", ParseItem{Stringy: "la", Floaty: 6}, false},

		{"inty,stringy,floaty\na,bla,6", ParseItem{}, true},
		{"inty,stringy,floaty\n123,b,6", ParseItem{}, true},
		{"inty,stringy,floaty\n123,,6", ParseItem{}, true},
	}
	testGeneric(t, true, tests)
}

func TestDecodeReader_parseNoHeader(t *testing.T) {
	tests := []testCase[SomewhatNumberedItem]{
		{"bla,123,543,3.14,666", SomewhatNumberedItem{
			Stringy: "bla", Inty: 123, Floaty: 3.14, Parsey: 5}, false},

		{"a,bla,6", SomewhatNumberedItem{}, true},
	}
	testGeneric(t, false, tests)
}

func TestDecodeReader_sameCol(t *testing.T) {
	tests := []testCase[SameColItem]{
		{"numby\n123", SameColItem{Stringy: "123", Inty: 123}, false},
	}
	testGeneric(t, true, tests)
}

func TestDecodeReader_bad(t *testing.T) {
	testGeneric(t, true, []testCase[BadParserItem1]{{"stringy\n", BadParserItem1{}, true}})
	testGeneric(t, true, []testCase[BadParserItem2]{{"stringy\n", BadParserItem2{}, true}})
	testGeneric(t, true, []testCase[BadParserItem3]{{"stringy\n", BadParserItem3{}, true}})
	testGeneric(t, true, []testCase[BadParserItem4]{{"stringy\n", BadParserItem4{}, true}})
	testGeneric(t, true, []testCase[string]{{"stringy\n", "", true}})
}

func testGeneric[T any](t *testing.T, header bool, tests []testCase[T]) {
tloop:
	for i, test := range tests {
		var got T
		r := bytes.NewBufferString(test.input)
		it := DecodeReader[T](r)
		if header {
			it = DecodeReaderHeader[T](r)
		}
		for x, err := range it {
			if err != nil {
				if test.wantErr {
					continue tloop
				}
				t.Fatalf("#%v Reader(%q) got err: %v",
					i+1, test.input, err)
			}
			got = x
		}
		if test.wantErr {
			t.Fatalf("#%v Reader(%q)=%v, want error",
				i+1, test.input, got)
		}
		if !reflect.DeepEqual(got, test.want) {
			t.Fatalf("#%v Reader(%q)=%v, want %v",
				i+1, test.input, got, test.want)
		}
	}
}

type testCase[T any] struct {
	input   string
	want    T
	wantErr bool
}

type BasicItem struct {
	Stringy string
	Inty    int
	Floaty  float64
}

type TaggedItem struct {
	Stringy string `csvdec:"1"`
	Inty    int    `csvdec:"intu"`
	Floaty  float64
}

type MissingItem struct {
	Stringy string
	Inty    int     `csvdec:"intu,optional"`
	Floaty  float64 `csvdec:",optional"`
	Nothing string  `csvdec:"-"`
}

type ParseItem struct {
	Stringy string `csvdec:"1,CheckString"`
	Inty    int    `csvdec:",ParseInt,allowempty"`
	Floaty  float64
}

func (t ParseItem) ParseInt(s string) (int, error) {
	return strconv.Atoi(s[:1])
}

func (t ParseItem) CheckString(s string) (string, error) {
	if len(s) <= 1 {
		return "", fmt.Errorf("")
	}
	return s[1:], nil
}

type SameColItem struct {
	Stringy string `csvdec:"0"`
	Inty    int    `csvdec:"0"`
}

type NumberedItem struct {
	Stringy string `csvdec:"1"`
	Inty    int    `csvdec:"0"`
}

type SomewhatNumberedItem struct {
	Stringy string
	Nothing int `csvdec:"-"`
	Inty    int
	Floaty  float64 `csvdec:"3"`
	Parsey  int     `csvdec:",ParseInt"`
}

func (t SomewhatNumberedItem) ParseInt(s string) (int, error) {
	return strconv.Atoi(s[:1])
}

type BadParserItem1 struct {
	Stringy string `csvdec:",Parse"`
}

func (BadParserItem1) Parse(string) string {
	return ""
}

type BadParserItem2 struct {
	Stringy string `csvdec:",Parse"`
}

func (BadParserItem2) Parse(string) (int, error) {
	return 0, nil
}

type BadParserItem3 struct {
	Stringy string `csvdec:",Parse"`
}

func (BadParserItem3) Parse(int) (string, error) {
	return "", nil
}

type BadParserItem4 struct {
	Stringy string `csvdec:",Parse1"`
}

func (BadParserItem4) Parse(string) (string, error) {
	return "", nil
}
