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

type Token struct {
	TokenType TokenType
	Value     string
	Line      int32
	Column    int
}

func CreateToken(tokenType TokenType, line int32, column int) Token {
	value := tokenTypes[tokenType]
	return Token{
		TokenType: tokenType,
		Value:     value,
		Line:      line,
		Column: column,
	}
}

func CreateLiteralToken(tokenType TokenType, value string, line int32, column int) Token {
	return Token{
		TokenType: tokenType,
		Value:     value,
		Line:      line,
		Column: column,

	}
}

// implements the fmt.Stringer interface to customize
// how the Token is rendered in the console.
func (t Token) String() string {
	return fmt.Sprintf("Token {Type: %s, Value: %q}", t.TokenType, t.Value)
}
