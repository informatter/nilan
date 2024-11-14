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

	EOF = "EOF"

	// keywords
	FUNC   = "fn"
	CONST  = "const"
	VAR    = "var"
	RETURN = "return"
	IF     = "if"
	ELSE   = "else"
	ELIF   = "elif"
	BREAK  = "break"
	TRUE   = "true"
	FALSE  = "false"
	NULL   = "null"
)

var KeyWords = map[string]TokenType{
	"fn":     FUNC,
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
}

func CreateToken(tokenType TokenType) Token {
	value := tokenTypes[tokenType]
	return Token{
		TokenType: tokenType,
		Value:     value,
	}
}

func CreateLiteralToken(tokenType TokenType, value string) Token {
	return Token{
		TokenType: tokenType,
		Value:     value,
	}
}

// implements the fmt.Stringer interface to customize
// how the Token is rendered in the console.
func (t Token) String() string {
	return fmt.Sprintf("Token {Type: %s, Value: %q}", t.TokenType, t.Value)
}
