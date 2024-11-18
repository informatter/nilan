package lexer

import (

	"nilan/token"
	"reflect"
	"testing"
)

func runTest(t *testing.T, testName string, scanner *Lexer, expected []token.Token) {

	t.Run(testName, func(t *testing.T) {
		got, err := scanner.Scan()
		if err != nil {
			t.Errorf("scanner.Scan() raised an error: %v", err)
		}

		// Use reflect.DeepEqual to compare the slices
		if !reflect.DeepEqual(got, expected) {
			t.Errorf("scanner.Scan() = %v, want %v", got, expected)
		}
	})
}

func TestScanLoose(t *testing.T) {
	testName := "TestScanLoose"
	expected := []token.Token{
		token.CreateToken(token.LPA),
		token.CreateToken(token.RPA),
		token.CreateToken(token.LCUR),
		token.CreateToken(token.RCUR),
		token.CreateToken(token.MULT),
		token.CreateToken(token.MULT),
		token.CreateToken(token.SEMICOLON),
		token.CreateToken(token.ADD),
		token.CreateToken(token.BANG),
		token.CreateToken(token.ASSIGN),
		token.CreateToken(token.LESS),
		token.CreateToken(token.EQUAL_EQUAL),
		token.CreateToken(token.EQUAL_EQUAL),
		token.CreateToken(token.EQUAL_EQUAL),
		token.CreateToken(token.EQUAL_EQUAL),
		token.CreateToken(token.NOT_EQUAL),
		token.CreateToken(token.LESS_EQUAL),
		token.CreateToken(token.LESS_EQUAL),
		token.CreateToken(token.LARGER_EQUAL),
		token.CreateToken(token.LARGER_EQUAL),
		token.CreateToken(token.EQUAL_EQUAL),
		token.CreateToken(token.NOT_EQUAL),
		token.CreateLiteralToken(token.IDENTIFIER, "my_var"),
		token.CreateToken(token.ASSIGN),
		token.CreateToken(token.LCUR),
		token.CreateToken(token.RCUR),
		token.CreateToken(token.EOF),
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
	# my_var = {}
	# some comment # # # # 
	my_var = {
	}
	`
	scanner := CreateLexer(test)
	runTest(t, testName, scanner, expected)

}

func TestLiteralStrings(t *testing.T){
	testName := "TestLiteralStrings"

	
	expected := []token.Token{
		token.CreateLiteralToken(token.VAR, "var"),
		token.CreateLiteralToken(token.IDENTIFIER, "myString"),
		token.CreateToken(token.ASSIGN),
		token.CreateLiteralToken(token.STRING, "hellow"),
		token.CreateLiteralToken(token.STRING, "hi"),
		token.CreateLiteralToken(token.VAR, "var"),
		token.CreateLiteralToken(token.IDENTIFIER, "tabedString"),
		token.CreateToken(token.ASSIGN),
		token.CreateLiteralToken(token.STRING, "tabed		"),
		token.CreateToken(token.EOF),
	}
	test := `
	var myString = "hellow" "hi"
	var tabedString = "tabed		"
	`
	scanner := CreateLexer(test)
	runTest(t, testName, scanner, expected)
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
            errMsg:  "unclosed string literal: unclosed\nline: 0",
        },
        {
            name:    "String literal with newline",
            input:   "\"new\nline\"",
            wantErr: true,
            errMsg:  "unclosed string literal: new\nline: 0",
        },
        {
            name:    "Only opening quote",
            input:   `"`,
            wantErr: true,
            errMsg:  "unclosed string literal: \nline: 0",
        },
        {
            name:    "String literal at end of input",
            input:   `hello "world`,
            wantErr: true,
            errMsg:  "unclosed string literal: world\nline: 0",
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

func TestScanSourceCode(t *testing.T) {
	testName := "TestScanSourceCode"
	expected := []token.Token{
		token.CreateLiteralToken(token.FUNC, "fn"),
		token.CreateLiteralToken(token.IDENTIFIER, "myFunction"),
		token.CreateToken(token.LPA),
		token.CreateLiteralToken(token.IDENTIFIER, "a"),
		token.CreateToken(token.COMMA),
		token.CreateLiteralToken(token.IDENTIFIER, "b"),
		token.CreateToken(token.RPA),
		token.CreateToken(token.LCUR),
		token.CreateLiteralToken(token.RETURN, "return"),
		token.CreateLiteralToken(token.IDENTIFIER, "a"),
		token.CreateToken(token.ADD),
		token.CreateLiteralToken(token.IDENTIFIER, "b"),
		token.CreateToken(token.RCUR),
		token.CreateLiteralToken(token.VAR, "var"),
		token.CreateLiteralToken(token.IDENTIFIER, "_foo_bar"),
		token.CreateToken(token.ASSIGN),
		token.CreateLiteralToken(token.FLOAT, "0.000001"),
		token.CreateLiteralToken(token.VAR, "var"),
		token.CreateLiteralToken(token.IDENTIFIER, "myInt"),
		token.CreateToken(token.ASSIGN),
		token.CreateLiteralToken(token.INT, "123"),
		token.CreateLiteralToken(token.VAR, "var"),
		token.CreateLiteralToken(token.IDENTIFIER, "myNegativeInt"),
		token.CreateToken(token.ASSIGN),
		token.CreateLiteralToken(token.INT, "-123"),
		token.CreateLiteralToken(token.VAR, "var"),
		token.CreateLiteralToken(token.IDENTIFIER, "myNegativeFloat"),
		token.CreateToken(token.ASSIGN),
		token.CreateLiteralToken(token.FLOAT, "-0.01"),
		token.CreateLiteralToken(token.VAR, "var"),
		token.CreateLiteralToken(token.IDENTIFIER, "myString"),
		token.CreateToken(token.ASSIGN),
		token.CreateLiteralToken(token.STRING, "hellow"),
		token.CreateToken(token.EOF),
	}

	test := `
	fn myFunction(a, b){
		return a + b
	}
	var _foo_bar = 0.000001
	var myInt = 123
	var myNegativeInt = -123
	var myNegativeFloat = -0.01
	var myString = "hellow"
	`

	scanner := CreateLexer(test)
	runTest(t, testName, scanner, expected)

}
