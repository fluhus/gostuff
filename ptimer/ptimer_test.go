package ptimer

import (
	"bytes"
	"fmt"
	"io"
	"regexp"
	"slices"
	"testing"
	"time"
)

const timePattern = "\\d\\d:\\d\\d:\\d\\d\\.\\d\\d\\d\\d\\d\\d"

func TestNew(t *testing.T) {
	want := "^"
	for _, i := range []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 20, 30, 35} {
		want += fmt.Sprintf("\r%s \\(%s\\) %d", timePattern, timePattern, i)
	}
	want += "\n$"

	got := bytes.NewBuffer(nil)
	pt := New()
	pt.W = got
	for i := 0; i < 35; i++ {
		pt.Inc()
	}
	pt.Done()

	if match, _ := regexp.MatchString(want, got.String()); !match {
		t.Fatalf("Inc()+Done()=%q, want %q", got.String(), want)
	}
}

func TestNewMessage(t *testing.T) {
	want := "^"
	for _, i := range []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 20, 30, 35} {
		want += fmt.Sprintf("\r%s \\(%s\\) hey %d ho",
			timePattern, timePattern, i)
	}
	want += "\n$"

	got := bytes.NewBuffer(nil)
	pt := NewMessage("hey {} ho")
	pt.W = got
	for i := 0; i < 35; i++ {
		pt.Inc()
	}
	pt.Done()

	if match, _ := regexp.MatchString(want, got.String()); !match {
		t.Fatalf("Inc()+Done()=%q, want %q", got.String(), want)
	}
}

func TestNewFunc(t *testing.T) {
	want := "^"
	for _, i := range []float64{1.5, 2.5, 3.5, 4.5, 5.5, 6.5, 7.5, 8.5, 9.5,
		10.5, 20.5, 30.5, 35.5} {
		want += fmt.Sprintf("\r%s \\(%s\\) ho ho %f",
			timePattern, timePattern, i)
	}
	want += "\n$"

	got := bytes.NewBuffer(nil)
	pt := NewFunc(func(i int) string {
		return fmt.Sprintf("ho ho %f", float64(i)+0.5)
	})
	pt.W = got
	for i := 0; i < 35; i++ {
		pt.Inc()
	}
	pt.Done()

	if match, _ := regexp.MatchString(want, got.String()); !match {
		t.Fatalf("Inc()+Done()=%q, want %q", got.String(), want)
	}
}

func TestDone(t *testing.T) {
	want := "^" + timePattern + " hello\n$"

	got := bytes.NewBuffer(nil)
	pt := NewMessage("hello")
	pt.W = got
	pt.Done()

	if match, _ := regexp.MatchString(want, got.String()); !match {
		t.Fatalf("Done()=%q, want %q", got.String(), want)
	}
}

func TestNextCheckpoint(t *testing.T) {
	want := []int{
		1, 2, 3, 4, 5, 6, 7, 8, 9,
		10, 20, 30, 40, 50, 60, 70, 80, 90,
		100, 200, 300, 400, 500}
	var got []int
	i := 0
	for range want {
		i = nextCheckpoint(i)
		got = append(got, i)
	}
	if !slices.Equal(got, want) {
		t.Fatalf("nextCheckpoint(...)=%v, want %v", got, want)
	}
}

func BenchmarkTimer_inc(b *testing.B) {
	pt := New()
	pt.W = io.Discard
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		pt.Inc()
	}
}

func Example() {
	pt := New()
	for i := 0; i < 45; i++ {
		time.Sleep(100 * time.Millisecond)
		pt.Inc()
	}
	pt.Done()
}
