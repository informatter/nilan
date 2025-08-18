package parser

import (
	"nilan/token"
)

// Expression is the core interface for all expression nodes in the Abstract Syntax Tree (AST).
// Any expression type (e.g., binary operation, literal, grouping, etc.) must implement this interface.
// The Accept method enables the Visitor design pattern so that operations can be performed on
// expressions without the expression types needing to know the details of those operations.
// The visitor pattern decoupled behaviour from data to easily allow adding the behaviour to objects
// without the need to change the objects themselves.
type Expression interface {
	// Accept dispatches the current expression node to the appropriate method on a Visitor.
	// v: the Visitor instance that defines behavior for this expression type
	// Returns: a generic result (any), since the Visitor may define its own return type
	Accept(v Visitor) any
}

// Binary represents a binary operation expression (e.g., "a + b").
// It consists of a left-hand side expression, an operator token (e.g., +, -, *, /),
// and a right-hand side expression.
type Binary struct {
	Left     Expression  // The left-hand expression (e.g., "a" in "a + b")
	Operator token.Token // The operator (e.g., "+")
	Right    Expression  // The right-hand expression (e.g., "b" in "a + b")
}

func (binary Binary) Accept(v Visitor) any {
	return v.VisitBinary(binary)
}

// Unary represents a unary operation expression (e.g., "!a" or "-b").
// It consists of an operator token and a single right-hand expression.
type Unary struct {
	Operator token.Token // The operator (e.g., "!" or "-")
	Right    Expression  // The expression the operator is applied to (e.g., "a" or "b")
}

func (unary Unary) Accept(v Visitor) any {
	return v.VisitUnary(unary)
}

// Literal represents a literal value in the source code
// (e.g., numbers, strings, booleans, or null).
type Literal struct {
	Value any // The literal value (Go's `any` allows different possible types)
}

func (literal Literal) Accept(v Visitor) any {
	return v.VisitLiteral(literal)
}

// Grouping represents a parenthesized expression (e.g., "(a + b)").
// Useful for controlling evaluation precedence.
type Grouping struct {
	Expression Expression // The inner expression inside the parentheses
}

func (grouping Grouping) Accept(v Visitor) any {
	return v.VisitGrouping(grouping)
}
