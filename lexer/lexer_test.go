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

func TestScanSuccess(t *testing.T) {
	testName := "TestScanSuccess"
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
