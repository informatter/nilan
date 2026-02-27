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

const MAX_PARAMETER_COUNT = 8

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
// `fromREPL` indicates whether the input is from a REPL session, which allows for more flexible parsing
// of statements (e.g., allowing expressions at the top level). If `fromREPL` is false, only declarations are allowed at
// the top level, and any non-declaration will produce a syntax error.
// This ensures that REPL sessions can execute arbitrary statements while regular src files enforce a stricter structure.
//
// Returns:
//   - []Stmt: the successfully parsed statements.
//   - []error: all errors that occurred during parsing.
func (parser *Parser) Parse(fromREPL bool) ([]ast.Stmt, []error) {
	statements := []ast.Stmt{}
	errors := []error{}

	for {
		if parser.isFinished() {
			break
		}
		statement, err := parser.parseDecl(!fromREPL)
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

// parseDecl parses either a top-level declaration or, when allowed, a statement.
// It inspects the next token and delegates to the appropriate method:
// - token.FUNC -> parseFuncDecl
// - token.VAR  -> variableDeclaration
// If topLevel is true, only declarations are accepted and any other token
// produces a syntax error.
// If topLevel is false and no declaration is found, the input is parsed as a statement.
// Returns the parsed ast.Stmt on success or an error describing the parse failure.
//
// Example: For input `fn foo() { print "hello" }`, parseDecl will recognize the `fn` token and
// delegate to parseFuncDecl to parse the function declaration and body, returning a FunctionStmt AST node.
// For input `foo()`, if topLevel is false, parseDecl will not find a declaration token, so it will delegate to
// statement() to parse the function call statement, returning a CallStmt AST node.
// If topLevel is true and the input is `foo()`, parseDecl will produce a syntax error since it expects a declaration at
// the top level. This ensures that function calls and other statements are not allowed at the top level.
// So programs such as:
// ```
// fn foo() { print "hello" }
// foo()
// ```
// will produce a syntax error for the `foo()` statement, since it's not a declaration.
//
// However, programs such as:
// ```
// fn foo() { print "hello" }
// var x = 10
// ```
// or
// ```
// fn foo() { print "hello" }
// fn main() { foo() }
// ```
// will be successfully parsed, since they only contain declarations at the top level.
func (parser *Parser) parseDecl(topLevel bool) (ast.Stmt, error) {

	if parser.isMatch([]token.TokenType{token.FUNC}) {
		return parser.parseFuncDecl()
	}
	if parser.isMatch([]token.TokenType{token.VAR}) {
		return parser.variableDeclaration()
	}
	// NOTE/TODO: ℹ️ Uncomment this once I have the necessary infra to enforce this syntax rule.
	// For exmample, if this is enabled, running `hellow_world.ni` breaks as it needs to be wrapped
	// within a function declaration which should get automatically executed during runtime.
	// For thos this is commented out to not block the function implementation development.  After implementing v1 of for functions (no closuers, no recursion, no first class functions)
	// Ineed to re-enable this and re-write test cases and the `hellow_world.ni` program to conform to this syntax rule.

	// if topLevel {
	// 	tok := parser.peek()
	// 	msg := fmt.Sprintf("expected declaration, found '%s'", tok.Lexeme)
	// 	return nil, CreateSyntaxError(tok.Line, tok.Column, msg)
	// }

	return parser.statement()
}

// parseFuncDecl parses a function declaration at the current parser position.
// If any expected token is missing or block parsing fails, it returns an error.
func (parser *Parser) parseFuncDecl() (ast.Stmt, error) {

	// TODO: Improve `consume` method so error handling is done in a more
	// encapsulated way, so we dont have to check for errors after every call
	// to `consume`

	name, err := parser.consume(token.IDENTIFIER, "expected function name")
	if err != nil {
		return nil, err
	}
	_, err = parser.consume(token.LPA, "expected '(' after function name")
	if err != nil {
		return nil, err
	}
	parameters := make([]token.Token, 0)

	if !parser.checkType(token.RPA) {
		for {
			if len(parameters) >= MAX_PARAMETER_COUNT {
				msg := fmt.Sprintf("A method can't have more than %d parameters.", MAX_PARAMETER_COUNT)
				return nil, CreateSyntaxError(name.Line, name.Column, msg)
			}
			param, err := parser.consume(token.IDENTIFIER, "expected parameter name")
			if err != nil {
				return nil, err
			}
			parameters = append(parameters, param)

			if !parser.isMatch([]token.TokenType{token.COMMA}) {
				break
			}
		}
	}
	_, err = parser.consume(token.RPA, "expected ')' after parameters")
	if err != nil {
		return nil, err
	}
	_, err = parser.consume(token.LCUR, "expected '{' before function body")
	if err != nil {
		return nil, err
	}
	bodyStmts, err := parser.block()
	if err != nil {
		return nil, err
	}
	return ast.FunctionStmt{
		Name:       name,
		Parameters: parameters,
		Body: ast.BlockStmt{
			Statements: bodyStmts,
		},
	}, nil
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

	var initialiser ast.Expr
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
		stmt, err := parser.parseDecl(false)
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
func (parser *Parser) expression() (ast.Expr, error) {
	return parser.assignment()
}

// parseCallExpr parses call expressions starting from a primary expression.
// Returns the resulting ast.Expr or an error encountered while parsing.
//
// Example: For input `foo(1, 2)`, parseCallExpr will first parse `foo` as a VariableExpr,
// then recognize the following `(` as the start of a call expression, and delegate to finishCall
// to parse the arguments `1` and `2` and construct a CallExpr AST node.
func (parser *Parser) parseCallExpr() (ast.Expr, error) {
	expr, err := parser.parsePrimaryExpr()
	if err != nil {
		return nil, err
	}

	// If the next token is a left parenthesis, we have a call expression. We delegate to finishCall to
	// parse the arguments and update the expression.
	// NOTE: Currently, This only supports parsing call expressions such as foo(1, 2) or foo() but not method calls such as
	// x.bar() or x.bar(1, 2), or chained calls `foo()(1)(2)`. To allow chaining, we would need to repeatedly check for a
	// left parenthesis after parsing the call expression, and delegate to finishCall each time.
	if parser.isMatch([]token.TokenType{token.LPA}) {
		expr, err = parser.finishCall(expr)
		if err != nil {
			return nil, err
		}
	}

	return expr, nil
}

// finishCall parses an expression call's argument list for the given callee and
// returns an ast.CallExpr or an error. It supports 0 or more comma-
// separated arguments.
func (parser *Parser) finishCall(callee ast.Expr) (ast.Expr, error) {
	args := make([]ast.Expr, 0)

	if !parser.checkType(token.RPA) {
		for {
			arg, err := parser.expression()
			if err != nil {
				return nil, err
			}
			args = append(args, arg)

			if !parser.isMatch([]token.TokenType{token.COMMA}) {
				break
			}
		}
	}

	paren, err := parser.consume(token.RPA, "expected ')' after arguments")
	if err != nil {
		return nil, err
	}

	return ast.CallExpr{
		Callee:    callee,
		Paren:     paren,
		Arguments: args,
	}, nil
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
func (parser *Parser) assignment() (ast.Expr, error) {
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
		case ast.VariableExpr:
			name := v.Name
			return ast.AssignExpr{Name: name, Value: value}, nil

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
func (parser *Parser) or() (ast.Expr, error) {
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
		expr = ast.LogicalExpr{
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
func (parser *Parser) and() (ast.Expr, error) {
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

		expr = ast.LogicalExpr{
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
func (parser *Parser) equality() (ast.Expr, error) {
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
		exp = ast.BinaryExpr{
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
func (parser *Parser) comparison() (ast.Expr, error) {
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
		exp = ast.BinaryExpr{
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
func (parser *Parser) term() (ast.Expr, error) {
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
		exp = ast.BinaryExpr{
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
func (parser *Parser) factor() (ast.Expr, error) {
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
		exp = ast.BinaryExpr{
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
func (parser *Parser) unary() (ast.Expr, error) {
	if parser.isMatch(unaryExpressionTypes) {
		operator := parser.previous()
		right, err := parser.unary()
		if err != nil {
			return nil, err
		}
		return ast.UnaryExpr{
			Operator: operator,
			Right:    right,
		}, nil
	}
	return parser.parseCallExpr()
}

// parsePrimaryExpr parses the most basic forms of expressions:
//   - Literals: true, false, null, strings, numbers
//   - Grouping: (expression)
//
// If no valid token matches, returns a syntax error.
//
// Returns:
//   - Expression: a Literal, Grouping expression .
//   - error: if no valid primary expression can be parsed.
func (parser *Parser) parsePrimaryExpr() (ast.Expr, error) {
	if parser.isMatch([]token.TokenType{token.FALSE}) {
		return ast.LiteralExpr{Value: false}, nil
	}
	if parser.isMatch([]token.TokenType{token.NULL}) {
		return ast.LiteralExpr{Value: nil}, nil
	}
	if parser.isMatch([]token.TokenType{token.TRUE}) {
		return ast.LiteralExpr{Value: true}, nil
	}

	if parser.isMatch([]token.TokenType{token.FLOAT, token.INT, token.STRING}) {
		return ast.LiteralExpr{Value: parser.previous().Literal}, nil
	}

	if parser.isMatch([]token.TokenType{token.IDENTIFIER}) {
		return ast.VariableExpr{Name: parser.previous()}, nil
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
		return ast.GroupingExpr{Expression: expr}, nil
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
