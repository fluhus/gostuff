package gnum

import (
	"math"
	"slices"
	"testing"
)

func TestEntropy(t *testing.T) {
	input1 := []int{1, 2, 3, 4}
	input2 := []uint{1, 2, 3, 4}
	input3 := []float32{1, 2, 3, 4}
	input4 := []float64{1, 2, 3, 4}
	want := 1.8464393
	if got := Entropy(input1); Diff(got, want) > 0.00000005 {
		t.Errorf("Entropy(%v)=%v, want %v", input1, got, want)
	}
	if got := Entropy(input2); Diff(got, want) > 0.00000005 {
		t.Errorf("Entropy(%v)=%v, want %v", input2, got, want)
	}
	if got := Entropy(input3); Diff(got, want) > 0.00000005 {
		t.Errorf("Entropy(%v)=%v, want %v", input3, got, want)
	}
	if got := Entropy(input4); Diff(got, want) > 0.00000005 {
		t.Errorf("Entropy(%v)=%v, want %v", input4, got, want)
	}
}

func TestIdiv(t *testing.T) {
	tests := []struct {
		a, b, want int
	}{
		{8, 1, 8},
		{8, 2, 4},
		{8, 3, 3},
		{8, 4, 2},
		{8, 5, 2},
		{8, 6, 1},
		{8, 7, 1},
		{8, 8, 1},
		{8, 9, 1},
		{8, 10, 1},
		{8, 11, 1},
		{8, 12, 1},
		{8, 13, 1},
		{8, 14, 1},
		{8, 15, 1},
		{8, 17, 0},
		{8, 18, 0},
		{8, 19, 0},
		{8, 20, 0},
	}
	for _, test := range tests {
		if got := Idiv(test.a, test.b); got != test.want {
			t.Errorf("Idiv(%v,%v)=%v, want %v", test.a, test.b, got, test.want)
		}
	}
}

func TestMinMax(t *testing.T) {
	tests := []struct {
		input            []int
		mn, mx, amn, amx int
	}{
		{nil, 0, 0, -1, -1},
		{[]int{42}, 42, 42, 0, 0},
		{[]int{42, 42}, 42, 42, 0, 0},
		{[]int{42, 42, 42}, 42, 42, 0, 0},
		{[]int{1, 2, 3}, 1, 3, 0, 2},
		{[]int{1, 3, 2}, 1, 3, 0, 1},
		{[]int{2, 1, 3}, 1, 3, 1, 2},
		{[]int{2, 3, 1}, 1, 3, 2, 1},
		{[]int{3, 1, 2}, 1, 3, 1, 0},
		{[]int{3, 2, 1}, 1, 3, 2, 0},
	}
	for _, test := range tests {
		mn, mx, amn, amx := Min(test.input), Max(test.input), ArgMin(test.input), ArgMax(test.input)
		if mn != test.mn {
			t.Errorf("Min(%v)=%v, want %v", test.input, mn, test.mn)
		}
		if mx != test.mx {
			t.Errorf("Max(%v)=%v, want %v", test.input, mx, test.mx)
		}
		if amn != test.amn {
			t.Errorf("ArgMin(%v)=%v, want %v", test.input, amn, test.amn)
		}
		if amx != test.amx {
			t.Errorf("ArgMax(%v)=%v, want %v", test.input, amx, test.amx)
		}
	}
}

func TestSum(t *testing.T) {
	tests := []struct {
		input []int
		want  int
	}{
		{nil, 0},
		{[]int{1}, 1},
		{[]int{1, 1}, 2},
		{[]int{1, 1, 1, 1}, 4},
		{[]int{6, 4, 1}, 11},
	}
	for _, test := range tests {
		if got := Sum(test.input); got != test.want {
			t.Errorf("Sum(%v)=%v, want %v", test.input, got, test.want)
		}
	}
}

func TestMean(t *testing.T) {
	tests := []struct {
		input []int
		want  float64
	}{
		{[]int{1}, 1},
		{[]int{1, 1}, 1},
		{[]int{1, 1, 1, 1}, 1},
		{[]int{6, 4, -1}, 3},
	}
	for _, test := range tests {
		if got := Mean(test.input); got != test.want {
			t.Errorf("Mean(%v)=%v, want %v", test.input, got, test.want)
		}
	}
}

func TestExpMean(t *testing.T) {
	tests := []struct {
		input []int
		want  float64
	}{
		{[]int{1}, 1},
		{[]int{1, 1}, 1},
		{[]int{3, 3, 3, 3}, 3},
		{[]int{10, 1000}, 100},
		{[]int{10, 100}, math.Sqrt(1000)},
		{[]int{10, 100, 1000}, 100},
	}
	const tolerance = 0.0000001
	for _, test := range tests {
		if got := ExpMean(test.input); Diff(got, test.want) > tolerance {
			t.Errorf("ExpMean(%v)=%v, want %v", test.input, got, test.want)
		}
	}
}

func FuzzSumMean(f *testing.F) {
	f.Add(0.0, 0.0, 0.0, 0.0)
	f.Fuzz(func(t *testing.T, a float64, b float64, c float64, d float64) {
		slice := []float64{a, b, c, d}
		want := a + b + c + d
		if got := Sum(slice); got != want {
			t.Fatalf("Sum([%v,%v,%v,%v])=%v, want %v", a, b, c, d, got, want)
		}
		want /= 4
		if got := Mean(slice); got != want {
			t.Fatalf("Mean([%v,%v,%v,%v])=%v, want %v", a, b, c, d, got, want)
		}
		if a > 0 && b > 0 && c > 0 && d > 0 {
			const tol = 0.0000001
			want = math.Pow(a*b*c*d, 0.25)
			if got := ExpMean(slice); Diff(got, want) > tol {
				t.Fatalf("ExpMean([%v,%v,%v,%v])=%v, want %v", a, b, c, d, got, want)
			}
		}
	})
}

func TestQuantiles(t *testing.T) {
	input := []int{1, 1, 2, 3, 5, 8, 13, 21, 34, 55, 89, 144}
	q := []float64{0.3, 0.5, 0, 1, 0.9}
	want1 := []int{3, 8, 1, 144, 89}
	want2 := []int{3, 13, 1, 144, 89}
	got := Quantiles(input, q...)
	if !slices.Equal(got, want1) && !slices.Equal(got, want2) {
		t.Fatalf("Quantiles(%v,%v)=%v, want %v or %v",
			input, q, got, want1, want2)
	}
}

func TestNQuantiles(t *testing.T) {
	input := []int{1, 1, 2, 3, 5, 8, 13, 21, 34, 55, 89, 144}
	want := []int{1, 5, 21, 144}
	got := NQuantiles(input, 3)
	if !slices.Equal(got, want) {
		t.Fatalf("Quantiles(%v,3)=%v, want %v", input, got, want)
	}
}

func TestLogFactorial(t *testing.T) {
	tests := [][]int{
		{0, 1}, {1, 1}, {2, 2}, {3, 6}, {4, 24}, {5, 120}, {6, 720},
	}
	for _, test := range tests {
		got := math.Exp(LogFactorial(test[0]))
		want := float64(test[1])
		if Diff(got, want) > want*0.05 {
			t.Errorf("lf(%v)=%v, want %v", test[0], got, test[1])
		}
	}
}

func BenchmarkLogFactorial(b *testing.B) {
	for i := range b.N {
		LogFactorial(i)
	}
}
