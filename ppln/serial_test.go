package ppln

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/fluhus/gostuff/gnum"
)

func ExampleSerial() {
	ngoroutines := 4
	var results []float64

	Serial(
		ngoroutines,
		// Read/generate input data.
		func(push func(int), stop func() bool) error {
			for i := 1; i <= 100; i++ {
				push(i)
			}
			return nil
		},
		// Some processing.
		func(a int, i, g int) (float64, error) {
			return float64(a*a) + 0.5, nil
		},
		// Accumulate/forward outputs.
		func(a float64) error {
			results = append(results, a)
			return nil
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
		func(push func(int), stop func() bool) error {
			for i := 1; i <= 100; i++ {
				push(i)
			}
			return nil
		},
		// Accumulate in goroutine-specific memory.
		func(a int, i, g int) (int, error) {
			results[g] += a
			return 0, nil // Unused.
		},
		// No outputs.
		func(a int) error { return nil })

	// Collect the results of all goroutines.
	fmt.Println("Sum of 1-100:", gnum.Sum(results))

	// Output:
	// Sum of 1-100: 5050
}

func TestSerial(t *testing.T) {
	for _, nt := range []int{1, 2, 4, 8} {
		t.Run(fmt.Sprint(nt), func(t *testing.T) {
			n := nt * 100
			var result []int
			err := Serial(nt, func(push func(int), stop func() bool) error {
				for i := 0; i < n; i++ {
					push(i)
				}
				return nil
			}, func(a int, i int, g int) (int, error) {
				time.Sleep(time.Millisecond * time.Duration(rand.Intn(3)))
				return a * a, nil
			}, func(i int) error {
				result = append(result, i)
				return nil
			})
			if err != nil {
				t.Fatalf("Serial2(...) failed: %d", err)
			}
			for i := range result {
				if result[i] != i*i {
					t.Errorf("result[%d]=%d, want %d", i, result[i], i*i)
				}
			}
		})
	}
}

func TestSerial_error(t *testing.T) {
	for _, nt := range []int{1, 2, 4, 8} {
		t.Run(fmt.Sprint(nt), func(t *testing.T) {
			n := nt * 100
			var result []int
			err := Serial(nt, func(push func(int), stop func() bool) error {
				for i := 0; i < n; i++ {
					if stop() {
						break
					}
					push(i)
				}
				return nil
			}, func(a int, i int, g int) (int, error) {
				time.Sleep(time.Millisecond * time.Duration(rand.Intn(3)))
				if a > 300 {
					return 0, fmt.Errorf("a too big: %d", a)
				}
				return a * a, nil
			}, func(i int) error {
				result = append(result, i)
				return nil
			})
			if nt <= 3 {
				if err != nil {
					t.Fatalf("Serial2(...) failed: %d", err)
				}
				for i := range result {
					if result[i] != i*i {
						t.Errorf("result[%d]=%d, want %d", i, result[i], i*i)
					}
				}
			} else { // n > 3
				if err == nil {
					t.Fatalf("Serial2(...) succeeded, want error")
				}
			}
		})
	}
}
