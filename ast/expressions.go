// expressions.go contains all the expression AST nodes. A expression node always evaluates to a value.

package ast

import (
	"nilan/token"
)


// Binary represents a binary operation expression (e.g., "a + b").
// It consists of a left-hand side expression, an operator token (e.g., +, -, *, /),
// and a right-hand side expression.
type Binary struct {
	Left     Expression  // The left-hand expression (e.g., "a" in "a + b")
	Operator token.Token // The operator (e.g., "+")
	Right    Expression  // The right-hand expression (e.g., "b" in "a + b")
}

func (binary Binary) Accept(v ExpressionVisitor) any {
	return v.VisitBinary(binary)
}

// Unary represents a unary operation expression (e.g., "!a" or "-b").
// It consists of an operator token and a single right-hand expression.
type Unary struct {
	Operator token.Token // The operator (e.g., "!" or "-")
	Right    Expression  // The expression the operator is applied to (e.g., "a" or "b")
}

func (unary Unary) Accept(v ExpressionVisitor) any {
	return v.VisitUnary(unary)
}

// Literal represents a literal value in the source code
// (e.g., numbers, strings, booleans, or null).
type Literal struct {
	Value any // The literal value (Go's `any` allows different possible types)
}

func (literal Literal) Accept(v ExpressionVisitor) any {
	return v.VisitLiteral(literal)
}

// Grouping represents a parenthesized expression (e.g., "(a + b)").
// Useful for controlling evaluation precedence.
type Grouping struct {
	Expression Expression // The inner expression inside the parentheses
}

func (grouping Grouping) Accept(v ExpressionVisitor) any {
	return v.VisitGrouping(grouping)
}

// Variable represents a value binded to a declared
// variable
type Variable struct {
	Name token.Token // An IDENTIFIER token
}

// Variable represents a variable expression in the abstract syntax tree (AST).
// It models the retrieval of a value previously bound to a variable name.
//
// Fields:
//   - Name: The token corresponding to the variable's identifier. This is an
//     IDENTIFIER token that holds the variable's name (lexeme).
func (variable Variable) Accept(v ExpressionVisitor) any {
	return v.VisitVariableExpression(variable)
}

// Assign represents an assignment expression in the abstract syntax tree (AST).
// It models the operation of assigning a new value to an existing variable.
//
// Fields:
//   - Name: The token corresponding to the variable's identifier.
//   - Value: The expression that produces the value being assigned to the variable.
//     This can be any valid expression node in the AST, which will be
//     evaluated and then stored in the environment.
type Assign struct {
	Name  token.Token
	Value Expression
}

func (assign Assign) Accept(v ExpressionVisitor) any {
	return v.VisitAssignExpression(assign)
}
