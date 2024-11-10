package lexer

import (
	"nilan/token"
	"reflect"
	"testing"
)


func runTestSuccess(t *testing.T, scanner *Lexer, expected []token.Token){

	t.Run("ValidTokenScan", func(t *testing.T) {
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

func TestOperatorsSuccess(t *testing.T){
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
	runTestSuccess(t,scanner,expected)

}

func TestScanSuccess(t *testing.T) {
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
	runTestSuccess(t,scanner,expected)

}
