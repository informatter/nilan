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

	// keywords reserved for nilan - i.e fn, const, return, if, else, break ...
	KEYWORD = "KEYWORD"

	INT = "INT"

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
)

var KeyWords = map[string]string{
	"fun":    "fun",
	"const":  "const",
	"return": "return",
	"if":     "if",
	"else":   "else",
	"elif":   "elif",
	"break":  "break",
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

// implements the fmt.Stringer interface to customize
// how the Token is rendered in the console.
func (t Token) String() string {
	return fmt.Sprintf("Token {Type: %s, Value: %q}", t.TokenType, t.Value)
}
