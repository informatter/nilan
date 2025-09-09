package interpreter

import (
	"fmt"
	"nilan/token"
)

// Environment represents a mapping of variable names to their values.
// It serves as a storage for variable bindings for the interpreter,
// where each key is a variable identifier and the value
// is the data associated with that variable.
type Environment struct {
	values map[string]any
}

// MakeEnvironment creates and returns a new, empty Environment instance.
//
// Returns:
//   - *Environment: A pointer to a newly allocated Environment structure
//     with an initialized values map.
func MakeEnvironment() *Environment {
	return &Environment{
		values: make(map[string]any),
	}
}

// assign attempts to update the value of an existing variable in the current environment.
//
// Parameters:
//   - name: A token representing the variable identifier. Its Lexeme field is used as the
//     variable's name.
//   - value: The new value to assign to the variable. It can be of any type (interface{}).
//
// Returns:
//   - error: Returns nil if the assignment succeeds. If the variable does not already
//     exist in the environment, returns a runtime error with source location
//     information (line and column).
func (env *Environment) assign(name token.Token, value any) error {
	_, ok := env.values[name.Lexeme]
	if ok {
		env.set(name.Lexeme, value)
		return nil
	}

	msg := fmt.Sprintf("Undefined variable: %s", name.Lexeme)
	return CreateRuntimeError(name.Line, name.Column, msg)
}

// Sets a variable in the environment
// Parameters:
//   - name: string
//     The name of the variable, i.e its indentifier
//   - value: any
//     The value assigned to the variable.
func (env *Environment) set(name string, value any) {
	env.values[name] = value
}

// Gets the value associated to a variable from the environment
// Parameters:
//   - name: token.Token
//     The variable to retrieve its value
//
// Returns:
//   - any: The value of the specified variable
//   - error: A RuntimeError if the variable has not been previously
//     assigned and its trying to be accessed.
func (env *Environment) get(name token.Token) (any, error) {
	value, ok := env.values[name.Lexeme]
	if ok {
		return value, nil
	}
	msg := fmt.Sprintf("Undefined variable: %s", name.Lexeme)
	return nil, CreateRuntimeError(name.Line, name.Column, msg)
}
