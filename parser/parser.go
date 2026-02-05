// Recursive descent parser
// https://en.wikipedia.org/wiki/Recursive_descent_parser

//	A Recursive descent parser is a top-down parser because it starts from the top
//
// grammar rule and works its way down in to the nested sub-experessions before reaching
// the leaves of the syntax tree (terminal rules)
package parser

import (
	"fmt"
	"nilan/ast"
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
func Make(tokens []token.Token) *Parser {
	return &Parser{
		tokens:   tokens,
		position: 0,
	}
}

// Print prints the AST as prettified JSON to standard output.
func (parser *Parser) Print(statements []ast.Stmt) {
	_, err := PrintASTJSON(statements)
	if err != nil {
		fmt.Println("error producing AST JSON:", err)
	}
}

// PrintToFile writes the AST for the provided statements to a .json file at the given path.
func (parser *Parser) PrintToFile(statements []ast.Stmt, path string) error {
	return WriteASTJSONToFile(statements, path)
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
func (parser *Parser) Parse() ([]ast.Stmt, []error) {
	statements := []ast.Stmt{}
	errors := []error{}

	for {
		if parser.isFinished() {
			break
		}
		statement, err := parser.declaration()
		if err != nil {
			errors = append(errors, err)
			if !parser.isFinished() {
				parser.position++
			}
			continue
		}
		statements = append(statements, statement)
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
func (parser *Parser) declaration() (ast.Stmt, error) {
	if parser.isMatch([]token.TokenType{token.VAR}) {
		return parser.variableDeclaration()
	}
	// TODO Add support for functions and classes
	return parser.statement()
}

// variableDeclaration parses a variable declaration statement.
// It expects an identifier token for the variable name
// followed by an optional '=' and an initializer expression.
// Returns:
//   - ast.VarStmt: A VarStmt AST node epresenting the variable declaration.
//   - error: A SyntaxError if parsing fails or if the variable has not been initialised.
func (parser *Parser) variableDeclaration() (ast.Stmt, error) {
	tok, consumeError := parser.consume(token.IDENTIFIER, "Expected variable name")
	if consumeError != nil {
		return nil, consumeError
	}

	var initialiser ast.Expression
	if parser.isMatch([]token.TokenType{token.ASSIGN}) {
		var err error
		initialiser, err = parser.expression()
		if err != nil {
			return nil, err
		}
	}

	return ast.VarStmt{
		Name:        tok,
		Initializer: initialiser,
	}, nil
}

// statement parses a single statement. Currently, this can be either
// a print statement ("print <expr>") an expression statement,
// a block statement or a conditional statement.
//
// Returns:
//   - Stmt: the parsed statement node.
//   - error: if parsing fails, otherwise nil.
func (parser *Parser) statement() (ast.Stmt, error) {

	if parser.isMatch([]token.TokenType{token.PRINT}) {
		printStatement, err := parser.printStatement()
		if err != nil {
			return nil, err
		}
		return printStatement, nil
	}

	if parser.isMatch([]token.TokenType{token.LCUR}) {
		statements, err := parser.block()
		if err != nil {
			return nil, err
		}
		return ast.BlockStmt{Statements: statements}, nil
	}

	if parser.isMatch([]token.TokenType{token.IF}) {
		return parser.ifStatement()
	}

	if parser.isMatch([]token.TokenType{token.WHILE}) {
		return parser.WhileStatement()
	}
	// TODO: Add more expression types.

	expression, err := parser.expression()
	if err != nil {
		return nil, err
	}
	exprStmt := ast.ExpressionStmt{Expression: expression}

	return exprStmt, nil
}

// printStatement parses a print statement of the form "print <expression>".
//
// Returns:
//   - Stmt: a PrintStmt containing the expression to print.
//   - error: if the inner expression fails to parse.
func (parser *Parser) printStatement() (ast.Stmt, error) {
	expression, err := parser.expression()
	if err != nil {
		return nil, err
	}
	return ast.PrintStmt{Expression: expression}, nil
}

// WhileStatement parses a while loop statement from the token stream.
// It parses a condition expression followed by a statement representing
// the loop body.
// Returns:
//   - ast.WhileStmt with the parsed condition and body.
//   - error: if parsing the condition or body fails.
func (parser *Parser) WhileStatement() (ast.Stmt, error) {

	expr, err := parser.expression()
	if err != nil {
		return nil, err
	}

	// NOTE: the statement contains the ast node encompassing all
	// the loops body.
	stmt, error := parser.statement()
	if error != nil {
		return nil, error
	}

	return ast.WhileStmt{
		Condition: expr,
		Body:      stmt,
	}, nil

}

// ifStatement parses an if-statement from the token stream.
// It expects a condition expression followed by a 'then' branch,
// and optionally parses an 'else' branch if present.
// Returns:
//   - ast.IfStmt: an IfStmt AST node.
//   - error: if any part fails to parse.
func (parser *Parser) ifStatement() (ast.Stmt, error) {

	conditionExpr, err := parser.expression()
	if err != nil {
		return nil, err
	}

	thenStmt, err := parser.statement()
	if err != nil {
		return nil, err
	}
	var elseStmt ast.Stmt = nil
	if parser.isMatch([]token.TokenType{token.ELSE}) {
		stmt, err := parser.statement()
		if err != nil {
			return nil, err
		}
		elseStmt = stmt
	}

	return ast.IfStmt{
		Condition: conditionExpr,
		Then:      thenStmt,
		Else:      elseStmt,
	}, nil
}

// expressionStatement parses a statement consisting of a single expression.
//
// Returns:
//   - Stmt: an ExpressionStmt wrapping the parsed expression.
//   - error: if the expression cannot be parsed.
func (parser *Parser) expressionStatement() (ast.Stmt, error) {
	expression, err := parser.expression()
	if err != nil {
		return nil, err
	}
	return ast.ExpressionStmt{Expression: expression}, nil
}

// block parser a block statement consisting of a list of
// statement AST nodes.
// Returns:
//   - [] Stmt: A list of parsed declarations or statements
//   - error: If the block statement cant be parsed.
func (parser *Parser) block() ([]ast.Stmt, error) {
	statements := []ast.Stmt{}

	for !parser.isMatch([]token.TokenType{token.RCUR}) && !parser.isFinished() {
		stmt, err := parser.declaration()
		if err != nil {
			return nil, err
		}
		statements = append(statements, stmt)

	}

	previousToken := parser.previous()
	if previousToken.TokenType != token.RCUR {
		errMsg := fmt.Sprintf("Expected '%s' after block.", token.RCUR)
		err := CreateSyntaxError(previousToken.Line, previousToken.Column, errMsg)
		return nil, err
	}
	return statements, nil
}

// expression is the entry point for parsing expressions. It begins at
// the assignment rule, which encompasses all lower-precedence rules.
//
// Returns:
//   - Expression: the parsed expression AST node.
//   - error: if parsing fails.
func (parser *Parser) expression() (ast.Expression, error) {
	return parser.assignment()
}

// assignment parses an assignment expression from the token stream.
//
// Steps:
//  1. First, parse the left-hand side (LHS) as an equality expression.
//     This ensures proper precedence, so assignment has lower precedence
//     than equality and arithmetic operators.
//  2. If the next token is an '=' (ASSIGN), then:
//     - Recursively call `assignment` to parse the right-hand side (RHS).
//     - Check if the LHS is a valid assignment target:
//     * If it's a Variable, produce an Assign AST node with the variable name
//     and the parsed RHS expression.
//     * Otherwise, produce a syntax error, since only variables can be assigned.
//  3. If no '=' follows, just return the previously parsed equality expression
//     as the result.
//
// Returns:
//   - Expression: Either an Assign node (for valid assignment expressions) or
//     the underlying expression if no assignment is found.
//   - error: Parsing errors such as invalid assignment targets or failed parsing of sub-expressions.
//
// Example:
// Input:  x = 10
// AST:    Assign{Name: x, Value: Literal(10)}
func (parser *Parser) assignment() (ast.Expression, error) {
	expression, err := parser.or()
	if err != nil {
		return nil, err
	}
	if parser.isMatch([]token.TokenType{token.ASSIGN}) {
		equalsToken := parser.previous()
		value, err := parser.assignment()
		if err != nil {
			return nil, err
		}
		switch v := expression.(type) {
		case ast.Variable:
			name := v.Name
			return ast.Assign{Name: name, Value: value}, nil

		default:
			msg := "Invalid assignment"
			return nil, CreateSyntaxError(equalsToken.Line, equalsToken.Column, msg)
		}
	}

	return expression, nil
}

// or parses a logical OR expression from the token stream.
// It first parses an AND expression on the left side, then consumes
// any sequence of OR operators, building a left-associative AST of logical expressions.
// Returns:
//   - ast.Expression: The constructed ast.Expression node
//   - error: An error if parsing fails.
func (parser *Parser) or() (ast.Expression, error) {
	expr, err := parser.and()
	if err != nil {
		return nil, err
	}

	for parser.isMatch([]token.TokenType{token.OR}) {
		op := parser.previous()
		rightExpr, err := parser.and()
		if err != nil {
			return nil, err
		}
		expr = ast.Logical{
			Left:     expr,
			Operator: op,
			Right:    rightExpr,
		}
	}

	return expr, nil
}

// and parses a logical AND expression from the token stream.
// It first parses an equality expression on the left side,
// then consumes any sequence of AND operators, building a left-associative
// abstract syntax tree (AST) of logical expressions.
// Returns:
//   - ast.Expression: The constructed ast.Expression node
//   - error: An error if parsing fails.
func (parser *Parser) and() (ast.Expression, error) {
	expr, err := parser.equality()
	if err != nil {
		return nil, err
	}

	for parser.isMatch([]token.TokenType{token.AND}) {
		op := parser.previous()
		rightExpr, err := parser.equality()
		if err != nil {
			return nil, err
		}

		expr = ast.Logical{
			Left:     expr,
			Operator: op,
			Right:    rightExpr,
		}
	}
	return expr, nil
}

// equality parses equality expressions using operators "==" and "!=".
//
// Returns:
//   - Expression: a Binary node (or sub-expression) representing equality comparison.
//   - error: if parsing fails.
func (parser *Parser) equality() (ast.Expression, error) {
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
		exp = ast.Binary{
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
func (parser *Parser) comparison() (ast.Expression, error) {
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
		exp = ast.Binary{
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
func (parser *Parser) term() (ast.Expression, error) {
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
		exp = ast.Binary{
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
func (parser *Parser) factor() (ast.Expression, error) {
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
		exp = ast.Binary{
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
func (parser *Parser) unary() (ast.Expression, error) {
	if parser.isMatch(unaryExpressionTypes) {
		operator := parser.previous()
		right, err := parser.unary()
		if err != nil {
			return nil, err
		}
		return ast.Unary{
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
func (parser *Parser) primary() (ast.Expression, error) {
	if parser.isMatch([]token.TokenType{token.FALSE}) {
		return ast.Literal{Value: false}, nil
	}
	if parser.isMatch([]token.TokenType{token.NULL}) {
		return ast.Literal{Value: nil}, nil
	}
	if parser.isMatch([]token.TokenType{token.TRUE}) {
		return ast.Literal{Value: true}, nil
	}

	if parser.isMatch([]token.TokenType{token.FLOAT, token.INT, token.STRING}) {
		return ast.Literal{Value: parser.previous().Literal}, nil
	}

	if parser.isMatch([]token.TokenType{token.IDENTIFIER}) {
		return ast.Variable{Name: parser.previous()}, nil
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
		return ast.Grouping{Expression: expr}, nil
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
