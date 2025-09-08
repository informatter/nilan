package interpreter

import (
	"fmt"
	"nilan/token"
)

// Defines the bindings that associate variables to values.
type Environment struct {
	values map[string]any
}

func MakeEnvironment() *Environment {
	return &Environment{
		values: make(map[string]any),
	}
}

// Sets a variable in the environment
// Parameters:
//	- name: string
//	  The name of the variable, i.e its indentifier
//  - value: any
//     The value assigned to the variable.
func (env *Environment) set(name string, value any) {
	env.values[name] = value
}

// Gets the value associated to a variable from the environment
// Parameters:
//   - name: token.Token
//     The variable to retrieve its value
// Returns:
//  - any: The value of the specified variable
//  - error: A RuntimeError if the variable has not been previously
//    assigned and its trying to be accessed.
func (env *Environment) get(name token.Token) (any, error) {
	value, ok := env.values[name.Lexeme]
	if ok {
		return value, nil
	}
	msg := fmt.Sprintf("Undefined variable: %s", name.Lexeme)
	return nil, CreateRuntimeError(name.Line, name.Column, msg)
}
