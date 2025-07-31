package parser

import (
	"nilan/token"
)

// Interface for all AST nodes to implement. All concrete implementations
// must implement the Accept method.
type Expression interface {

	// Calls a method on the Visitor interface
	// which performs an action on an expression
	Accept(v Visitor) string
}

type Binary struct {
	Left     Expression
	Right    Expression
	Operator token.Token
}

func (binary Binary) Accept(v Visitor) string {
	return v.VisitBinary(binary)
}

type Unary struct {
	Right    Expression
	Operator token.Token
}

func (unary Unary) Accept(v Visitor) string {
	return v.VisitUnary(unary)
}

type Literal struct {
	Value string
}

func (literal Literal) Accept(v Visitor) string {
	return v.VisitLiteral(literal)
}

type Grouping struct {
	Expression Expression
}

func (grouping Grouping) Accept(v Visitor) string {
	return v.VisitGrouping(grouping)
}
