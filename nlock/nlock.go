// A lock that can be locked n times simultanously.
package nlock

import (
	"fmt"
	"sync"
)

// An NLock is a lock that allows n holders to hold the lock simultanously.
type NLock struct {
	c *sync.Cond // Sync object.
	n int        // Max number of holders.
	i int        // Current number of holders.
}

// Creates a new lock for n maximum holders.
func New(n int) *NLock {
	if n <= 0 {
		panic(fmt.Sprintf("Bad n: %v, needs to be at least 1.", n))
	}
	return &NLock{sync.NewCond(&sync.Mutex{}), n, 0}
}

// Locks the lock. Will block if n calls to lock were made, that were not
// unlocked.
func (n *NLock) Lock() {
	n.c.L.Lock()
	defer n.c.L.Unlock()
	for n.i >= n.n {
		n.c.Wait()
	}
	n.i++
}

// Releases one holder of the lock. Panics if lock has 0 holders.
func (n *NLock) Unlock() {
	n.c.L.Lock()
	defer n.c.L.Unlock()
	if n.i <= 0 {
		panic("Attempt to unlock a not locked lock.")
	}
	n.i--
	n.c.Signal()
}

// Attempts to obtain lock without waiting. Returns true if succeeded, or
// false if not.
func (n *NLock) TryLock() bool {
	n.c.L.Lock()
	defer n.c.L.Unlock()
	if n.i >= n.n {
		return false
	}
	n.i++
	return true
}
