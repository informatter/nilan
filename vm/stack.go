package vm

type Stack []any

// Check if the stack is empty
func (s *Stack) IsEmpty() bool {
	return len(*s) == 0
}

// Push a new value onto the stack
func (s *Stack) Push(value any) {
	*s = append(*s, value)
}

// Removes and returns the top element of the stack
func (s *Stack) Pop() (any, bool) {
	if s.IsEmpty() {
		return nil, false
	}
	index := len(*s) - 1
	element := (*s)[index]
	*s = (*s)[:index]
	return element, true
}

// Returns the top element without removing it
func (s *Stack) Peek() (any, bool) {
	if s.IsEmpty() {
		return nil, false
	}
	index := len(*s) - 1
	return (*s)[index], true
}
