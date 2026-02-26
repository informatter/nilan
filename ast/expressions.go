// expressions.go contains all the expression AST nodes. A expression node always evaluates to a value.

package ast

import (
	"nilan/token"
)

// BinaryExpr represents a binary operation expression (e.g., "a + b").
// It consists of a left-hand side expression, an operator token (e.g., +, -, *, /),
// and a right-hand side expression.
type BinaryExpr struct {
	Left     Expr        // The left-hand expression (e.g., "a" in "a + b")
	Operator token.Token // The operator (e.g., "+")
	Right    Expr        // The right-hand expression (e.g., "b" in "a + b")
}

func (binary BinaryExpr) Accept(v ExprVisitor) any {
	return v.VisitBinaryExpr(binary)
}

// UnaryExpr represents a unary operation expression (e.g., "!a" or "-b").
// It consists of an operator token and a single right-hand expression.
type UnaryExpr struct {
	Operator token.Token // The operator (e.g., "!" or "-")
	Right    Expr        // The expression the operator is applied to (e.g., "a" or "b")
}

func (unary UnaryExpr) Accept(v ExprVisitor) any {
	return v.VisitUnaryExpr(unary)
}

// LiteralExpr represents a literal value in the source code
// (e.g., numbers, strings, booleans, or null).
type LiteralExpr struct {
	Value any // The literal value (Go's `any` allows different possible types)
}

func (literal LiteralExpr) Accept(v ExprVisitor) any {
	return v.VisitLiteralExpr(literal)
}

// GroupingExpr represents a parenthesized expression (e.g., "(a + b)").
// Useful for controlling evaluation precedence.
type GroupingExpr struct {
	Expression Expr // The inner expression inside the parentheses
}

func (grouping GroupingExpr) Accept(v ExprVisitor) any {
	return v.VisitGroupingExpr(grouping)
}

// VariableExpr represents a value binded to a declared
// variable
type VariableExpr struct {
	Name token.Token // An IDENTIFIER token
}

// Variable represents a variable expression in the abstract syntax tree (AST).
// It models the retrieval of a value previously bound to a variable name.
//
// Fields:
//   - Name: The token corresponding to the variable's identifier. This is an
//     IDENTIFIER token that holds the variable's name (lexeme).
func (variable VariableExpr) Accept(v ExprVisitor) any {
	return v.VisitVariableExpr(variable)
}

// AssignExpr represents an assignment expression in the abstract syntax tree (AST).
// It models the operation of assigning a new value to an existing variable.
//
// Fields:
//   - Name: The token corresponding to the variable's identifier.
//   - Value: The expression that produces the value being assigned to the variable.
//     This can be any valid expression node in the AST, which will be
//     evaluated and then stored in the environment.
type AssignExpr struct {
	Name  token.Token
	Value Expr
}

func (assign AssignExpr) Accept(v ExprVisitor) any {
	return v.VisitAssignExpr(assign)
}

// LogicalExpr represents a logical expression in the abstract syntax tree (AST).
// It models the logical `if`, `else` expressions.
//
// Fields:
//   - Left: The left side expression.
//   - Operator: The operator being represented, currently it can be either `or` or `and`.
//   - Right: The right side expression, for example.
//
// Example:
// >>> `if a<b or b>10 and b!=0`
type LogicalExpr struct {
	Left     Expr
	Operator token.Token
	Right    Expr
}

func (logical LogicalExpr) Accept(v ExprVisitor) any {
	return v.VisitLogicalExpr(logical)
}

// CallExpr represents a function call expression in the abstract syntax tree (AST).
// It models the invocation of a function with a set of arguments, or a function without arguments.
// Example:
// >>> `add(1, 2)` represents a call to the function `add` with the arguments `1` and `2`.
type CallExpr struct {
	// Callee is the expression representing the function being called.
	// For example, in `add(1, 2)`, the callee is the variable expression `add`.
	Callee Expr
	// Parent is the token representing the closing parenthesis of the function call.
	Paren token.Token
	// Arguments is a slice of expressions representing the arguments passed to the function.
	// For example, in `add(1, 2)`, the arguments are the literal expressions `1` and `2`.
	Arguments []Expr
}

func (call CallExpr) Accept(v ExprVisitor) any {
	return v.VisitCallExpr(call)
}
