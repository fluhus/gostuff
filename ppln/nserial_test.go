package ppln

import (
	"fmt"
	"math"
	"testing"
)

func TestNonSerial(t *testing.T) {
	want := 21082009.0
	for _, nt := range []int{1, 2, 4, 8} {
		t.Run(fmt.Sprint(nt), func(t *testing.T) {
			got := 0.0
			NonSerial[int, float64](
				nt,
				RangeInput(1, 100001),
				func(a, g int) (float64, error) {
					return math.Sqrt(float64(a)), nil
				},
				func(a float64) error {
					got += a
					return nil
				},
			)
			if math.Round(got) != want {
				t.Fatalf("NonSerial: got %f, want %f", got, want)
			}
		})
	}
}

func TestNonSerial_inputError(t *testing.T) {
	for _, nt := range []int{1, 2, 4, 8} {
		t.Run(fmt.Sprint(nt), func(t *testing.T) {
			got := 0.0
			err := NonSerial[int, float64](
				nt,
				func(yield func(int, error) bool) {
					for i, err := range RangeInput(1, 100001) {
						if i == 1000 {
							yield(0, fmt.Errorf("oh no"))
							return
						}
						if !yield(i, err) {
							return
						}
					}
				},
				func(a, g int) (float64, error) {
					return math.Sqrt(float64(a)), nil
				},
				func(a float64) error {
					got += a
					return nil
				},
			)
			if err == nil {
				t.Fatalf("NonSerial succeeded, want error")
			}
		})
	}
}

func TestNonSerial_transformError(t *testing.T) {
	for _, nt := range []int{1, 2, 4, 8} {
		t.Run(fmt.Sprint(nt), func(t *testing.T) {
			got := 0.0
			err := NonSerial[int, float64](
				nt,
				RangeInput(1, 100001),
				func(a, g int) (float64, error) {
					if a == 1000 {
						return 0, fmt.Errorf("oh no")
					}
					return math.Sqrt(float64(a)), nil
				},
				func(a float64) error {
					got += a
					return nil
				},
			)
			if err == nil {
				t.Fatalf("NonSerial succeeded, want error")
			}
		})
	}
}

func TestNonSerial_outputError(t *testing.T) {
	for _, nt := range []int{1, 2, 4, 8} {
		t.Run(fmt.Sprint(nt), func(t *testing.T) {
			got := 0.0
			err := NonSerial[int, float64](
				nt,
				RangeInput(1, 100001),
				func(a, g int) (float64, error) {
					return math.Sqrt(float64(a)), nil
				},
				func(a float64) error {
					if a == 32 {
						return fmt.Errorf("oh no")
					}
					got += a
					return nil
				},
			)
			if err == nil {
				t.Fatalf("NonSerial succeeded, want error")
			}
		})
	}
}

// TODO(amit): Error tests.
