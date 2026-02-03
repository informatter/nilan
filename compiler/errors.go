package compiler

import "fmt"

type SemanticError struct {
	Message string
}

func (e SemanticError) Error() string {
	return fmt.Sprintf("💥 SemanticError: %s", e.Message)
}

type DeveloperError struct {
	Message string
}

func (e DeveloperError) Error() string {
	return fmt.Sprintf("🤖 DeveloperError: %s", e.Message)
}
