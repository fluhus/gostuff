package ppln

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/fluhus/gostuff/gnum"
)

func TestSerial(t *testing.T) {
	for _, nt := range []int{1, 2, 4, 8} {
		t.Run(fmt.Sprint(nt), func(t *testing.T) {
			n := nt * 100
			var result []int
			Serial(nt, func(push func(int), s Stopper) {
				for i := 0; i < n; i++ {
					push(i)
				}
			}, func(a int, i int, g int, s Stopper) int {
				time.Sleep(time.Millisecond * time.Duration(rand.Intn(5)))
				return a * a
			}, func(i int, s Stopper) {
				result = append(result, i)
			})
			for i := range result {
				if result[i] != i*i {
					t.Errorf("result[%d]=%d, want %d", i, result[i], i*i)
				}
			}
		})
	}
}

func ExampleSerial() {
	ngoroutines := 4
	var results []float64

	Serial(
		ngoroutines,
		// Read/generate input data.
		func(push func(int), s Stopper) {
			for i := 1; i <= 100; i++ {
				push(i)
			}
		},
		// Some processing.
		func(a int, i, g int, s Stopper) float64 {
			return float64(a*a) + 0.5
		},
		// Accumulate/forward outputs.
		func(a float64, s Stopper) {
			results = append(results, a)
		})

	fmt.Println(results[:3], results[len(results)-3:])

	// Output:
	// [1.5 4.5 9.5] [9604.5 9801.5 10000.5]
}

func ExampleSerial_parallelAggregation() {
	ngoroutines := 4
	results := make([]int, ngoroutines) // Goroutine-specific data and objects.

	Serial(
		ngoroutines,
		// Read/generate input data.
		func(push func(int), s Stopper) {
			for i := 1; i <= 100; i++ {
				push(i)
			}
		},
		// Accumulate in goroutine-specific memory.
		func(a int, i, g int, s Stopper) int {
			results[g] += a
			return 0 // Unused.
		},
		// No outputs.
		func(a int, s Stopper) {})

	// Collect the results of all goroutines.
	fmt.Println("Sum of 1-100:", gnum.Sum(results))

	// Output:
	// Sum of 1-100: 5050
}
