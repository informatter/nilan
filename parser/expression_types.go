package parser

import (
	"nilan/token"
)

type Expression interface {
	Print() string
}

type Binary struct {
	Left     Expression
	Right    Expression
	Operator token.Token
}

func (binaryExpression Binary) Print() string {
	return ""
}

type Unary struct {
	Right    Expression
	Operator token.Token
}

func (expression Unary) Print() string {
	return ""
}

type Literal struct {
	Value string
}

func (expression Literal) Print() string {
	return ""
}

type Grouping struct {
	Expression Expression
}

func (expression Grouping) Print() string {
	return ""
}
