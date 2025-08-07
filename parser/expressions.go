package parser

import (
	"nilan/token"
)

// Interface for all AST nodes to implement. All concrete implementations
// must implement the Accept method.
type Expression interface {

	// Calls a method on the Visitor interface
	// which performs an action on an expression
	Accept(v Visitor) any
}

type Binary struct {
	Left     Expression
	Right    Expression
	Operator token.Token
}

func (binary Binary) Accept(v Visitor) any {
	return v.VisitBinary(binary)
}

type Unary struct {
	Right    Expression
	Operator token.Token
}

func (unary Unary) Accept(v Visitor) any {

	return v.VisitUnary(unary)
}

type Literal struct {
	Value string
}

func (literal Literal) Accept(v Visitor) any {
	return v.VisitLiteral(literal)
}

type Grouping struct {
	Expression Expression
}

func (grouping Grouping) Accept(v Visitor) any {
	return v.VisitGrouping(grouping)
}
