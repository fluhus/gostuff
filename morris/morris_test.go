package morris

import "testing"

func TestRestore(t *testing.T) {
	tests := []struct {
		i    byte
		m    uint
		want uint
	}{
		{0, 1, 0},
		{1, 1, 1},
		{2, 1, 4},
		{3, 1, 9},
		{4, 1, 19},
		{0, 10, 0},
		{1, 10, 1},
		{2, 10, 2},
		{10, 10, 10},
		{15, 10, 20},
		{20, 10, 31},
		{25, 10, 51},
		{30, 10, 72},
	}
	for _, test := range tests {
		if got := Restore(test.i, test.m); got != test.want {
			t.Errorf("Restore(%d,%d)=%d, want %d",
				test.i, test.m, got, test.want)
		}
	}
}

func TestRaise_overflow(t *testing.T) {
	if !checkOverFlow {
		t.Skip()
	}
	defer func() { recover() }()
	got := Raise(byte(255), 10)
	t.Fatalf("Raise(byte(255)=%d, want fail", got)
}

func TestRaise(t *testing.T) {
	const reps = 1000
	want := map[byte]int{
		1: 1, 2: 2, 3: 3, 4: 4, 5: 5, 6: 6, 7: 7, 8: 8, 9: 9, 10: 10,
		15: 20, 20: 30, 25: 50, 30: 70, 35: 110, 40: 150,
	}
	margins := map[byte]int{
		1: 0, 2: 0, 3: 0, 4: 0, 5: 0, 6: 0, 7: 0, 8: 0, 9: 0, 10: 0,
		15: 2, 20: 2, 25: 2, 30: 2, 35: 2, 40: 2,
	}
	got := map[byte]int{}

	for i := 0; i < reps; i++ {
		a, n := byte(0), 0
		for a < 40 {
			n++
			aa := Raise(a, 10)
			if aa == a {
				continue
			}
			a = aa
			if _, ok := want[a]; !ok {
				continue
			}
			got[a] += n
		}
	}

	for k, n := range got {
		got := n / reps
		if got < want[k]-margins[k] || got > want[k]+margins[k] {
			t.Errorf("Raise() to %d took %d, want %d-%d", k, got,
				want[k]-margins[k], want[k]+margins[k])
		}
	}
}
