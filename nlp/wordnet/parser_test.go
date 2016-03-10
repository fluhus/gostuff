package wordnet

import (
	"encoding/json"
	"reflect"
	"strings"
	"testing"
)

func TestDataParser(t *testing.T) {
	expected := map[string]*Synset{
		"v.111": &Synset{
			"v",
			[]*DataWord{
				&DataWord{"foo", 1},
				&DataWord{"bar", 3},
				&DataWord{"baz", 5},
			},
			[]*DataPtr{
				&DataPtr{"!", "n.123", 0, 0},
				&DataPtr{"@", "a.321", 1, 2},
			},
			[]*DataFrame{
				&DataFrame{4, 5},
				&DataFrame{6, 7},
			},
			"hello world",
		},
	}

	actual := map[string]*Synset{}
	err := parseDataFile(strings.NewReader(testData), "v", actual)
	if err != nil {
		t.Fatal("Parsing error:", err)
	}
	if !reflect.DeepEqual(expected, actual) {
		t.Error("Non-equal values:")
		t.Error(stringify(expected))
		t.Error(stringify(actual))
	}
}

func TestIndexParser(t *testing.T) {
	expected := map[string]*Lemma{
		"n.yoink": &Lemma{
			[]string{"#", "$"},
			[]string{"n.123", "n.456", "n.789"},
		},
	}

	actual := map[string]*Lemma{}
	err := parseIndexFile(strings.NewReader(testIndex), actual)
	if err != nil {
		t.Fatal("Parsing error:", err)
	}
	if !reflect.DeepEqual(expected, actual) {
		t.Error("Non-equal values:")
		t.Error(stringify(expected))
		t.Error(stringify(actual))
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

func stringify(a interface{}) string {
	j, _ := json.Marshal(a)
	return string(j)
}

var testData = `  copyright line
111 1 v 3 foo 1 bar 3 baz 5 2 ! 123 n 0000 @ 321 a 0102 2 + 4 5 + 6 7 | hello world`

var testIndex = `  copyright line
yoink n 3 2 # $ 3 1 123 456 789`

var testException = `foo bar
baz bla blu
`
