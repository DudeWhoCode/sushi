package repl

// Stack implementation used for basic brackets/linefeed handling in REPL

import "container/list"

type stack struct {
	count int
	items *list.List
}

func NewStack() *stack {
	return &stack{
		items: list.New(),
	}
}

func (s *stack) push(v byte) {
	s.items.PushBack(v)
	s.count++
}

func (s *stack) pop() byte {
	item := s.items.Front()
	if item == nil {
		return 0
	}
	s.items.Remove(item)
	s.count--
	return item.Value.(byte)
}

func (s *stack) peek() byte {
	item := s.items.Front()
	if item != nil {
		return item.Value.(byte)
	}
	return 0
}
