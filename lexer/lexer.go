package lexer

import (
	"fmt"
	"nilan/token"
	"strconv"
)

const (
	COMMENT_CHAR = '#'
)

func isLetter(char rune) bool {
	return rune('a') <= char && char <= rune('z') || rune('A') <= char && char <= rune('Z') || char == rune('_')
}

func isNumber(char rune) bool {
	return rune('0') <= char && char <= rune('9')
}

// Lexer represents a lexical scanner for processing input text into tokens.
// It maintains the current scanning state, including the position within the
// input, the current character, and metadata for line/column tracking.
// The Lexer also records tokens and errors encountered during scanning.
type Lexer struct {
	// rune slice of the input string being scanned.
	characters []rune

	// Total number of runes in the input.
	totalChars int

	// Stores the sequence of tokens produced during lexing.
	tokens []token.Token

	// The index of the character that was previously read
	position int

	// The current character being examined.
	currentChar rune

	// The index of the next position where the next character
	// will be read
	readPosition int

	// Tracks the number of lines processed (incremented on newline).
	lineCount int32

	// Tracks the character's position within the current line.
	// Gets reset on every new line back to 0
	column int

	// Stores any scanning errors that occur during lexing.
	errors []error
}

// Initializes and returns a new Lexer instance.
//
// Parameters:
//   - input: string
//     The the source code as a string to be lexically analyzed.
//
// Returns:
//   - *Lexer: A pointer to a newly created Lexer instance.
func New(input string) *Lexer {
	lexer := &Lexer{
		characters: []rune(input),
	}
	lexer.totalChars = len(lexer.characters)
	lexer.readChar()
	return lexer
}

// Updates the `Lexer`'s reading position forward by one character.
//
// Behavior:
//   - Sets `position` to the current `readPosition“
//   - Increments `readPosition` by 1, so the lexer is ready to read the next
//     character on the following call.
//   - Updates the `column` to match `readPosition`, keeping track of the
//     character's position within the line.
func (lexer *Lexer) advance() {
	lexer.position = lexer.readPosition
	lexer.readPosition++
	lexer.column = lexer.readPosition
}

// Determines of the lexer has finished scanning all the source code.
//
// Returns:
//   - bool: true if the lexer has finished scanning, false otherwise
func (lexer *Lexer) isFinished() bool {
	return lexer.readPosition >= lexer.totalChars
}

// Reads the character at the `Lexer`'s `readPosition`. If there
// are no more characters to parse, it sets the `Lexer`'s current
// character to null.
func (lexer *Lexer) readChar() {

	if lexer.isFinished() {
		lexer.currentChar = rune(0)
	} else {
		lexer.currentChar = lexer.characters[lexer.readPosition]
	}
	lexer.advance()
}

// Reads a sequence of characters from the input until a whitespace
// character or end-of-file marker (rune(0)) is encountered. This method is
// typically used to capture tokens or substrings that do not match any valid
// lexical category (i.e., "illegal" tokens).
//
// Parameters:
//   - startPos (int): The index in the character slice where the illegal token begins.
//
// Returns:
//   - string: The substring of characters between startPos (inclusive) and the
//     current read position, representing the
//     illegal token.
func (lexer *Lexer) readIllegal(startPos int) string {
	for !lexer.isWhiteSpace(lexer.currentChar) && !lexer.isFinished() {
		lexer.readChar()
	}
	// return string(lexer.characters[startPos:lexer.readPosition -1])
	return string(lexer.characters[startPos:lexer.readPosition])

}

// Returns the character at the `Lexer`s `readPosition` without consiming the character
//
// Returns:
//   - rune: The next character in the input stream.
//     If the lexer has reached the end of the input, it returns 0 (null)
func (lexer *Lexer) peek() rune {
	if lexer.isFinished() {
		return rune(0)
	}
	return lexer.characters[lexer.readPosition]
}

// Returns the next character from the `Lexer`'s `readPosition` without consiming the character
// Returns:
//   - rune: The next character in the input stream.
//     If the lexer has reached the end of the input, it returns 0 (null)
func (lexer *Lexer) peekNext() rune {
	nextReadPos := lexer.readPosition + 1
	if nextReadPos >= lexer.totalChars {
		return rune(0)
	}
	return lexer.characters[nextReadPos]
}

// handleComment processes a comment in the input stream.
//
// This method is responsible for handling comments in the lexical analysis.
// It checks if the current character is a comment character and, if so,
// consumes all characters until the end of the line or end of input,
// while advancing the `Lexer`'s position
func (lexer *Lexer) handleComment() {
	for lexer.currentChar != rune('\n') && !lexer.isFinished() {
		lexer.readChar()
	}
}

// handleNumber scans a sequence of digits (and at most one decimal point) from
// the input and creates an integer or floating-point literal token accordingly.
//
// The method starts scanning from the current lexer position and continues
// advancing until it encounters a character that is not a digit or a decimal
// point (`.`). A decimal point is allowed only once within the number.
//
// Validation rules:
//   - A number ending with a decimal point (e.g., "1.") without further digits
//     results in an error.
//   - Multiple decimal points (e.g., "1.1.") are considered invalid and cause
//     an error.
//
// Returns:
//   - nil if the token was successfully created and added
//   - an error if the number format is invalid
func (lexer *Lexer) handleNumber() error {
	initPos := lexer.position
	decimalCount := 0

	for {
		nextChar := lexer.peek()
		if nextChar == rune(0) || nextChar == rune('\n') || !isNumber(nextChar) && nextChar != rune('.') {
			break
		}
		if nextChar == '.' {
			// handles numbers such as 1.
			if lexer.peekNext() == rune(0) {
				illegalNumber := string(lexer.characters[initPos : lexer.readPosition+1])
				return fmt.Errorf("invalid number: '%s', line: %v", string(illegalNumber), lexer.lineCount)
			}
			// handles numbers such as 1.1.
			if decimalCount == 1 {
				illegalNumber := lexer.readIllegal(initPos)
				return fmt.Errorf("invalid number: '%s', line: %v", string(illegalNumber), lexer.lineCount)

			}
			decimalCount++
		}
		// handles numbers such as .2
		if lexer.currentChar == rune('.') && isNumber(nextChar) {
			decimalCount++
		}

		lexer.advance()
	}
	number := string(lexer.characters[initPos:lexer.readPosition])
	var tok token.Token

	if decimalCount == 0 {
		result, _ := strconv.ParseInt(number, 0, 64)
		tok = token.CreateLiteralToken(token.INT, result, number, lexer.lineCount, lexer.column)
	} else {
		result, _ := strconv.ParseFloat(number, 64)
		tok = token.CreateLiteralToken(token.FLOAT, result, number, lexer.lineCount, lexer.column)
	}
	lexer.tokens = append(lexer.tokens, tok)

	return nil
}

// handleIdentifier processes a user identifier or a
// language keyword in the source code.
func (lexer *Lexer) handleIdentifier() {

	initPos := lexer.position
	for {
		result := lexer.peek()
		if result == rune(0) || result == rune('\n') || !isLetter(result) {
			break
		}
		lexer.advance()
	}

	identifier := lexer.characters[initPos:lexer.readPosition]
	lexeme := token.Token{
		TokenType: token.IDENTIFIER,
		Lexeme:    string(identifier),
	}

	if keywordType, exists := token.KeyWords[lexeme.Lexeme]; exists {
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
		return fmt.Errorf("unclosed string literal: '%s', line: %v", string(lexer.characters[initPos+1:lexer.readPosition]), lexer.lineCount)
	}

	// NOTE: `initPos+1`` and `lexer.position` is to ignore escape characters.
	// as we dont need to store them for a literal string token
	// "\"foo\"" -> "foo"
	stringLiteral := string(lexer.characters[initPos+1 : lexer.position])
	lexer.tokens = append(lexer.tokens, token.CreateLiteralToken(token.STRING, stringLiteral, stringLiteral, lexer.lineCount, lexer.column))
	return nil
}

// Determines if the next character in the source code
// matches the `expected` character.
func (lexer *Lexer) isMatch(expected rune) bool {

	if lexer.isFinished() {
		return false
	}

	if lexer.characters[lexer.readPosition] == expected {
		lexer.readPosition++
		return true
	}
	return false

}

// isWhiteSpace determines whether a given rune represents whitespace in the input stream.
// In Nilan, whitespace is considered to be the following characters:
//   - carriage return ('\r')
//   - tab ('\t')
//   - newline ('\n')
//   - ASCII space (' ')
//
// Parameters:
//   - char (rune): The character being evaluated.
//
// Returns:
//   - bool: true if the character is considered whitespace, otherwise false.
func (lexer *Lexer) isWhiteSpace(char rune) bool {

	if char == rune(' ') || char == rune('\r') || char == rune('\t') {
		return true
	}
	if lexer.currentChar == rune('\n') {
		// increment line count and reset column back to zero
		lexer.lineCount++
		lexer.column = 0
		return true
	}
	return false
}

// Skips all whitespaces in the input while advancing the `Lexer`'s position
func (lexer *Lexer) skipWhiteSpace() {
	for lexer.isWhiteSpace(lexer.currentChar) {
		lexer.readChar()
	}
}

// Processes the current character and creates a token if applicable.
//
// This method is responsible for identifying and creating tokens based on the current
// character in the input stream.
func (lexer *Lexer) createToken() {

	lexer.skipWhiteSpace()

	switch lexer.currentChar {
	case rune('('):
		tok := token.CreateToken(token.LPA, lexer.lineCount, lexer.column)
		lexer.tokens = append(lexer.tokens, tok)
	case rune(')'):
		tok := token.CreateToken(token.RPA, lexer.lineCount, lexer.column)
		lexer.tokens = append(lexer.tokens, tok)
	case rune('{'):
		tok := token.CreateToken(token.LCUR, lexer.lineCount, lexer.column)
		lexer.tokens = append(lexer.tokens, tok)
	case rune('}'):
		tok := token.CreateToken(token.RCUR, lexer.lineCount, lexer.column)
		lexer.tokens = append(lexer.tokens, tok)
	case rune(';'):
		tok := token.CreateToken(token.SEMICOLON, lexer.lineCount, lexer.column)
		lexer.tokens = append(lexer.tokens, tok)
	case rune(','):
		tok := token.CreateToken(token.COMMA, lexer.lineCount, lexer.column)
		lexer.tokens = append(lexer.tokens, tok)
	case rune('*'):
		tok := token.CreateToken(token.MULT, lexer.lineCount, lexer.column)
		lexer.tokens = append(lexer.tokens, tok)
	case rune('+'):
		tok := token.CreateToken(token.ADD, lexer.lineCount, lexer.column)
		lexer.tokens = append(lexer.tokens, tok)
	case rune('-'):
		tok := token.CreateToken(token.SUB, lexer.lineCount, lexer.column)
		lexer.tokens = append(lexer.tokens, tok)
	case rune('/'):
		tok := token.CreateToken(token.DIV, lexer.lineCount, lexer.column)
		lexer.tokens = append(lexer.tokens, tok)
	case rune('='):
		tok := token.CreateToken(token.ASSIGN, lexer.lineCount, lexer.column)
		if lexer.isMatch(rune('=')) {
			tok = token.CreateToken(token.EQUAL_EQUAL, lexer.lineCount, lexer.column)
		}
		lexer.tokens = append(lexer.tokens, tok)
	case rune('!'):
		tok := token.CreateToken(token.BANG, lexer.lineCount, lexer.column)
		if lexer.isMatch(rune('=')) {
			tok = token.CreateToken(token.NOT_EQUAL, lexer.lineCount, lexer.column)
		}
		lexer.tokens = append(lexer.tokens, tok)
	case rune('<'):
		tok := token.CreateToken(token.LESS, lexer.lineCount, lexer.column)
		if lexer.isMatch(rune('=')) {
			tok = token.CreateToken(token.LESS_EQUAL, lexer.lineCount, lexer.column)
		}
		lexer.tokens = append(lexer.tokens, tok)
	case rune('>'):
		tok := token.CreateToken(token.LARGER, lexer.lineCount, lexer.column)
		if lexer.isMatch(rune('=')) {
			tok = token.CreateToken(token.LARGER_EQUAL, lexer.lineCount, lexer.column)
		}
		lexer.tokens = append(lexer.tokens, tok)
	case rune('"'):
		err := lexer.handleStringLiteral()
		if err != nil {

			lexer.errors = append(lexer.errors, err)
		}

	case rune(COMMENT_CHAR):
		lexer.handleComment()
	default:
		if isLetter(lexer.currentChar) {
			lexer.handleIdentifier()
		} else if isNumber(lexer.currentChar) || lexer.currentChar == rune('.') {
			err := lexer.handleNumber()
			if err != nil {
				lexer.errors = append(lexer.errors, err)
			}
		} else if !lexer.isFinished() {

			position := lexer.position
			column := lexer.column
			currentChar := lexer.currentChar
			illegal := lexer.readIllegal(position)

			err := fmt.Errorf("unexpected character: '%c' in: '%s', line: %v, column: %v", currentChar, illegal, lexer.lineCount, column)
			lexer.errors = append(lexer.errors, err)
		}
	}

	lexer.readChar()
}

// Scan performs lexical analysis on the input and returns a slice of tokens.
//
// This method is the main entry point for the lexical analysis process. It iterates
// through the input, tokenizing it and collecting all tokens until the end of the input
// is reached or an error occurs.
//
// Returns:
//   - []token.Token: A slice containing all tokens found in the input.
//   - error: An error if any issues occurred during lexing, or nil if successful.
func (lexer *Lexer) Scan() ([]token.Token, error) {

	if lexer.totalChars > 1 {
		for lexer.currentChar != rune(0) {
			lexer.createToken()
			if len(lexer.errors) == 1 {
				return lexer.tokens, lexer.errors[0]
			}
		}
	} else {
		// special handling for inputs with a single character or empty inputs.
		lexer.createToken()
		if len(lexer.errors) == 1 {
			return lexer.tokens, lexer.errors[0]
		}
	}
	lexer.tokens = append(lexer.tokens, token.CreateToken(token.EOF, lexer.lineCount, lexer.column))
	return lexer.tokens, nil
}
