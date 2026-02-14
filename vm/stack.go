package vm

type Stack []any

func (s *Stack) IsEmpty() bool {
	return len(*s) == 0
}

func (s *Stack) Push(value any) {
	*s = append(*s, value)
}

// Removes and returns the top element of the stack
func (s *Stack) Pop() any {
	if s.IsEmpty() {
		return nil
	}
	index := len(*s) - 1
	element := (*s)[index]
	*s = (*s)[:index]
	return element
}

// Returns the top element without removing it
func (s *Stack) Peek() any {
	if s.IsEmpty() {
		return nil
	}
	index := len(*s) - 1
	return (*s)[index]
}
