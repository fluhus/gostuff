package ppln

// A Stopper is used in pipelines to communicate that processing should stop.
type Stopper chan struct{}

// Stop sets Stopped to true. After calling this function no further calls to
// mapper and puller will be made. Pusher should stop itself if stopped.
func (s Stopper) Stop() {
	close(s)
}

// Stopped returns whether Stop was called.
func (s Stopper) Stopped() bool {
	select {
	case <-s:
		return true
	default:
		return false
	}
}
