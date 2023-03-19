package types

var STACK_SIZE = 1024

type Stack struct {
	data []any
	sp   int
}

func NewStack(size int) *Stack {
	return &Stack{
		data: make([]any, size),
		sp:   0,
	}
}

func (s *Stack) Top() any {
	return s.data[0]
}

func (s *Stack) Pop() any {
	value := s.data[s.sp-1]
	s.data = append(s.data[:s.sp-1], s.data[s.sp+1:]...)
	s.sp--
	return value
}

func (s *Stack) Get() []any {
	return s.data
}

func (s *Stack) Push(v any) {
	s.data[s.sp] = v
	s.sp++
}

func (s *Stack) IsEmpty() bool {
	return s.sp <= 0
}

func (s *Stack) Clear() {
	s.sp = 0
	s.data = []any{}
}
