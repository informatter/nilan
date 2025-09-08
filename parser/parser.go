// Recursive descent parser
// https://en.wikipedia.org/wiki/Recursive_descent_parser

//	A Recursive descent parser is a top-down parser because it starts from the top
//
// grammar rule and works its way down in to the nested sub-experessions before reaching
// the leaves of the syntax tree (terminal rules)
package parser

import (
	"fmt"
	"nilan/token"
)

var comparisonTokenTypes = []token.TokenType{
	token.LARGER,
	token.LARGER_EQUAL,
	token.LESS,
	token.LESS_EQUAL,
}

var equalityTokenTypes = []token.TokenType{
	token.NOT_EQUAL,
	token.EQUAL_EQUAL,
}

var termTokenTypes = []token.TokenType{
	token.SUB,
	token.ADD,
}

var factorExpressionTypes = []token.TokenType{
	token.MULT,
	token.DIV,
}

var unaryExpressionTypes = []token.TokenType{
	token.BANG,
	token.SUB,

	// NOTE: not supported operands on unary expressions are included
	// So they can be parsed, but then the interpreter can throw a more detailed
	// runtime error message. This is known as "error productions"
	token.MULT,
	token.ADD,
	token.DIV,
}

type Parser struct {
	tokens   []token.Token
	position int
}

// NOTE: The parsers position is always one unit ahead of the
// current token

// Initializes and returns a new Parser instance.
//
// Parameters:
//   - tokens: []token.Token
//     The tokens created by the lexer.
//   - position: int
//     The position of the parser in respect to the current token being
//     looked at.
//
// Returns:
//   - *Parser: A pointer to a newly created Parser instance.
func Create(tokens []token.Token) *Parser {
	return &Parser{
		tokens:   tokens,
		position: 0,
	}
}

// Print prints the string representations of a slice of Stmt nodes
// using the AST printer. Each statement is visited with an astPrinter,
// and the resulting string is output using fmt.Println.
//
// This method does not return any value; its purpose is to output
// formatted representations of the AST nodes to standard output.
func (parser *Parser) Print(statements []Stmt) {
	for _, s := range statements {
		result := s.Accept(astPrinter{})
		fmt.Println(result)
	}
}

// Peeks the token at the parser's current position,
// without advancing the parser's position.
// Returns:
//   - token.Token: The token at the parser's current position
func (parser *Parser) peek() token.Token {
	return parser.tokens[parser.position]
}

// Retrieves the token at the parser's previous position
// (position -1)
//
// Returns:
//   - token.Token: The token at the previous position
func (parser *Parser) previous() token.Token {
	return parser.tokens[parser.position-1]
}

// Increments the parser's position by one unit and
// consumes the current token
//
// Returns:
//   - token.Token: The token at the previous position
func (parser *Parser) advance() token.Token {
	if !parser.isFinished() {
		parser.position++
	}
	return parser.previous()
}

// Determines of the parser has finished scanning all the tokens.
//
// Returns:
//   - bool: true if the parser has finished scanning, false otherwise
func (parser *Parser) isFinished() bool {
	tok := parser.peek()
	return tok.TokenType == token.EOF
}

// Determines if the provided tokenType matches the TokenType
// at the parser's current position
//
// Returns
//   - bool: true if the TokenType matches, false otherwise
func (parser *Parser) checkType(tokeType token.TokenType) bool {
	if parser.isFinished() {
		return false
	}
	tok := parser.peek()
	return tok.TokenType == tokeType
}

// Determines if the TokenType at the current
// position matches any of the provided tokenTypes. If a match is
// found the parser increments its position and consumes the
// current token
//
// Returns
//   - bool: true if a match was found, false otherwise
func (parser *Parser) isMatch(tokenTypes []token.TokenType) bool {
	for i := range tokenTypes {
		tokenType := tokenTypes[i]

		if parser.checkType(tokenType) {
			parser.advance()
			return true
		}
	}
	return false
}

// Parse parses the entire token stream into a slice of Stmt (statement) nodes,
// continuing until the end of input. Errors during parsing are collected
// but parsing continues to find additional errors where possible.
//
// Returns:
//   - []Stmt: the successfully parsed statements.
//   - []error: all errors that occurred during parsing.
func (parser *Parser) Parse() ([]Stmt, []error) {
	statements := []Stmt{}
	errors := []error{}

	for {
		if parser.isFinished() {
			break
		}
		statement, err := parser.declaration()
		if err != nil {
			errors = append(errors, err)
			continue
		}
		statements = append(statements, statement)
	}

	for _, err := range errors {
		fmt.Println(err.Error())
	}
	return statements, errors
}

// declaration parses a declaration statement.
// 
// It first checks if the next token is a variable declaration keyword (e.g., `var`).
// If so, it calls the variableDeclaration method to parse the variable declaration statement.
// 
// TODO: Support for function and class declarations will be added later.
//
// If the next token is not a variable declaration, it parses a general statement.
//
// Returns the parsed statement (Stmt) or an error if parsing fails.
//
func (parser *Parser) declaration() (Stmt, error) {
	if parser.isMatch([]token.TokenType{token.VAR}) {
		return parser.variableDeclaration()
	}
	// TODO Add support for functions and classes
	return parser.statement()
}

// variableDeclaration parses and creates a variable declaration statement.
//
// It expects the next token to be an identifier representing the variable's name.
// If the token is not an identifier, it returns an error indicating the expected variable name.
//
// After successfully consuming the identifier, it optionally parses an initializer expression
// if an assignment operator (=) is found. If an initializer is present, it is parsed as an expression.
//
// It returns a VarStmt representing the variable declaration with the variable name and optional initialiser
// expression, or an error if parsing fails at any point.
// 
// Returns:
//  - VarStmt: The variable declaration statement AST node.
//
// Example input:
//   
//   >>> var x = 10
// 
func (parser *Parser) variableDeclaration() (Stmt, error) {
	tok, consumeError := parser.consume(token.IDENTIFIER, "Expected variable name")
	if consumeError != nil {
		return nil, consumeError
	}

	var initialiser Expression
	if parser.isMatch([]token.TokenType{token.ASSIGN}) {
		var err error
		initialiser, err = parser.expression()
		if err != nil {
			return nil, err
		}
	}

	return VarStmt{
		Name:        tok,
		Initializer: initialiser,
	}, nil
}

// statement parses a single statement. Currently, this can be either
// a print statement ("print <expr>") or an expression statement.
//
// Returns:
//   - Stmt: the parsed statement node.
//   - error: if parsing fails, otherwise nil.
func (parser *Parser) statement() (Stmt, error) {

	if parser.isMatch([]token.TokenType{token.PRINT}) {
		printStatement, err := parser.printStatement()
		if err != nil {
			return nil, err
		}
		return printStatement, nil
	}
	// TODO: Add more expression types.
	exprStatement, err := parser.expressionStatement()
	if err != nil {
		return nil, err
	}
	return exprStatement, nil
}

// printStatement parses a print statement of the form "print <expression>".
//
// Returns:
//   - Stmt: a PrintStmt containing the expression to print.
//   - error: if the inner expression fails to parse.
func (parser *Parser) printStatement() (Stmt, error) {
	expression, err := parser.expression()
	if err != nil {
		return nil, err
	}
	return PrintStmt{Expression: expression}, nil
}

// expressionStatement parses a statement consisting of a single expression.
//
// Returns:
//   - Stmt: an ExpressionStmt wrapping the parsed expression.
//   - error: if the expression cannot be parsed.
func (parser *Parser) expressionStatement() (Stmt, error) {
	expression, err := parser.expression()
	if err != nil {
		return nil, err
	}
	return ExpressionStmt{Expression: expression}, nil
}

// expression is the entry point for parsing expressions. It begins at
// the equality rule, which encompasses all lower-precedence rules.
//
// Returns:
//   - Expression: the parsed expression AST node.
//   - error: if parsing fails.
func (parser *Parser) expression() (Expression, error) {
	return parser.equality()
}

// equality parses equality expressions using operators "==" and "!=".
//
// Returns:
//   - Expression: a Binary node (or sub-expression) representing equality comparison.
//   - error: if parsing fails.
func (parser *Parser) equality() (Expression, error) {
	exp, err := parser.comparison()
	if err != nil {
		return nil, err
	}
	for parser.isMatch(equalityTokenTypes) {
		operator := parser.previous()
		right, err := parser.comparison()
		if err != nil {
			return nil, err
		}
		exp = Binary{
			Left:     exp,
			Operator: operator,
			Right:    right,
		}
	}
	return exp, nil
}

// comparison parses comparison expressions using operators "<", "<=", ">", ">=".
//
// Returns:
//   - Expression: a Binary node (or sub-expression) representing a comparison.
//   - error: if parsing fails.
func (parser *Parser) comparison() (Expression, error) {
	exp, err := parser.term()
	if err != nil {
		return nil, err
	}
	for parser.isMatch(comparisonTokenTypes) {
		operator := parser.previous()
		right, err := parser.term()
		if err != nil {
			return nil, err
		}
		exp = Binary{
			Left:     exp,
			Operator: operator,
			Right:    right,
		}
	}
	return exp, nil
}

// term parses addition and subtraction expressions using operators "+" and "-".
//
// Returns:
//   - Expression: a Binary node (or sub-expression) representing addition or subtraction.
//   - error: if parsing fails.
func (parser *Parser) term() (Expression, error) {
	exp, err := parser.factor()
	if err != nil {
		return nil, err
	}
	for parser.isMatch(termTokenTypes) {
		operator := parser.previous()
		right, err := parser.factor()
		if err != nil {
			return nil, err
		}
		exp = Binary{
			Left:     exp,
			Operator: operator,
			Right:    right,
		}
	}
	return exp, nil
}

// factor parses multiplication and division expressions using operators "*" and "/".
//
// Returns:
//   - Expression: a Binary node (or sub-expression) representing multiplication or division.
//   - error: if parsing fails.
func (parser *Parser) factor() (Expression, error) {
	exp, err := parser.unary()
	if err != nil {
		return nil, err
	}
	for parser.isMatch(factorExpressionTypes) {
		operator := parser.previous()
		right, err := parser.unary()
		if err != nil {
			return nil, err
		}
		exp = Binary{
			Left:     exp,
			Operator: operator,
			Right:    right,
		}
	}
	return exp, nil
}

// unary parses unary prefix expressions using operators "!" or "-".
// Examples: "!true", "-x".
//
// Returns:
//   - Expression: a Unary node if a unary operator was found, otherwise defers to primary().
//   - error: if parsing fails.
func (parser *Parser) unary() (Expression, error) {
	if parser.isMatch(unaryExpressionTypes) {
		operator := parser.previous()
		right, err := parser.unary()
		if err != nil {
			return nil, err
		}
		return Unary{
			Operator: operator,
			Right:    right,
		}, nil
	}
	return parser.primary()
}

// primary parses the most basic forms of expressions:
//   - Literals: true, false, null, strings, numbers
//   - Grouping: (expression)
//
// If no valid token matches, returns a syntax error.
//
// Returns:
//   - Expression: a Literal, Grouping expression .
//   - error: if no valid primary expression can be parsed.
func (parser *Parser) primary() (Expression, error) {
	if parser.isMatch([]token.TokenType{token.FALSE}) {
		return Literal{Value: false}, nil
	}
	if parser.isMatch([]token.TokenType{token.NULL}) {
		return Literal{Value: nil}, nil
	}
	if parser.isMatch([]token.TokenType{token.TRUE}) {
		return Literal{Value: true}, nil
	}

	if parser.isMatch([]token.TokenType{token.FLOAT, token.INT, token.STRING}) {
		return Literal{Value: parser.previous().Literal}, nil
	}

	if parser.isMatch([]token.TokenType{token.IDENTIFIER}) {
		return Variable{Name: parser.previous()}, nil
	}

	if parser.isMatch([]token.TokenType{token.LPA}) {
		expr, err := parser.expression()
		if err != nil {
			return nil, err
		}
		_, consumeErr := parser.consume(token.RPA, fmt.Sprintf("expression is missing '%s'", token.RPA))
		if consumeErr != nil {
			return nil, consumeErr
		}
		return Grouping{Expression: expr}, nil
	}

	currentToken := parser.peek()
	return nil, CreateSyntaxError(currentToken.Line, currentToken.Column, "Unrecognised expression.")
}

// Consumes the current token by advancing the parsers current position by
// one unit if the `tokenType` matches the token type of the parsers current
// position.
//
//	Returns:
//	- A SyntaxError if the provided `tokenType` does not match the `TokenType`
//		at the parsers current position
func (parser *Parser) consume(tokenType token.TokenType, errorMessage string) (token.Token, error) {
	if parser.checkType(tokenType) {
		return parser.advance(), nil
	}
	currentToken := parser.peek()
	return token.CreateToken(token.EOF, 0, 0), CreateSyntaxError(currentToken.Line, currentToken.Column, errorMessage)
}
