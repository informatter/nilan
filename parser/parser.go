package parser

import (
	"nilan/token"
)

// Recursive descent parser
// https://en.wikipedia.org/wiki/Recursive_descent_parser

//  A Recursive descent parser is a top-down parser because it starts from the top
// grammar rule and works is way down in to the netsed sub-experessions before reaching
// the leaves of the syntax tree (terminal rules)

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
//	 - position: int
//	   The position of the parser in respect to the current token being
//     looked at.
//
// Returns:
//   - *Parser: A pointer to a newly created Parser instance.
func CreateParser(tokens []token.Token) *Parser {
	return &Parser{
		tokens:   tokens,
		position: 0,
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

func (parser *Parser) expression() Expression {
	return parser.equality()
}

// equality = comparison { ("!=" | "==") comparison }
func (parser *Parser) equality() Expression {

	exp := parser.comparison()
	for parser.isMatch(equalityTokenTypes) {

		operator := parser.previous()
		right := parser.comparison()
		exp = Binary{
			Left:     exp,
			Operator: operator,
			Right:    right,
		}
	}
	return exp
}

// comparison = term { (">" | ">=" | "<" | "<=") term }
func (parser *Parser) comparison() Expression {

	exp := parser.term()
	for parser.isMatch(comparisonTokenTypes) {

		operator := parser.previous()
		right := parser.term()
		exp = Binary{
			Left:     exp,
			Operator: operator,
			Right:    right,
		}
	}
	return exp
}
