package parser

import (
	"fmt"
)

// A struct which implements the Visitor interface
// and prints an Abstract Syntax Tree (AST)
type astPrinter struct{}

func (p astPrinter) VisitExpressionStmt(exprStmt ExpressionStmt) any {
	return p.parenthesize("expression", exprStmt.Expression)
}

func (p astPrinter) VisitPrintStmt(printStmt PrintStmt) any {

	return p.parenthesize("print", printStmt.Expression)
}

func (p astPrinter) VisitVarStmt(varStmt VarStmt) any {
	return p.parenthesize(varStmt.Name.Lexeme, varStmt.Initializer)
}

func (p astPrinter) VisitVariableExpression(variale Variable) any {
	return fmt.Sprintf("%s", variale.Name.Literal)
}

func (p astPrinter) VisitBinary(b Binary) any {

	return p.parenthesize(b.Operator.Lexeme, b.Left, b.Right)
}
func (p astPrinter) VisitUnary(u Unary) any {
	return p.parenthesize(u.Operator.Lexeme, u.Right)
}
func (p astPrinter) VisitLiteral(l Literal) any {
	return fmt.Sprintf("%v", l.Value)
}
func (p astPrinter) VisitGrouping(g Grouping) any {
	return p.parenthesize("group", g.Expression)
}

// parenthesize creates an S-expression in order to visualise
// the expression presedence order within the AST.
func (p astPrinter) parenthesize(name string, expressions ...Expression) string {
	astString := "(" + name
	for _, expression := range expressions {
		astString += " " + expression.Accept(p).(string)
	}
	astString += ")"
	return astString
}
