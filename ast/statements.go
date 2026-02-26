// statements.go contains all the statement AST nodes. A statement node does not produce a value.

package ast

import "nilan/token"

// ExpressionStmt represents a statement that consists of a single expression.
// Example: `foo + bar;`
// This evaluates the expression and discards the result.
type ExpressionStmt struct {
	Expression Expr // The expression used as a statement
}

func (e ExpressionStmt) Accept(v StmtVisitor) any {
	return v.VisitExpressionStmt(e)
}

// PrintStmt represents a print statement that outputs the result
// of evaluating an expression. Example: `print foo + bar;`
type PrintStmt struct {
	Expression Expr // The expression whose result will be printed
}

func (p PrintStmt) Accept(v StmtVisitor) any {
	return v.VisitPrintStmt(p)
}

// VarStmt represents a variable declaration statement, its composed
// of the name of the variable and the expression it binds to. A declaration
// statement declares functions, variables and classes.
type VarStmt struct {

	// Name is the token representing the variable's identifier.
	Name token.Token

	// Initialiser is the expression assigned to the variable when declared.
	// For example, `var x=5` the initialser is `5`. Since this is an expression,
	// this is also supported `var x = 5+3`.
	Initializer Expr
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
	Condition Expr
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
	Condition Expr
	Body      Stmt
}

func (stmt WhileStmt) Accept(v StmtVisitor) any {
	return v.VisitWhileStmt(stmt)
}

// FunctionStmt represents a function declaration statement, which consists of the function's name,
// its parameters, and its body.
type FunctionStmt struct {

	// Name is the token representing the function's identifier.
	Name token.Token
	// Body is the block statement representing the function's body, which contains the statements that define the
	// function's behavior.
	Body BlockStmt
	// Parameters is a list of tokens representing the function's parameters.
	Parameters []token.Token
}

func (stmt FunctionStmt) Accept(v StmtVisitor) any {
	return v.VisitFunctionStmt(stmt)
}

// ReturnStmt represents a return statement, which consists of the return keyword and the expression whose value
// is to be returned.
type ReturnStmt struct {
	// Keyword is the token representing the 'return' keyword.
	Keyword token.Token
	// Value is the expression whose value is to be returned when the return statement is executed.
	// TODO: Add handling for returning functions.
	Value Expr
}

func (stmt ReturnStmt) Accept(v StmtVisitor) any {
	return v.VisitReturnStmt(stmt)
}
