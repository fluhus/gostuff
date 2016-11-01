package wordnet

import (
	"encoding/json"
	"reflect"
	"strings"
	"testing"
)

func TestDataParser(t *testing.T) {
	expected := map[string]*Synset{
		"v111": &Synset{
			"111",
			"v",
			[]string{
				"foo",
				"bar",
				"baz",
			},
			[]*Pointer{
				{"!", "n123", -1, -1},
				{"@", "a321", 0, 1},
			},
			[]*Frame{
				{4, 4},
				{6, 6},
			},
			"hello world",
			nil,
		},
	}

	actual := map[string]*Synset{}
	err := parseDataFile(strings.NewReader(testData), "v", map[string][]int{},
		actual)
	if err != nil {
		t.Fatal("Parsing error:", err)
	}
	if !reflect.DeepEqual(expected, actual) {
		t.Error("Non-equal values:")
		t.Error(stringify(expected))
		t.Error(stringify(actual))
	}
	for key, ss := range actual {
		if ss.Id() != key {
			t.Errorf("ss.Id()=%v, want %v", ss.Id(), key)
		}
	}
}

func TestExceptionParser(t *testing.T) {
	expected := map[string][]string{
		"n.foo": []string{"n.bar"},
		"n.baz": []string{"n.bla", "n.blu"},
	}

	actual := map[string][]string{}
	err := parseExceptionFile(strings.NewReader(testException), "n", actual)
	if err != nil {
		t.Fatal("Parsing error:", err)
	}
	if !reflect.DeepEqual(expected, actual) {
		t.Error("Non-equal values:")
		t.Error(stringify(expected))
		t.Error(stringify(actual))
	}
}

func TestExampleIndexParser(t *testing.T) {
	expected := map[string][]int{
		"abash.37.0": []int{126, 127},
		"abhor.37.0": []int{138, 139, 15},
	}

	actual, err := parseExampleIndex(strings.NewReader(testExampleIndex))
	if err != nil {
		t.Fatal("Parsing error:", err)
	}
	if !reflect.DeepEqual(expected, actual) {
		t.Error("Non-equal values:")
		t.Error(expected)
		t.Error(actual)
	}
}

func TestExampleParser(t *testing.T) {
	expected := map[string]string{
		"111": "hello world",
		"222": "goodbye universe",
	}

	actual, err := parseExamples(strings.NewReader(testExamples))
	if err != nil {
		t.Fatal("Parsing error:", err)
	}
	if !reflect.DeepEqual(expected, actual) {
		t.Error("Non-equal values:")
		t.Error(expected)
		t.Error(actual)
	}
}

func TestIndexParser(t *testing.T) {
	expected := map[string][]string{
		"n.thing":  {"na", "nb"},
		"v.thing2": {"vc", "vd"},
	}

	actual, err := parseIndex(strings.NewReader(testIndex))
	if err != nil {
		t.Fatal("Parsing error:", err)
	}
	if !reflect.DeepEqual(expected, actual) {
		t.Error("Non-equal values:")
		t.Error(expected)
		t.Error(actual)
	}
}

func stringify(a interface{}) string {
	j, _ := json.Marshal(a)
	return string(j)
}

var testData = `  copyright line
111 1 v 3 foo 1 bar 3 baz 5 2 ! 123 n 0000 @ 321 a 0102 2 + 4 5 + 6 7 | hello world`

var testException = `foo bar
baz bla blu`

var testExampleIndex = `abash%2:37:00:: 126,127
abhor%2:37:00:: 138,139,15`

var testExamples = `111 hello world
222 goodbye universe`

var testIndex = `  copyright line
thing n 2 3 x y z 2 2 a b
thing2 v 4 1 x 4 2 c d e f`
