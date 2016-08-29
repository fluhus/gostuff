package nlock

import (
	"testing"
	"time"
)

func TestNLock(t *testing.T) {
	lock := New(3)
	expected := []bool{true, true, true, false, false}
	for _, e := range expected {
		if ok := lock.TryLock(); ok != e {
			t.Fatalf("TryLock() = %v, want %v", ok, e)
		}
	}
	lock.Unlock()
	lock.Unlock()
	lock.Unlock()

	defer func() {
		if recover() == nil {
			t.Fatalf("Unlock() returned normally, want panic")
		}
	}()
	lock.Unlock()
}

func TestNLock_async(t *testing.T) {
	x := 0
	lock := New(3)
	go func() {
		for i := 0; i < 5; i++ {
			lock.Lock()
			x++
		}
	}()
	expected := []int{3, 4, 5, 5, 5}
	for _, e := range expected {
		time.Sleep(100 * time.Millisecond)
		if x != e {
			t.Fatalf("Lock() -> x = %v, want %v", x, e)
		}
		lock.Unlock()
	}
}

func TestNLock_badN(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatalf("New(0) returned normally, want panic")
		}
	}()
	New(0)
}
