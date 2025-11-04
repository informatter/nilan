package compiler

import "fmt"

type SyntaxError struct {
	Message string
}

func (e SyntaxError) Error() string {
	return fmt.Sprintf("💥 SyntaxError: %s", e.Message)
}
