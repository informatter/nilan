package token

import (
	"fmt"
)

const (

	// special characters
	LPA       = "("
	RPA       = ")"
	COMMA     = ","
	SEMICOLON = ";"
	RCUR      = "}"
	LCUR      = "{"

	// naming given by programmer i.e myVar, myFunc, add ..ect
	IDENTIFIER = "IDENTIFIER"
	STRING     = "STRING"

	// operators
	ASSIGN       = "="
	MULT         = "*"
	ADD          = "+"
	SUB          = "-"
	DIV          = "/"
	NOT_EQUAL    = "!="
	EQUAL_EQUAL  = "=="
	BANG         = "!"
	LESS         = "<"
	LESS_EQUAL   = "<="
	LARGER       = ">"
	LARGER_EQUAL = ">="

	FLOAT = "FLOAT"
	INT   = "INT"

	EOF = "EOF"

	// keywords
	FUNC   = "FUNCTION"
	OR     = "OR"
	AND    = "AND"
	FOR    = "FOR"
	WHILE  = "WHILE"
	CONST  = "CONST"
	VAR    = "VAR"
	RETURN = "RETURN"
	IF     = "IF"
	ELSE   = "ELSE"
	ELIF   = "ELIF"
	BREAK  = "BREAK"
	TRUE   = "TRUE"
	FALSE  = "FALSE"
	NULL   = "NULL"
)

// KeyWords maps reserved keyword strings in Nilan to their
// corresponding TokenType values.
//
// During lexical analysis, if an identifier matches one of these strings,
// it is classified as the appropriate keyword token, rather than a generic
// identifier.
//
// Example:
//
//	tokenType, ok := KeyWords[lexeme]
//	if ok {
//	    // lexeme is a known keyword; use tokenType
//	} else {
//	    // lexeme is a regular identifier
//	}
var KeyWords = map[string]TokenType{
	"fn":     FUNC,
	"or":     OR,
	"and":    AND,
	"while":  WHILE,
	"for":    FOR,
	"var":    VAR,
	"const":  CONST,
	"return": RETURN,
	"if":     IF,
	"else":   ELSE,
	"elif":   ELIF,
	"break":  BREAK,
	"false":  FALSE,
	"true":   TRUE,
	"null":   NULL,
}

// tokenTypes maps single and multi-character symbols in Nilan
// to their corresponding TokenType values.
//
// The lexer uses this map when it encounters punctuation, operators, or other
// syntax characters, allowing it to quickly look up the correct token type.
//
// Note: This map is typically used alongside the KeyWords map to differentiate
// between symbols/operators and keywords.
var tokenTypes = map[TokenType]string{
	"EOF": EOF,
	"(":   LPA,
	")":   RPA,
	"{":   LCUR,
	"}":   RCUR,
	";":   SEMICOLON,
	",":   COMMA,
	"=":   ASSIGN,
	"*":   MULT,
	"+":   ADD,
	"-":   SUB,
	"/":   DIV,
	"!":   BANG,
	"!=":  NOT_EQUAL,
	"<":   LESS,
	"<=":  LESS_EQUAL,
	">":   LARGER,
	">=":  LARGER_EQUAL,
	"==":  EQUAL_EQUAL,
}

type TokenType string

// Token represents a lexical token identified during lexical analysis
// (tokenization) of a source file. It encapsulates the token's type, its
// original textual representation, any literal value it may hold, and its
// position within the source.
//
// Fields:
//   - TokenType: The category or classification of the token, such as
//     keyword, identifier, string, number, or symbol.
//   - Lexeme: The exact string from the source code that was matched to form
//     this token.
//   - Literal: The interpreted value of the token, if applicable.
//     For example, a number token might have an integer or float value here.
//     This is stored as `any`.
//   - Line: The source line (o-based index) where the token appears.
//   - Column: The character position (0-based index) within the line where
//     the token starts.
type Token struct {
	TokenType TokenType
	Lexeme    string
	Literal   any
	Line      int32
	Column    int
}

// CreateToken constructs and returns a new Token instance for the given
// token type and position within the source code.
//
// Parameters:
//   - tokenType: The classification of the token to create.
//   - line:      The 0-based line number where the token begins.
//   - column:    The 0-based column number where the token begins.
//
// Returns:
//
//	A Token with the specified type, position, and an empty Literal value.
//
// Note:
//
//	The `tokenTypes` map must contain an entry for the given tokenType;
//	otherwise, the Lexeme will be an empty string.
func CreateToken(tokenType TokenType, line int32, column int) Token {
	lexeme := tokenTypes[tokenType]
	return Token{
		TokenType: tokenType,
		Lexeme:    lexeme,
		Literal:   nil,
		Line:      line,
		Column:    column,
	}
}

// CreateLiteralToken constructs and returns a new Token instance that includes
// an explicit literal value and lexeme.
//
// Unlike CreateToken, which derives the lexeme from a predefined mapping,
// CreateLiteralToken requires both the lexeme (source text) and the literal
// (interpreted value) to be provided directly. This is useful for tokens where
// the value is computed or parsed during tokenization rather than looked up.
//
// Parameters:
//   - tokenType: The classification of the token being created.
//   - literal:   The parsed or computed value this token represents. Can be
//     any type (string, number, bool, etc.).
//   - lexeme:    The original source text that produced this token.
//   - line:      The 0-based source line number where the token begins.
//   - column:    The 0-based column number where the token begins.
//
// Returns:
//
//	A Token with the specified type, lexeme, literal, and position.
func CreateLiteralToken(tokenType TokenType, literal any, lexeme string, line int32, column int) Token {
	return Token{
		TokenType: tokenType,
		Lexeme:    lexeme,
		Literal:   literal,
		Line:      line,
		Column:    column,
	}
}

// String returns a human-readable representation of the Token,
// formatted to show its type and lexeme. he output is primarily intended for debugging
// and logging purposes
//
// Example:
//
//	tok := CreateLiteralToken(NUMBER, 123, "123", 3, 10)
//	fmt.Println(tok)
//	// Output: Token {Type: NUMBER, Value: "123"}
func (t Token) String() string {
	return fmt.Sprintf("Token {Type: %s, Value: %q}", t.TokenType, t.Lexeme)
}
