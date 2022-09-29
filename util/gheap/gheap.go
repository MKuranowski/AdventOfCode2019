package gheap

import (
	"container/heap"

	"golang.org/x/exp/constraints"
)

type heapImpl[T any] struct {
	arr  []T
	less func(T, T) bool
}

func (h heapImpl[T]) Len() int           { return len(h.arr) }
func (h heapImpl[T]) Less(i, j int) bool { return h.less(h.arr[i], h.arr[j]) }
func (h heapImpl[T]) Swap(i, j int)      { h.arr[i], h.arr[j] = h.arr[j], h.arr[i] }
func (h *heapImpl[T]) Push(x any)        { h.arr = append(h.arr, x.(T)) }

func (h *heapImpl[T]) Pop() any {
	var zero T
	l := len(h.arr)
	r := h.arr[l-1]
	h.arr[l-1] = zero
	h.arr = h.arr[:l-1]
	return r
}

func DefaultLess[T constraints.Ordered]() func(T, T) bool {
	return func(t1, t2 T) bool { return t1 < t2 }
}

func GetInterface[T any](less func(T, T) bool) heap.Interface {
	return &heapImpl[T]{nil, less}
}

type GenericHeap[T any] struct {
	i heap.Interface
}

func NewGenericHeap[T any](less func(T, T) bool) GenericHeap[T] {
	return GenericHeap[T]{GetInterface(less)}
}

func NewGenericHeapWithInterface[T any](i heap.Interface) GenericHeap[T] {
	return GenericHeap[T]{i}
}

func (h GenericHeap[T]) Len() int { return h.i.Len() }
func (h GenericHeap[T]) Pop() T   { return heap.Pop(h.i).(T) }
func (h GenericHeap[T]) Push(x T) { heap.Push(h.i, x) }
