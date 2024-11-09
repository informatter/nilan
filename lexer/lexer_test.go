package lexer

import (
	"nilan/token"
	"reflect"
	"testing"
)

func TestScan(t *testing.T) {
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
