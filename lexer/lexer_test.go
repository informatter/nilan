package lexer

import (
	"nilan/token"
	"testing"
)

func runTest(expected []token.Token, scanner *Lexer, t *testing.T) {

	result, err := scanner.Scan()
	if err != nil {
		t.Errorf("scanner.Scan() raised an error: %v", err)
	}

	for i, tt := range expected {
		tok := result[i]

		if tok.TokenType != tt.TokenType {
			t.Fatalf("Wrong token type - expected: %q, got: %q", tt.TokenType, tok.TokenType)
		}
		if tok.Lexeme != tt.Lexeme {
			t.Fatalf("Wrong lexeme - expected: %q, got: %q", tt.Lexeme, tok.Lexeme)
		}

		if tok.Literal != tt.Literal {

			t.Fatalf("Wrong literal - expected: %q, got: %q", tt.Lexeme, tok.Lexeme)
		}

	}
}

func TestComments(t *testing.T) {

	expected := []token.Token{
		token.CreateToken(token.LPA, 0, 0),
		token.CreateToken(token.RPA, 0, 0),
		token.CreateToken(token.LCUR, 0, 0),
		token.CreateToken(token.RCUR, 0, 0),
		token.CreateToken(token.MULT, 0, 0),
		token.CreateToken(token.MULT, 0, 0),
		token.CreateToken(token.SEMICOLON, 0, 0),
		token.CreateToken(token.ADD, 0, 0),
		token.CreateToken(token.BANG, 0, 0),
		token.CreateToken(token.ASSIGN, 0, 0),
		token.CreateToken(token.LESS, 0, 0),
		token.CreateToken(token.EQUAL_EQUAL, 0, 0),
		token.CreateToken(token.EQUAL_EQUAL, 0, 0),
		token.CreateToken(token.EQUAL_EQUAL, 0, 0),
		token.CreateToken(token.EQUAL_EQUAL, 0, 0),
		token.CreateToken(token.NOT_EQUAL, 0, 0),
		token.CreateToken(token.LESS_EQUAL, 0, 0),
		token.CreateToken(token.LESS_EQUAL, 0, 0),
		token.CreateToken(token.LARGER_EQUAL, 0, 0),
		token.CreateToken(token.LARGER_EQUAL, 0, 0),
		token.CreateToken(token.EQUAL_EQUAL, 0, 0),
		token.CreateToken(token.NOT_EQUAL, 0, 0),
		{
			TokenType: token.IDENTIFIER, Lexeme: "my_var",
		},
		token.CreateToken(token.ASSIGN, 0, 0),
		token.CreateToken(token.LCUR, 0, 0),
		token.CreateToken(token.RCUR, 0, 0),
		token.CreateToken(token.EOF, 0, 0),
	}

	test := `
	(){}
	**;+
	!
	=
	<
	== == ==
	==
	!=
	<= <=
	>=
	>===!=
	my_var = {}
	# some comment # # # # 
	#my_var = {
	#}
	`

	scanner := CreateLexer(test)
	runTest(expected, scanner, t)

}

func TestLiteralStrings(t *testing.T) {

	multiLine := `
	 this is a multi line comment
	 which continues here
	`
	expected := []token.Token{
		{
			TokenType: token.VAR, Lexeme: "var",
		},
		{
			TokenType: token.IDENTIFIER, Lexeme: "myString",
		},
		token.CreateToken(token.ASSIGN, 0, 0),
		token.CreateLiteralToken(token.STRING, "hellow", "hellow", 0, 0),

		token.CreateLiteralToken(token.STRING, "hi", "hi", 0, 0),
		{
			TokenType: token.VAR, Lexeme: "var",
		},
		{
			TokenType: token.IDENTIFIER, Lexeme: "tabedString",
		},
		token.CreateToken(token.ASSIGN, 0, 0),
		token.CreateLiteralToken(token.STRING, "tabed	", "tabed	", 0, 0),
		token.CreateLiteralToken(token.STRING, multiLine, multiLine, 0, 0),
		token.CreateToken(token.EOF, 0, 0),
	}
	test := `
	var myString = "hellow" "hi"
	var tabedString = "tabed	"
	"
	 this is a multi line comment
	 which continues here
	"
	`
	scanner := CreateLexer(test)
	runTest(expected, scanner, t)
}

func TestHandleStringLiteralErrors(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
		errMsg  string
	}{
		{
			name:    "Unclosed string literal",
			input:   `var c ="unclosed`,
			wantErr: true,
			errMsg:  "unclosed string literal: 'unclosed', line: 0",
		},
		{
			name:    "Only opening quote",
			input:   `"`,
			wantErr: true,
			errMsg:  "unclosed string literal: '', line: 0",
		},
		{
			name:    "String literal at end of input",
			input:   `hello "world`,
			wantErr: true,
			errMsg:  "unclosed string literal: 'world', line: 0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			scanner := CreateLexer(tt.input)

			_, err := scanner.Scan()

			if err == nil {
				t.Errorf("handleStringLiteral() error = nil, wantErr %v", tt.wantErr)
				return
			}
			if err.Error() != tt.errMsg {
				t.Errorf("handleStringLiteral() error = %v, wantErr %v", err, tt.errMsg)
			}

		})
	}
}

func TestHandleNumberErrors(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
		errMsg  string
	}{
		{
			name:    "Malformed decimal number A",
			input:   `1.11.`,
			wantErr: true,
			errMsg:  "invalid number: '1.11.', line: 0",
		},
		{
			name:    "Malformed decimal number A",
			input:   `0.000.111`,
			wantErr: true,
			errMsg:  "invalid number: '0.000.111', line: 0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			scanner := CreateLexer(tt.input)

			_, err := scanner.Scan()

			if err == nil {
				t.Errorf("handleNumber() error = nil, wantErr %v", tt.wantErr)
				return
			}
			if err.Error() != tt.errMsg {
				t.Errorf("handleNumber() error = %v, wantErr %v", err, tt.errMsg)
			}

		})
	}
}

func TestHandleNumber(t *testing.T) {
	expected := []token.Token{
		token.CreateLiteralToken(token.FLOAT, float64(0.2), ".2", 0, 0),
		token.CreateLiteralToken(token.FLOAT, float64(0.0001), "0.0001", 0, 0),
		token.CreateLiteralToken(token.INT, int64(1000), "1000", 0, 0),
	}
	test := `
	.2
	0.0001
	1000
	`
	scanner := CreateLexer(test)
	runTest(expected, scanner, t)
}

func TestScanSourceCode(t *testing.T) {

	expected := []token.Token{

		{
			TokenType: token.FUNC, Lexeme: "fn",
		},
		{
			TokenType: token.IDENTIFIER, Lexeme: "myFunction",
		},
		token.CreateToken(token.LPA, 0, 0),
		{
			TokenType: token.IDENTIFIER, Lexeme: "a",
		},
		token.CreateToken(token.COMMA, 0, 0),
		{
			TokenType: token.IDENTIFIER, Lexeme: "b",
		},
		token.CreateToken(token.RPA, 0, 0),
		token.CreateToken(token.LCUR, 0, 0),
		{
			TokenType: token.RETURN, Lexeme: "return",
		},
		{
			TokenType: token.IDENTIFIER, Lexeme: "a",
		},
		token.CreateToken(token.ADD, 0, 0),
		{
			TokenType: token.IDENTIFIER, Lexeme: "b",
		},
		token.CreateToken(token.RCUR, 0, 0),

		{
			TokenType: token.VAR, Lexeme: "var",
		},
		{
			TokenType: token.IDENTIFIER, Lexeme: "result",
		},
		token.CreateToken(token.ASSIGN, 0, 0),
		{
			TokenType: token.IDENTIFIER, Lexeme: "myFunction",
		},
		token.CreateToken(token.LPA, 0, 0),
		token.CreateLiteralToken(token.INT, int64(2), "2", 0, 0),
		token.CreateToken(token.ADD, 0, 0),
		token.CreateLiteralToken(token.INT, int64(5), "5", 0, 0),
		token.CreateToken(token.RPA, 0, 0),
		{
			TokenType: token.VAR, Lexeme: "var",
		},
		{
			TokenType: token.IDENTIFIER, Lexeme: "_foo_bar",
		},
		token.CreateToken(token.ASSIGN, 0, 0),
		token.CreateLiteralToken(token.FLOAT, 0.000001, "0.000001", 0, 0),
		{
			TokenType: token.VAR, Lexeme: "var",
		},
		{
			TokenType: token.IDENTIFIER, Lexeme: "myInt",
		},
		token.CreateToken(token.ASSIGN, 0, 0),
		token.CreateLiteralToken(token.INT, int64(123), "123", 0, 0),
		{
			TokenType: token.VAR, Lexeme: "var",
		},
		{
			TokenType: token.IDENTIFIER, Lexeme: "myNegativeInt",
		},
		token.CreateToken(token.ASSIGN, 0, 0),
		{
			TokenType: token.SUB, Lexeme: "-",
		},
		token.CreateLiteralToken(token.INT, int64(123), "123", 0, 0),
		{
			TokenType: token.VAR, Lexeme: "var",
		},
		{
			TokenType: token.IDENTIFIER, Lexeme: "myNegativeFloat",
		},
		token.CreateToken(token.ASSIGN, 0, 0),
		{
			TokenType: token.SUB, Lexeme: "-",
		},
		token.CreateLiteralToken(token.FLOAT, 0.01, "0.01", 0, 0),
		{
			TokenType: token.VAR, Lexeme: "var",
		},
		{
			TokenType: token.IDENTIFIER, Lexeme: "myString",
		},
		token.CreateToken(token.ASSIGN, 0, 0),
		token.CreateLiteralToken(token.STRING, "hellow", "hellow", 0, 0),

		{
			TokenType: token.IF, Lexeme: "if",
		},
		{
			TokenType: token.AND, Lexeme: "and",
		},
		{
			TokenType: token.OR, Lexeme: "or",
		},
		{
			TokenType: token.WHILE, Lexeme: "while",
		},
		{
			TokenType: token.FOR, Lexeme: "for",
		},
		token.CreateToken(token.EOF, 0, 0),
	}
	test := `
	fn myFunction(a, b){
		return a + b
	}
	var result = myFunction(2+5)
	var _foo_bar = 0.000001
	var myInt = 123
	var myNegativeInt = -123
	var myNegativeFloat = -0.01
	var myString = "hellow"
	
	if and or while for
	`

	scanner := CreateLexer(test)
	runTest(expected, scanner, t)

}
