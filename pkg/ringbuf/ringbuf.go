package ringbuf

import (
	"fmt"

	"github.com/gammazero/deque"
)

// if accessed from different go routines, lock access via mutex!

type Ringbuf[T any] struct {
	q    *deque.Deque[T]
	size int
}

func New[T any](size int) *Ringbuf[T] {
	return &Ringbuf[T]{
		q:    deque.New[T](size),
		size: size,
	}
}

func (r *Ringbuf[T]) Push(v T) {
	r.q.PushBack(v)
	if r.q.Len() > r.size {
		r.q.PopFront()
	}
}

// At returns the element at index i.
// index 0 is the oldest element.
// accessing an element > Len()-1 panics
func (r *Ringbuf[T]) At(i int) T {
	return r.q.At(i)
}

func (r *Ringbuf[T]) Len() int {
	return r.q.Len()
}

func (r *Ringbuf[T]) String(s string) string {
	rv := fmt.Sprintf("%s: %d\n", s, r.Len())
	for i := 0; i < r.Len(); i++ {
		rv += fmt.Sprintf("%d:%+v\n", i, r.At(i))
	}
	return rv
}
