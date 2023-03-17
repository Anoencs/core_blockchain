package types

import (
	"fmt"
	"reflect"
)

type List[T any] struct {
	data []T
}

func NewList[T any](data ...T) *List[T] {
	if data != nil {
		return &List[T]{
			data: data,
		}
	}
	return &List[T]{
		data: []T{},
	}
}

func (l *List[T]) GetData() []T {
	return l.data
}

func (l *List[T]) Get(index int) T {
	if index > len(l.data)-1 {
		err := fmt.Sprintf("the given index (%d) is higher than the length (%d)", index, len(l.data))
		panic(err)
	}
	return l.data[index]
}

func (l *List[T]) Insert(v T) {
	l.data = append(l.data, v)
}

func (l *List[T]) Clear() {
	l.data = []T{}
}

func (l *List[T]) GetIndex(v T) int {
	for i := 0; i < len(l.data); i++ {
		if reflect.DeepEqual(l.data[i], v) {
			return i
		}
	}
	return -1
}

func (l *List[T]) Remove(v T) {
	idx := l.GetIndex(v)
	if idx == -1 {
		return
	}
	l.Pop(idx)
}

func (l *List[T]) Pop(idx int) {
	l.data = append(l.data[:idx], l.data[idx+1:]...)
}

func (l *List[T]) Contains(v T) bool {
	for _, val := range l.data {
		if reflect.DeepEqual(val, v) {
			return true
		}
	}
	return false
}

func (l *List[T]) Last() T {
	return l.data[len(l.data)-1]
}

func (l *List[T]) Len() int {
	return len(l.data)
}
