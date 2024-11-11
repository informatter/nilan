package lexer

import (
	"fmt"
	"nilan/token"
)

const (
	COMMENT_CHAR = '#'
)

type Lexer struct {
	Input     string
	tokens    []token.Token
	position  int
	lineCount int
}

// Initializes and returns a new Lexer instance.
//
// Parameters:
//   - input: string
//     The the source code as a string to be lexically analyzed.
//
// Returns:
//   - *Lexer: A pointer to a newly created Lexer instance.
func CreateLexer(input string) *Lexer {
	return &Lexer{
		Input:     input,
		lineCount: 0,
	}
}

// Determines of the lexer has finished scanning all the source code.
//
// Returns:
//   - bool: true if the lexer has finished scanning, false otherwise
func (lexer *Lexer) isFinished() bool {
	return lexer.position >= len(lexer.Input)
}

// Increments the lexer's position by one unit
func (lexer *Lexer) advance() {
	lexer.position++
}

// Gets the character at the lexer's current position
//
// Returns:
//   - byte: The character at the lexer's current position.
func (lexer *Lexer) getCurrentChar() byte {

	return lexer.Input[lexer.position]
}

// Returns the next character in the input without advancing the lexer's position.
//
// This method allows the lexer to look ahead at the next character in the input stream
// without consuming it.
//
// Returns:
//   - byte: The next character in the input stream.
//     If the lexer has reached the end of the input, it returns 0 (null byte )
func (lexer *Lexer) peek() byte {
	if lexer.isFinished() {
		return 0
	}
	return lexer.Input[lexer.position]
}

// handleComment processes a comment in the input stream.
//
// This method is responsible for handling comments in the lexical analysis.
// It checks if the current character is a comment character and, if so,
// consumes all characters until the end of the line or end of input.
//
// Parameters:
//   - char: The current character being processed.
//
// Returns:
//   - bool: true if a comment was processed, false otherwise.
func (lexer *Lexer) handleComment(char byte) bool {
	if char != COMMENT_CHAR {
		return false
	}

	for {
		result := lexer.peek()
		if result == 0 || result == '\n' {
			break
		}
		lexer.advance()
	}

	return true
}

// Determines if the next character in the source code
// matches the `expected` character.
func (lexer *Lexer) isMatch(expected byte) bool {

	if lexer.isFinished() {
		return false
	}
	nextIndex := lexer.position + 1
	if nextIndex >= len(lexer.Input) {
		return false
	}
	if lexer.Input[nextIndex] != expected {
		return false
	}
	lexer.advance()
	return true
}

// Processes the current character and creates a token if applicable.
//
// This method is responsible for identifying and creating tokens based on the current
// character in the input stream.
//
// Returns:
//   - error: An error if an unexpected character is encountered, nil otherwise.
func (lexer *Lexer) scanToken() error {

	char := lexer.getCurrentChar()
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
	case '*':
		tok = token.CreateToken(token.MULT)
	case '+':
		tok = token.CreateToken(token.ADD)
	case '-':
		tok = token.CreateToken(token.SUB)
	case '/':
		tok = token.CreateToken(token.DIV)
	case '=':
		tok = token.CreateToken(token.ASSIGN)
		if lexer.isMatch('=') {
			tok = token.CreateToken(token.EQUAL_EQUAL)
		}
	case '!':
		tok = token.CreateToken(token.BANG)
		if lexer.isMatch('=') {
			tok = token.CreateToken(token.NOT_EQUAL)
		}
	case '<':
		tok = token.CreateToken(token.LESS)
		if lexer.isMatch('=') {
			tok = token.CreateToken(token.LESS_EQUAL)
		}
	case '>':
		tok = token.CreateToken(token.LARGER)
		if lexer.isMatch('=') {
			tok = token.CreateToken(token.LARGER_EQUAL)
		}

	case '\n':
		lexer.lineCount++
	case ' ': // ignores empty space
	case '\r': // ingores carriage return
	case '\t': // ignores tab
		break
	default:

		result := lexer.handleComment(char)
		if result {
			break
		}
		return fmt.Errorf("unexpected character: %c", char)
	}
	if tok.Value != "" {
		lexer.tokens = append(lexer.tokens, tok)
	}

	return nil
}

// Scan performs lexical analysis on the input and returns a slice of tokens.
//
// This method is the main entry point for the lexical analysis process. It iterates
// through the input, tokenizing it and collecting all tokens until the end of the input
// is reached or an error occurs.
//
// Returns:
//   - []token.Token: A slice containing all tokens found in the input.
//   - error: An error if any issues occurred during scanning, or nil if successful.
func (lexer *Lexer) Scan() ([]token.Token, error) {

	for !lexer.isFinished() {
		err := lexer.scanToken()
		if err != nil {
			return lexer.tokens, err
		}
		lexer.advance()
	}
	lexer.tokens = append(lexer.tokens, token.CreateToken(token.EOF))
	return lexer.tokens, nil
}
