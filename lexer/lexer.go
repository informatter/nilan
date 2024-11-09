package lexer

import (
	"fmt"
	"nilan/token"
)

type Lexer struct {
	Input       string
	tokens      []token.Token
	start       int
	current     int
	currentChar byte
}

func CreateLexer(input string) *Lexer {
	return &Lexer{
		Input:   input,
		start:   0,
		current: 0,
	}
}

func isWhiteSpace(char byte) bool {
	if char == ' ' {
		return true
	}
	return false
}

func (lexer *Lexer) isAtEnd() bool {
	return lexer.current >= len(lexer.Input)
}

// Returns the next character in the source code.
func (lexer *Lexer) nextChar() byte {

	char := lexer.Input[lexer.current]
	lexer.current++
	return char

}

func (lexer *Lexer) match(expected byte) bool {

	if lexer.isAtEnd() {
		return false
	}
	if lexer.Input[lexer.current] != expected {
		return false
	}
	lexer.current++
	return true
}

func (lexer *Lexer) scanToken() error {

	char := lexer.nextChar()

	var tok token.Token
	switch char {
	case '(':
		tok = token.CreateToken(token.LPA)
	case ')':
		tok = token.CreateToken(token.RPA)
	case '{':
		tok = token.CreateToken(token.LCUR)
	case '}':
		tok = token.CreateToken(token.RCUR)
	case ';':
		tok = token.CreateToken(token.SEMICOLON)
	case ',':
		tok = token.CreateToken(token.COMMA)
	case '=':
		tok = token.CreateToken(token.ASSIGN)
	case '*':
		tok = token.CreateToken(token.MULT)
	case '+':
		tok = token.CreateToken(token.ADD)
	case '-':
		tok = token.CreateToken(token.SUB)
	case '/':
		tok = token.CreateToken(token.DIV)

	case '!':
		tok = token.CreateToken(token.BANG)
		if lexer.match('=') {
			tok = token.CreateToken(token.NOT_EQUAL)
		}
	case '<':
		tok = token.CreateToken(token.LESS)
		if lexer.match('=') {
			tok = token.CreateToken(token.LESS_EQUAL)
		}
	default:
		return fmt.Errorf("unexpected character: %c", char)
	}

	lexer.tokens = append(lexer.tokens, tok)
	return nil
}

// Creates a list of Token's from the source code
func (lexer *Lexer) Scan() ([]token.Token, error) {

	for !lexer.isAtEnd() {
		lexer.start = lexer.current
		err := lexer.scanToken()
		if err != nil {
			return lexer.tokens, err
		}
	}
	lexer.tokens = append(lexer.tokens, token.CreateToken(token.EOF))
	return lexer.tokens, nil
}
