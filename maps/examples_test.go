package maps

import (
	"fmt"
)

func ExampleKeys() {
	m := map[string]int{
		"a": 1,
		"c": 3,
		"b": 2,
	}
	keys := Keys(m).([]string)
	fmt.Println(keys)
	// Output: [a b c]
}

func ExampleOf() {
	people := []string{"alice", "bob"}
	m := Of(people, true).(map[string]bool)

	fmt.Println("'alice' in map:", m["alice"])
	fmt.Println("'bob' in map:", m["bob"])
	fmt.Println("'charles' in map:", m["charles"])
	// Output:
	// 'alice' in map: true
	// 'bob' in map: true
	// 'charles' in map: false
}
