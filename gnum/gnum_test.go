package gnum

import "testing"

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
