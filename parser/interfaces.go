package parser

// An interface to be used to implement the Visitor
// design pattern
type Visitor interface {
	VisitBinary(binary Binary) any
	VisitUnary(unary Unary) any
	VisitLiteral(literal Literal) any
	VisitGrouping(grouping Grouping) any
	// TODO: Add further grammar production rules.
}
