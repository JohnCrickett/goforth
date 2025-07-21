package interpreter

import "fmt"

type Stack[T any] struct {
	items []T
}

func (s *Stack[T]) Push(item T) {
	s.items = append(s.items, item)
}

func (s *Stack[T]) Pop() {
	if len(s.items) == 0 {
		return
	}

	s.items = s.items[:len(s.items)-1]
}

func (s *Stack[T]) Top() (T, error) {
	if len(s.items) == 0 {
		var t T
		return t, fmt.Errorf("stack underflow")
	}
	return s.items[len(s.items)-1], nil
}
