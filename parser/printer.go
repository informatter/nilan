package parser

import (
	"fmt"
	"nilan/ast"
)

// A struct which implements the Visitor interface
// and prints an Abstract Syntax Tree (AST)
type astPrinter struct{}

func (p astPrinter) VisitExpressionStmt(exprStmt ast.ExpressionStmt) any {
	return p.parenthesize("expression", exprStmt.Expression)
}

func (p astPrinter) VisitPrintStmt(printStmt ast.PrintStmt) any {

	return p.parenthesize("print", printStmt.Expression)
}

func (p astPrinter) VisitVarStmt(varStmt ast.VarStmt) any {
	return p.parenthesize(varStmt.Name.Lexeme, varStmt.Initializer)
}

func (p astPrinter) VisitAssignExpression(assign ast.Assign) any {
	return p.parenthesize(assign.Name.Lexeme, assign.Value)
}

func (p astPrinter) VisitVariableExpression(variale ast.Variable) any {
	return fmt.Sprint(variale.Name.Lexeme)
}

func (p astPrinter) VisitBinary(b ast.Binary) any {

	return p.parenthesize(b.Operator.Lexeme, b.Left, b.Right)
}
func (p astPrinter) VisitUnary(u ast.Unary) any {
	return p.parenthesize(u.Operator.Lexeme, u.Right)
}
func (p astPrinter) VisitLiteral(l ast.Literal) any {
	return fmt.Sprintf("%v", l.Value)
}
func (p astPrinter) VisitGrouping(g ast.Grouping) any {
	return p.parenthesize("group", g.Expression)
}

// parenthesize creates an S-expression in order to visualise
// the expression presedence order within the AST.
func (p astPrinter) parenthesize(name string, expressions ...ast.Expression) string {
	astString := "(" + name
	for _, expression := range expressions {
		if expression == nil {
			astString += " " + "null"
		} else {
			astString += " " + expression.Accept(p).(string)
		}
	}
	astString += ")"
	return astString
}
