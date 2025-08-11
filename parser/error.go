package parser

import "fmt"

// Defines the struct for all syntax errors in the Parser
type SyntaxError struct {
	Line    int32
	Column  int
	Message string
}

func CreateSyntaxError(line int32, column int, message string) SyntaxError {
	return SyntaxError{
		Line:    line,
		Column:  column,
		Message: message,
	}
}

func (e SyntaxError) Error() string {
	return fmt.Sprintf("💥 Nilan Syntax error:\nline:%d, column:%d - %s", e.Line, e.Column, e.Message)
}
