package vm

import "fmt"

type RuntimeError struct {
	Message string
}

func (e RuntimeError) Error() string {
	return fmt.Sprintf("💥 RuntimeError: %s", e.Message)
}
