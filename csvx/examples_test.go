package csvx

import (
	"fmt"
	"strings"
)

func ExampleDecodeReader() {
	type person struct {
		Name string
		Age  int
	}
	input := strings.NewReader("alice,30\nbob,25")

	for p, err := range DecodeReader[person](input) {
		if err != nil {
			panic(err)
		}
		fmt.Println(p.Name, "is", p.Age, "years old")
	}

	//Output:
	//alice is 30 years old
	//bob is 25 years old
}

func ExampleDecodeReaderHeader() {
	type person struct {
		Name string
		Age  int
	}
	input := strings.NewReader(
		"user_id,age,city,name\n" +
			"111,30,paris,alice\n" +
			"222,25,london,bob")

	for p, err := range DecodeReaderHeader[person](input) {
		if err != nil {
			panic(err)
		}
		fmt.Println(p.Name, "is", p.Age, "years old")
	}

	//Output:
	//alice is 30 years old
	//bob is 25 years old
}
