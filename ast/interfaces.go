// interfaces.go contains all visitor interfaces that any code traversing expression and statement AST nodes must implement.
// It also contains the interfaces that all statement and expression AST nodes must implement which also follows the
// visitor design pattern

package ast

// ExpressionVisitor is the interface for operating on all Expression AST nodes.
// Any type that wants to perform an operation on expressions (e.g., an interpreter,
// ast-printer, or type checker) must implement this interface.
//
// Each Visit method corresponds to a distinct Expression type.
type ExpressionVisitor interface {
	// VisitBinary is called when visiting a Binary expression (e.g., "a + b").
	VisitBinary(binary Binary) any

	// VisitUnary is called when visiting a Unary expression (e.g., "!a" or "-b").
	VisitUnary(unary Unary) any

	// VisitLiteral is called when visiting a Literal expression (e.g., a number, string, or boolean).
	VisitLiteral(literal Literal) any

	// VisitGrouping is called when visiting a Grouping expression (expressions wrapped in parentheses).
	VisitGrouping(grouping Grouping) any

	VisitVariableExpression(variable Variable) any

	VisitAssignExpression(assign Assign) any

	VisitLogicalExpression(logical Logical) any

	// TODO: Add further Visit methods as new expression grammar rules are introduced.
}

// StmtVisitor is the interface for operating on all Statement AST nodes.
// Like ExpressionVisitor, it defines one Visit method per statement type.
// This separation between expressions and statements mirrors the grammar structure.
type StmtVisitor interface {
	// VisitExpressionStmt is called when visiting an Expression statement.
	// Example: "foo + bar;"
	VisitExpressionStmt(exprStmt ExpressionStmt) any

	// VisitPrintStmt is called when visiting a Print statement.
	// Example: "print foo + bar;"
	VisitPrintStmt(printStmt PrintStmt) any

	// visitVarStmt is called when visiting a declaration statement.
	// Example: "name = 'foo'"
	VisitVarStmt(varStmt VarStmt) any

	// VisitBlockStmt is called when visiting a block statement.
	VisitBlockStmt(blockStmt BlockStmt) any

	VisitIfStmt(stmt IfStmt) any

	VisitWhileStmt (stmt WhileStmt) any

	// TODO: Add further visit methods as new statement grammar rules are introduced.
}

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
	Accept(v ExpressionVisitor) any
}
