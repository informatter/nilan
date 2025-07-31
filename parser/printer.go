package parser

import (
	"fmt"
)

// A struct which implements the Visitor interface
// and prints an Abstract Syntax Tree (AST)
type astPrinter struct{}

func (p astPrinter) VisitBinary(b Binary) string {

	return parenthesize(b.Operator.Value, b.Left, b.Right)
}
func (p astPrinter) VisitUnary(u Unary) string {
	return parenthesize(u.Operator.Value, u.Right)
}
func (p astPrinter) VisitLiteral(l Literal) string {
	return fmt.Sprintf("%v", l.Value)
}
func (p astPrinter) VisitGrouping(g Grouping) string {
	return parenthesize("group", g.Expression)
}

// parenthesize creates an S-expression in order to visualise
// the expression presedence order within the AST.
func parenthesize(name string, expressions ...Expression) string {
	astString := "(" + name
	for _, expression := range expressions {
		astString += " " + expression.Accept(astPrinter{})
	}
	astString += ")"
	return astString
}
