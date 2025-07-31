package lexer

import (
	"fmt"
	"nilan/token"
	"strconv"
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

func isNumber(char byte) bool {
	return '0' <= char && char <= '9'
}

func convertToInt(s string) (int, error) {
	num, err := strconv.Atoi(s)
	if err != nil {
		return 0, err
	}
	return num, nil
}
func convertTofloat64(s string) (float64, error) {
	num, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, err
	}
	return num, nil
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

func (lexer *Lexer) peekNext() byte {
	if lexer.readPosition+1 >= len(lexer.Input) {
		return 0
	}
	return lexer.Input[lexer.readPosition+1]
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

// handleNumber processes and tokenizes a number in the input string.
// It handles both integer and floating-point numbers, including negative numbers.
//
// The method scans the input from the current position, identifying valid number formats.
// It supports the following number formats:
//
//   - Integers: e.g., 123, -456
//   - Floating-point numbers: e.g., 3.14, -0.5
//
// Returns an error if an invalid number format is encountered.
func (lexer *Lexer) handleNumber() error {
	initPos := lexer.position
	decimalCount := 0
	negativeCount := 0
	isNegative := false
	if lexer.Input[lexer.position] == '-' {
		isNegative = true
	}
	for {
		result := lexer.peek()
		if result == 0 || result == '\n' || !isNumber(result) && result != '.' && result != '-' {
			break
		}
		if result == '.' {
			// handles numbers such as 1.
			if lexer.peekNext() == 0 {
				return fmt.Errorf("invalid number in line: %v", lexer.lineCount)
			}
			// handles numbers such as 1.1.
			if decimalCount == 1 {
				return fmt.Errorf("invalid number in line: %v", lexer.lineCount)
			}
			decimalCount++
		}
		if result == '-' {
			if isNegative {
				return fmt.Errorf("invalid number in line: %v", lexer.lineCount)
			}

			pNextResult := lexer.peekNext()
			// handles numbers such as 2-2 or 2-!
			if pNextResult == 0 || isNumber(pNextResult) || !isNumber(pNextResult) {
				return fmt.Errorf("invalid number in line: %v", lexer.lineCount)
			}

			if negativeCount == 1 {
				return fmt.Errorf("invalid number in line: %v", lexer.lineCount)
			}
			negativeCount++

		}

		lexer.advance()
	}
	substring := lexer.Input[initPos:lexer.readPosition]
	var tokenType token.TokenType
	if decimalCount == 0 {
		tokenType = token.INT
	} else {
		tokenType = token.FLOAT
	}
	lexer.tokens = append(lexer.tokens, token.CreateLiteralToken(tokenType, substring))

	return nil
}

// handleIdentifier processes a user identifier or a
// language keyword in the source code.
func (lexer *Lexer) handleIdentifier() {

	initPos := lexer.position
	for {
		result := lexer.peek()
		if result == 0 || result == '\n' || !isLetter(result) {
			break
		}
		lexer.advance()
	}

	substring := lexer.Input[initPos:lexer.readPosition]
	lexeme := token.Token{
		TokenType: token.IDENTIFIER,
		Value:     substring,
	}

	if keywordType, exists := token.KeyWords[lexeme.Value]; exists {
		lexeme.TokenType = keywordType
	}

	lexer.tokens = append(lexer.tokens, lexeme)
}

// handleStringLiteral processes string literals in the input.
//
// Returns:
//   - nil if the string literal is properly closed and processed
//   - error if the string literal is unclosed or has new lines
func (lexer *Lexer) handleStringLiteral() error {

	initPos := lexer.position
	isClosed := false
	for {
		result := lexer.peek()
		if result == 0 {
			break
		}

		lexer.advance()
		if result == '"' {
			isClosed = true
			break
		}
	}

	if !isClosed {
		return fmt.Errorf("unclosed string literal: %s\nline: %v", lexer.Input[initPos+1:lexer.readPosition], lexer.lineCount)
	}
	substring := lexer.Input[initPos+1 : lexer.position]
	lexer.tokens = append(lexer.tokens, token.CreateLiteralToken(token.STRING, substring))
	return nil
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
		if isNumber(lexer.peek()) {
			err := lexer.handleNumber()
			if err != nil {
				return err
			}
			return nil
		}
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
	case '"':
		err := lexer.handleStringLiteral()
		if err != nil {
			return err
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
		if isNumber(char) {
			err := lexer.handleNumber()
			if err != nil {
				return err
			}
			return nil
		}

		return fmt.Errorf("unexpected character: %c\nline: %v", char, lexer.lineCount)
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
