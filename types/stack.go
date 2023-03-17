package types

import "fmt"

var STACK_SIZE = 1024

type Stack[T any] struct {
	data []T
	size int
	cap  int
}

func NewStack[T any](data ...T) *Stack[T] {
	if data == nil {
		return &Stack[T]{
			data: data,
			size: 0,
			cap:  STACK_SIZE,
		}
	}
	return &Stack[T]{
		data: data,
		size: len(data),
		cap:  STACK_SIZE,
	}
}

func (s *Stack[T]) Top() T {
	if s.IsEmpty() {
		err := fmt.Errorf("stack underflow")
		panic(err)
	}

	return s.data[s.size-1]
}

func (s *Stack[T]) Pop() {
	if s.IsEmpty() {
		err := fmt.Errorf("stack underflow")
		panic(err)
	}
	s.data = s.data[:s.size-1]
	s.size--
}

func (s *Stack[T]) Push(v T) {
	if s.size == s.cap {
		err := fmt.Errorf("stack overflow")
		panic(err)
	}

	s.data = append(s.data, v)
	s.size++
}

func (s *Stack[T]) IsEmpty() bool {
	return s.size <= 0
}

func (s *Stack[T]) Clear() {
	s.size = 0
	s.data = []T{}
}
