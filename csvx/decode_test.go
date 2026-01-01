package csvx

import (
	"bytes"
	"fmt"
	"reflect"
	"strconv"
	"testing"
)

func TestDecodeReader_basic(t *testing.T) {
	tests := []testCase[basicItem]{
		{"stringy,inty,floaty\nbla,5,6", basicItem{Stringy: "bla", Inty: 5, Floaty: 6}, false},
		{"inty,floaty,stringy\n5,6,bla", basicItem{Stringy: "bla", Inty: 5, Floaty: 6}, false},
		{"floaty,stringy,inty\n6,bla,5", basicItem{Stringy: "bla", Inty: 5, Floaty: 6}, false},
		{"InTy,STRINGY,floaty\n5,bla,6", basicItem{Stringy: "bla", Inty: 5, Floaty: 6}, false},

		{"stringy,inty,floaty\nbla,a,6", basicItem{}, true},
		{"stringy,inty,floaty\nbla,5,a", basicItem{}, true},
		{"stringy,inty\nbla,5", basicItem{}, true},
		{"stringyy,inty,floaty\nbla,5,6", basicItem{}, true},
		{"stringy,intyy,floaty\nbla,5,6", basicItem{}, true},
		{"stringy,intyy,floatyy\nbla,5,6", basicItem{}, true},
	}
	testGeneric(t, true, tests)
}

func TestDecodeReader_basicNoHeader(t *testing.T) {
	tests := []testCase[basicItem]{
		{"bla,5,6", basicItem{Stringy: "bla", Inty: 5, Floaty: 6}, false},
		{"bla", basicItem{}, true},
	}
	testGeneric(t, false, tests)
}

func TestDecodeReader_tagged(t *testing.T) {
	tests := []testCase[taggedItem]{
		{"intu,stringy,floaty\n5,bla,6", taggedItem{Stringy: "bla", Inty: 5, Floaty: 6}, false},
		{"intu,blabla,floaty\n5,bla,6", taggedItem{Stringy: "bla", Inty: 5, Floaty: 6}, false},
		{"floaty,stringy,intu\n6,bla,5", taggedItem{Stringy: "bla", Inty: 5, Floaty: 6}, false},

		{"inty,stringy,floaty\n5,bla,6", taggedItem{}, true},
		{"intU,stringy,floaty\n5,bla,6", taggedItem{}, true},
		{"inty,floaty,stringy\n5,6,bla", taggedItem{}, true},
	}
	testGeneric(t, true, tests)
}

func TestDecodeReader_taggedNoHeader(t *testing.T) {
	tests := []testCase[numberedItem]{
		{"5,bla", numberedItem{Stringy: "bla", Inty: 5}, false},
	}
	testGeneric(t, false, tests)
}

func TestDecodeReader_missing(t *testing.T) {
	tests := []testCase[missingItem]{
		{"stringy\nbla", missingItem{Stringy: "bla"}, false},
		{"stringy,nothing\nbla,blu", missingItem{Stringy: "bla"}, false},
		{"stringy,inty\nbla,blu", missingItem{Stringy: "bla"}, false},
		{"stringy,inty\nbla,5", missingItem{Stringy: "bla"}, false},
		{"stringy,intu\nbla,5", missingItem{Stringy: "bla", Inty: 5}, false},
		{"stringy,floaty\nbla,6", missingItem{Stringy: "bla", Floaty: 6}, false},
		{"stringy,intu,floaty\nbla,5,6", missingItem{Stringy: "bla", Inty: 5, Floaty: 6}, false},
		{"intu,floaty,stringy\n5,6,bla", missingItem{Stringy: "bla", Inty: 5, Floaty: 6}, false},

		{"intu,floaty\n5,6", missingItem{}, true},
	}
	testGeneric(t, true, tests)
}

func TestDecodeReader_parse(t *testing.T) {
	tests := []testCase[parseItem]{
		{"inty,stringy,floaty\n123,bla,6", parseItem{Stringy: "la", Inty: 1, Floaty: 6}, false},
		{"inty,stringy,floaty\n,bla,6", parseItem{Stringy: "la", Floaty: 6}, false},

		{"inty,stringy,floaty\na,bla,6", parseItem{}, true},
		{"inty,stringy,floaty\n123,b,6", parseItem{}, true},
		{"inty,stringy,floaty\n123,,6", parseItem{}, true},
	}
	testGeneric(t, true, tests)
}

func TestDecodeReader_parseNoHeader(t *testing.T) {
	tests := []testCase[somewhatNumberedItem]{
		{"bla,123,543,3.14,666", somewhatNumberedItem{
			Stringy: "bla", Inty: 123, Floaty: 3.14, Parsey: 5}, false},

		{"a,bla,6", somewhatNumberedItem{}, true},
	}
	testGeneric(t, false, tests)
}

func TestDecodeReader_sameCol(t *testing.T) {
	tests := []testCase[sameColItem]{
		{"numby\n123", sameColItem{Stringy: "123", Inty: 123}, false},
	}
	testGeneric(t, true, tests)
}

func TestDecodeReader_bad(t *testing.T) {
	testGeneric(t, true, []testCase[badParserItem1]{{"stringy\n", badParserItem1{}, true}})
	testGeneric(t, true, []testCase[badParserItem2]{{"stringy\n", badParserItem2{}, true}})
	testGeneric(t, true, []testCase[badParserItem3]{{"stringy\n", badParserItem3{}, true}})
	testGeneric(t, true, []testCase[badParserItem4]{{"stringy\n", badParserItem4{}, true}})
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

type basicItem struct {
	Stringy string
	Inty    int
	Floaty  float64
}

type taggedItem struct {
	Stringy string `csvx:"1"`
	Inty    int    `csvx:"intu"`
	Floaty  float64
}

type missingItem struct {
	Stringy string
	Inty    int     `csvx:"intu,optional"`
	Floaty  float64 `csvx:",optional"`
	Nothing string  `csvx:"-"`
}

type parseItem struct {
	Stringy string `csvx:"1,CheckString"`
	Inty    int    `csvx:",ParseInt,allowempty"`
	Floaty  float64
}

func (t parseItem) ParseInt(s string) (int, error) {
	return strconv.Atoi(s[:1])
}

func (t parseItem) CheckString(s string) (string, error) {
	if len(s) <= 1 {
		return "", fmt.Errorf("")
	}
	return s[1:], nil
}

type sameColItem struct {
	Stringy string `csvx:"0"`
	Inty    int    `csvx:"0"`
}

type numberedItem struct {
	Stringy string `csvx:"1"`
	Inty    int    `csvx:"0"`
}

type somewhatNumberedItem struct {
	Stringy string
	Nothing int `csvx:"-"`
	Inty    int
	Floaty  float64 `csvx:"3"`
	Parsey  int     `csvx:",ParseInt"`
}

func (t somewhatNumberedItem) ParseInt(s string) (int, error) {
	return strconv.Atoi(s[:1])
}

type badParserItem1 struct {
	Stringy string `csvx:",Parse"`
}

func (badParserItem1) Parse(string) string {
	return ""
}

type badParserItem2 struct {
	Stringy string `csvx:",Parse"`
}

func (badParserItem2) Parse(string) (int, error) {
	return 0, nil
}

type badParserItem3 struct {
	Stringy string `csvx:",Parse"`
}

func (badParserItem3) Parse(int) (string, error) {
	return "", nil
}

type badParserItem4 struct {
	Stringy string `csvx:",Parse1"`
}

func (badParserItem4) Parse(string) (string, error) {
	return "", nil
}
