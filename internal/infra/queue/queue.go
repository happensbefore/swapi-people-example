package queue

import (
	"sync"
)

const minCapacity = 16

type Queue[T any] struct {
	buf   []T
	head  int
	tail  int
	count int

	mu sync.Mutex
}

func NewQueue[T any]() *Queue[T] {
	return &Queue[T]{
		buf: make([]T, minCapacity),
	}
}

func (q *Queue[T]) Add(data T) {
	q.mu.Lock()
	defer q.mu.Unlock()

	q.push(data)
}

func (q *Queue[T]) GetAll() []T {
	q.mu.Lock()
	defer q.mu.Unlock()

	if q.count <= 0 {
		return nil
	}

	ret := make([]T, q.count)
	for i := 0; i < q.count; i++ {
		ret[i] = q.pop()
	}

	return ret
}

func (q *Queue[T]) Get(n int) []T {
	q.mu.Lock()
	defer q.mu.Unlock()

	if q.count <= 0 || q.count < n {
		return nil
	}

	ret := make([]T, n)
	for i := 0; i < n; i++ {
		ret[i] = q.pop()
	}

	return ret
}

func (q *Queue[T]) push(elem T) {
	q.growIfNeed()

	q.buf[q.tail] = elem
	q.tail = q.next(q.tail)
	q.count++
}

func (q *Queue[T]) pop() T {
	ret := q.buf[q.head]

	var empty T
	q.buf[q.head] = empty

	q.head = q.next(q.head)
	q.count--

	q.shrinkIfNeed()

	return ret
}

func (q *Queue[T]) growIfNeed() {
	if q.count != len(q.buf) {
		return
	}

	if len(q.buf) == 0 {
		q.buf = make([]T, minCapacity)

		return
	}

	q.resize()
}

func (q *Queue[T]) shrinkIfNeed() {
	if len(q.buf) > minCapacity && (q.count<<2) == len(q.buf) {
		q.resize()
	}
}

func (q *Queue[T]) resize() {
	newBuf := make([]T, q.count<<1)
	if q.tail > q.head {
		copy(newBuf, q.buf[q.head:q.tail])
	} else {
		n := copy(newBuf, q.buf[q.head:])
		copy(newBuf[n:], q.buf[:q.tail])
	}

	q.head = 0
	q.tail = q.count
	q.buf = newBuf
}

func (q *Queue[T]) next(i int) int {
	return (i + 1) & (len(q.buf) - 1)
}
