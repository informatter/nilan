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

func TestComment(t *testing.T) {
	testName := "TestComment"
	expected := []token.Token{
		token.CreateToken(token.EOF),
	}
	test :=
		`
	# fmt.Println()
	# comment
	`
	scanner := CreateLexer(test)
	runTest(t, testName, scanner, expected)
}

func TestOperatorsSuccess(t *testing.T) {
	testName := "TestOperatorsSuccess"
	expected := []token.Token{
		token.CreateToken(token.EQUAL_EQUAL),
		token.CreateToken(token.DIV),
		token.CreateToken(token.ASSIGN),
		token.CreateToken(token.MULT),
		token.CreateToken(token.ADD),
		token.CreateToken(token.LARGER),
		token.CreateToken(token.SUB),
		token.CreateToken(token.LESS),
		token.CreateToken(token.NOT_EQUAL),
		token.CreateToken(token.LESS_EQUAL),
		token.CreateToken(token.LARGER_EQUAL),
		token.CreateToken(token.BANG),
		token.CreateToken(token.BANG),
		token.CreateToken(token.EOF),
	}
	scanner := CreateLexer("==/=*+>-<!=<=>=!!")
	runTest(t, testName, scanner, expected)

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
		token.CreateToken(token.NOT_EQUAL),
		token.CreateToken(token.LESS_EQUAL),
		token.CreateToken(token.EOF),
	}

	scanner := CreateLexer("(){}**;+!=<=")
	runTest(t, testName, scanner, expected)

}
