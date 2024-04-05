package containers

type Stack[T any] struct {
	items []T
}

func (s *Stack[T]) Push(item T) int {
	s.items = append(s.items, item)
	return s.Size()
}

func (s *Stack[T]) Pop() (popped T) {
	if s.Empty() {
		return
	}
	s.items, popped = s.items[0:s.Size()-1], s.items[s.Size()-1]
	return popped
}

func (s *Stack[T]) Peek() (peeked T) {
	if s.Empty() {
		return
	}
	return s.items[s.Size()-1]
}

func (s *Stack[T]) Size() int {
	return len(s.items)
}

func (s *Stack[T]) Empty() bool {
	return s.Size() == 0
}
