// statements.go contains all the statement AST nodes. A statement node does not produce a value.

package ast

import "nilan/token"

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

// BlockStmt represents a block statement containing a list
// of statement expression AST nodes.
type BlockStmt struct {
	Statements []Stmt
}

func (blockStmt BlockStmt) Accept(v StmtVisitor) any {
	return v.VisitBlockStmt(blockStmt)
}

// IfStmt represents an if statement containing the expression
// to evaluate the statement to execute if the expression is true
// or the statement to execute of the expression is false.
type IfStmt struct {
	Condition Expression
	Then      Stmt
	Else      Stmt
}

func (stmt IfStmt) Accept(v StmtVisitor) any {
	return v.VisitIfStmt(stmt)
}

// WhileStmt represents a while loop AST node.
//
// Fields:
//   - Condition: The expression evaluated before each iteration of the loop.
//     If this expression evaluates to true, the loop body executes;
//     otherwise, the loop terminates.
//   - Body: The block statement representing the loop body,
type WhileStmt struct {
	Condition Expression
	Body      Stmt
}

func (stmt WhileStmt) Accept(v StmtVisitor) any {
	return v.VisitWhileStmt(stmt)
}
