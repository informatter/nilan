package interpreter

import "fmt"

// Defines the struct for all runtime errors in the Parser
type RuntimeError struct {
	Line    int32
	Column  int
	Message string
}

func CreateRuntimeError(line int32, column int, message string) RuntimeError {
	return RuntimeError{
		Line:    line,
		Column:  column,
		Message: message,
	}
}

func (e RuntimeError) Error() string {
	return fmt.Sprintf("💥 Nilan Runtime error:\nline:%d, column:%d - %s", e.Line, e.Column, e.Message)
}
