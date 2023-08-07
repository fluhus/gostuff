// Package ptimer provides a progress timer for iterative processes.
//
// A timer prints how much time passed since its creation at exponentially
// growing time-points.
// Percisely, prints are triggered after i calls to Inc, if i has only one non-zero
// digit. That is: 1, 2, 3 .. 9, 10, 20, 30 .. 90, 100, 200, 300...
//
// # Output Format
//
// For a regular use:
//
//	00:00:00.000000 (00:00:00.000000) message
//	|                          |         |
//	Total time since creation  |         |
//	                           |         |
//	Average time per call to Inc         |
//	                                     |
//	User-defined message ----------------|
//	(default message is number of calls to Inc)
//
// When calling Done without calling Inc:
//
//	00:00:00.000000 message
package ptimer

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

// A Timer measures time during iterative processes and prints the progress on
// exponential checkpoints.
type Timer struct {
	N int              // Current count, incremented with each call to Inc
	W io.Writer        // Timer's output, defaults to stderr
	t time.Time        // Start time
	f func(int) string // Message function
	c int              // Next checkpoint
}

// Prints the progress.
func (t *Timer) print() {
	since := time.Since(t.t)
	if t.N == 0 { // Happens when calling Done without Inc.
		fmt.Fprintf(t.W, "%s %s", fmtDuration(since), t.f(t.N))
		return
	}
	fmt.Fprintf(t.W, "\r%s (%s) %s", fmtDuration(since),
		fmtDuration(since/time.Duration(t.N)), t.f(t.N))
}

// Formats a duration in constant-width format.
func fmtDuration(d time.Duration) string {
	return fmt.Sprintf("%02d:%02d:%02d.%06d",
		d/time.Hour,
		d%time.Hour/time.Minute,
		d%time.Minute/time.Second,
		d%time.Second/time.Microsecond,
	)
}

// NewFunc returns a new timer that calls f with the current count on checkpoints,
// and prints its output.
func NewFunc(f func(i int) string) *Timer {
	return &Timer{0, os.Stderr, time.Now(), f, 0}
}

// NewMessasge returns a new timer that prints msg on checkpoints.
// A "{}" in msg will be replaced with the current count.
func NewMessasge(msg string) *Timer {
	return NewFunc(func(i int) string {
		return strings.ReplaceAll(msg, "{}", fmt.Sprint(i))
	})
}

// New returns a new timer that prints the current count on checkpoints.
func New() *Timer {
	return NewMessasge("{}")
}

// Inc increments t's counter and prints progress if reached a checkpoint.
func (t *Timer) Inc() {
	t.N++
	if t.N >= t.c {
		t.print()
		for t.c <= t.N {
			t.c = nextCheckpoint(t.c)
		}
	}
}

// Done prints progress as if a checkpoint was reached.
func (t *Timer) Done() {
	t.print()
	fmt.Fprintln(t.W)
}

// Returns the next i in which the timer should print, given the current i.
func nextCheckpoint(i int) int {
	m := 1
	for m*10 <= i {
		m *= 10
	}
	if i%m != 0 {
		panic(fmt.Sprintf(
			"bad checkpoint: %d, should be a multiple of a power of 10", i))
	}
	return i + m
}
