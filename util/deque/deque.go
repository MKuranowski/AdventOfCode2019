package deque

import "container/list"

type Deque[T any] struct {
	l *list.List
}

func NewDeque[T any]() Deque[T]  { return Deque[T]{&list.List{}} }
func (d Deque[T]) Len() int      { return d.l.Len() }
func (d Deque[T]) PeekFront() T  { return d.l.Front().Value.(T) }
func (d Deque[T]) PeekBack() T   { return d.l.Back().Value.(T) }
func (d Deque[T]) PushFront(x T) { d.l.PushFront(x) }
func (d Deque[T]) PushBack(x T)  { d.l.PushBack(x) }
func (d Deque[T]) PopFront() T   { return d.l.Remove(d.l.Front()).(T) }
func (d Deque[T]) PopBack() T    { return d.l.Remove(d.l.Back()).(T) }
