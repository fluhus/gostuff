package bits

import (
	"slices"
	"testing"
)

func TestSet_true(t *testing.T) {
	tests := []struct {
		n    int
		want []byte
	}{
		{0, []byte{1, 0}},
		{1, []byte{2, 0}},
		{2, []byte{4, 0}},
		{3, []byte{8, 0}},
		{4, []byte{16, 0}},
		{5, []byte{32, 0}},
		{6, []byte{64, 0}},
		{7, []byte{128, 0}},
		{8, []byte{0, 1}},
		{9, []byte{0, 2}},
		{10, []byte{0, 4}},
		{11, []byte{0, 8}},
		{12, []byte{0, 16}},
		{13, []byte{0, 32}},
		{14, []byte{0, 64}},
		{15, []byte{0, 128}},
	}
	for _, test := range tests {
		b := []byte{0, 0}
		Set(b, test.n, true)
		if !slices.Equal(b, test.want) {
			t.Errorf("Set(%v,%v,%v)=%v, want %v",
				[]byte{0, 0}, test.n, true, b, test.want)
		}
	}
}

func TestSet_false(t *testing.T) {
	tests := []struct {
		n    int
		want []byte
	}{
		{0, []byte{255 - 1, 255}},
		{1, []byte{255 - 2, 255}},
		{2, []byte{255 - 4, 255}},
		{3, []byte{255 - 8, 255}},
		{4, []byte{255 - 16, 255}},
		{5, []byte{255 - 32, 255}},
		{6, []byte{255 - 64, 255}},
		{7, []byte{255 - 128, 255}},
		{8, []byte{255, 255 - 1}},
		{9, []byte{255, 255 - 2}},
		{10, []byte{255, 255 - 4}},
		{11, []byte{255, 255 - 8}},
		{12, []byte{255, 255 - 16}},
		{13, []byte{255, 255 - 32}},
		{14, []byte{255, 255 - 64}},
		{15, []byte{255, 255 - 128}},
	}
	for _, test := range tests {
		b := []byte{255, 255}
		Set(b, test.n, false)
		if !slices.Equal(b, test.want) {
			t.Errorf("Set(%v,%v,%v)=%v, want %v",
				[]byte{0, 0}, test.n, false, b, test.want)
		}
	}
}

func TestGet(t *testing.T) {
	input := []byte{0b10011001, 0b01100111}
	want := []int{1, 0, 0, 1, 1, 0, 0, 1,
		1, 1, 1, 0, 0, 1, 1, 0}
	for i := range want {
		if got := Get(input, i); got != want[i] {
			t.Errorf("Get(%v,%v)=%v, want %v", input, i, got, want[i])
		}
	}
}

func TestSum(t *testing.T) {
	tests := []struct {
		input []byte
		want  int
	}{
		{nil, 0},
		{[]byte{0}, 0},
		{[]byte{0b01100111}, 5},
		{[]byte{0b10011001, 0b01100111}, 9},
	}
	for _, test := range tests {
		if got := Sum(test.input); got != test.want {
			t.Errorf("Sum(%v)=%v, want %v", test.input, got, test.want)
		}
	}
}

func TestByteOnes(t *testing.T) {
	tests := []struct {
		i    int
		want []int
	}{
		{0, nil},
		{1, []int{0}},
		{2, []int{1}},
		{3, []int{0, 1}},
		{4, []int{2}},
		{5, []int{0, 2}},
		{255, []int{0, 1, 2, 3, 4, 5, 6, 7}},
	}
	for _, test := range tests {
		if got := byteOnes[test.i]; !slices.Equal(got, test.want) {
			t.Errorf("byteOnes(%v)=%v, want %v", test.i, got, test.want)
		}
	}
}

func TestOnesZeros(t *testing.T) {
	input := []byte{0b01011100, 0b11101010}
	wantOnes := []int{2, 3, 4, 6, 9, 11, 13, 14, 15}
	wantZeros := []int{0, 1, 5, 7, 8, 10, 12}

	gotOnes := slices.Collect(Ones(input))
	if !slices.Equal(gotOnes, wantOnes) {
		t.Errorf("Ones(%v)=%v, want %v", input, gotOnes, wantOnes)
	}

	gotZeros := slices.Collect(Zeros(input))
	if !slices.Equal(gotZeros, wantZeros) {
		t.Errorf("Zeros(%v)=%v, want %v", input, gotZeros, wantZeros)
	}
}
