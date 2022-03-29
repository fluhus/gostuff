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
