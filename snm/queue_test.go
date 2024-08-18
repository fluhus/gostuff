package snm

import "testing"

func TestQueue(t *testing.T) {
	q := &Queue[int]{}

	qEnqueue(q, 11)
	qExpect(q, t, 11)

	qEnqueue(q, 22, 33)
	qExpect(q, t, 22, 33)

	qEnqueue(q, 44, 55, 66)
	qExpect(q, t, 44)
	qEnqueue(q, 77)
	qExpect(q, t, 55)
	qEnqueue(q, 88)
	qExpect(q, t, 66, 77, 88)
}

func qEnqueue(q *Queue[int], x ...int) {
	for _, xx := range x {
		q.Enqueue(xx)
	}
}

func qExpect(q *Queue[int], t *testing.T, x ...int) {
	for _, xx := range x {
		if got := q.Peek(); got != xx {
			t.Fatalf("q.pull()=%v, want %v", got, xx)
		}
		if got := q.Dequeue(); got != xx {
			t.Fatalf("q.pull()=%v, want %v", got, xx)
		}
	}
}

func FuzzQueue(f *testing.F) {
	f.Add(1, 1, 1, 1, 1, 1, 1, 1, 1, 1)
	f.Fuzz(func(t *testing.T, a, b, c, d, e, f, g, h, i, j int) {
		q := &Queue[int]{}
		qEnqueue(q, a)
		qExpect(q, t, a)
		qEnqueue(q, b, c)
		qExpect(q, t, b, c)
		qEnqueue(q, d, e)
		qExpect(q, t, d)
		qEnqueue(q, f)
		qExpect(q, t, e)
		qEnqueue(q, g, h, i)
		qExpect(q, t, f, g)
		qEnqueue(q, j)
		qExpect(q, t, h, i, j)
	})
}
