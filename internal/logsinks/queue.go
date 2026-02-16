package logsinks

import "sync"

type Queue[T any] struct {
	ch   chan T
	once sync.Once
}

func NewQueue[T any](size int) *Queue[T] {
	return &Queue[T]{ch: make(chan T, size)}
}

func (q *Queue[T]) Push(v T) bool {
	select {
	case q.ch <- v:
		return true
	default:
		select {
		case <-q.ch:
		default:
		}

		select {
		case q.ch <- v:
			return true
		default:
			return false
		}
	}
}

func (q *Queue[T]) Chan() <-chan T {
	return q.ch
}

func (q *Queue[T]) Close() {
	q.once.Do(func() {
		close(q.ch)
	})
}
