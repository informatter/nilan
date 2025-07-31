package parser


// An interface to be used to implement the Visitor
// design pattern
type Visitor interface {
	VisitBinary(binary Binary) string
	VisitUnary(unary Unary) string
	VisitLiteral(literal Literal) string
	VisitGrouping(grouping Grouping) string
	// TODO: Add further grammar production rules.
}
