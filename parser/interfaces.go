package parser

// Visitor is the interface for operating on all Expression AST nodes.
// Any type that wants to perform an operation on expressions (e.g., an interpreter,
// ast-printer, or type checker) must implement this interface.
//
// Each Visit method corresponds to a distinct Expression type.
type Visitor interface {
	// VisitBinary is called when visiting a Binary expression (e.g., "a + b").
	VisitBinary(binary Binary) any

	// VisitUnary is called when visiting a Unary expression (e.g., "!a" or "-b").
	VisitUnary(unary Unary) any

	// VisitLiteral is called when visiting a Literal expression (e.g., a number, string, or boolean).
	VisitLiteral(literal Literal) any

	// VisitGrouping is called when visiting a Grouping expression (expressions wrapped in parentheses).
	VisitGrouping(grouping Grouping) any

	// TODO: Add further Visit methods as new expression grammar rules are introduced.
}

// StmtVisitor is the interface for operating on all Statement AST nodes.
// Like Visitor, it defines one Visit method per statement type.
// This separation between expressions and statements mirrors the grammar structure.
type StmtVisitor interface {
	// VisitExpressionStmt is called when visiting an Expression statement.
	// Example: "foo + bar;"
	VisitExpressionStmt(exprStmt ExpressionStmt) any

	// VisitPrintStmt is called when visiting a Print statement.
	// Example: "print foo + bar;"
	VisitPrintStmt(printStmt PrintStmt) any
	// TODO: Add further visit methods as new statement grammar rules are introduced.
}
