package lexer

import (
	"fmt"
	"nilan/token"
)

const (
	COMMENT_CHAR = '#'
)

type Lexer struct {
	Input        string
	tokens       []token.Token
	position     int
	readPosition int
	lineCount    int
}

func isLetter(char byte) bool {
	return 'a' <= char && char <= 'z' || 'A' <= char && char <= 'Z' || char == '_'
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
		Input:        input,
		lineCount:    0,
		position:     0,
		readPosition: 0,
	}
}

func (lexer *Lexer) advance() {
	lexer.position = lexer.readPosition
	lexer.readPosition++
}

// Determines of the lexer has finished scanning all the source code.
//
// Returns:
//   - bool: true if the lexer has finished scanning, false otherwise
func (lexer *Lexer) isFinished() bool {
	return lexer.readPosition >= len(lexer.Input)
}

// Gets the character at the lexer's current position
//
// Returns:
//   - byte: The character at the lexer's current position.
func (lexer *Lexer) readChar() byte {

	if lexer.isFinished() {
		return 0
	}

	char := lexer.Input[lexer.readPosition]
	lexer.advance()

	return char
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
	return lexer.Input[lexer.readPosition]
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

		// could be wrapped in to a method
		lexer.advance()
	}

	return true
}

func (lexer *Lexer) handleIdentifier() {

	initPos := lexer.position
	for {
		result := lexer.peek()
		if result == 0 || result == '\n' || !isLetter(result) {
			break
		}

		// could be wrapped in to a method
		lexer.position = lexer.readPosition
		lexer.readPosition++
	}

	substring := lexer.Input[initPos:lexer.readPosition]
	lexeme := token.CreateLiteralToken(token.IDENTIFIER, substring)
	lexer.tokens = append(lexer.tokens, lexeme)

}

// Determines if the next character in the source code
// matches the `expected` character.
func (lexer *Lexer) isMatch(expected byte) bool {

	nextIndex := lexer.readPosition
	if nextIndex >= len(lexer.Input) {
		return false
	}

	if lexer.Input[nextIndex] == expected {
		lexer.readPosition++
		return true
	}
	return false

}

func (lexer *Lexer) isWhiteSpace(char byte) bool {

	if char == ' ' || char == '\r' || char == '\t' {
		return true
	}
	if char == '\n' {
		lexer.lineCount++
		return true
	}
	return false
}

// Processes the current character and creates a token if applicable.
//
// This method is responsible for identifying and creating tokens based on the current
// character in the input stream.
//
// Returns:
//   - error: An error if an unexpected character is encountered, nil otherwise.
func (lexer *Lexer) scanToken() error {

	char := lexer.readChar()
	if lexer.isWhiteSpace(char) {
		return nil
	}
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
		fmt.Println(tok)
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

	case COMMENT_CHAR:
		if lexer.handleComment(char) {
			return nil
		}
	default:
		if isLetter(char) {
			lexer.handleIdentifier()
			return nil
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
	}
	lexer.tokens = append(lexer.tokens, token.CreateToken(token.EOF))
	return lexer.tokens, nil
}
