package main

import "fmt"

type Stack struct {
	items []int
}

func (s *Stack) Push(item int) {
	s.items = append(s.items, item)
}

func (s *Stack) Pop() {
	if len(s.items) == 0 {
		return
	}

	s.items = s.items[:len(s.items)-1]
}

func (s *Stack) Top() (int, error) {
	if len(s.items) == 0 {
		return 0, fmt.Errorf("stack underflow")
	}
	return s.items[len(s.items)-1], nil
}
