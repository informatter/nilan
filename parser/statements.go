package parser

import "nilan/token"

// Stmt is the base interface for all statement nodes in the AST.
// Like Expression, it follows the Visitor design pattern where each
// statement type implements Accept, calling back into the correct
// Visit method on a StmtVisitor.
//
// A statement represents an action in a program (e.g., printing,
// evaluating an expression, variable declaration). Unlike expressions,
// statements typically do not produce a value.
type Stmt interface {
	// Accept dispatches this statement to the appropriate Visit method
	// of the provided StmtVisitor implementation.
	Accept(v StmtVisitor) any
}

// ExpressionStmt represents a statement that consists of a single expression.
// Example: `foo + bar;`
// This evaluates the expression and discards the result.
type ExpressionStmt struct {
	Expression Expression // The expression used as a statement
}

func (e ExpressionStmt) Accept(v StmtVisitor) any {
	return v.VisitExpressionStmt(e)
}

// PrintStmt represents a print statement that outputs the result
// of evaluating an expression. Example: `print foo + bar;`
type PrintStmt struct {
	Expression Expression // The expression whose result will be printed
}

func (p PrintStmt) Accept(v StmtVisitor) any {
	return v.VisitPrintStmt(p)
}

// VarStmt represents a variable declaration statement, its composed
// of the name of the variable and the expression it binds to. A declaration
// statement declares functions, variables and classes.
type VarStmt struct {
	Name        token.Token
	Initializer Expression
}

func (varStmt VarStmt) Accept(v StmtVisitor) any {
	return v.VisitVarStmt(varStmt)
}
