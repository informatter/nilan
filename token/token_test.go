package token

import (
	"testing"
)

func TestCreateToken(t *testing.T) {
	tests := []struct {
		name      string
		tokenType TokenType
		value     string
		want      Token
	}{
		{
			name:      "Create ASSIGN token",
			tokenType: TokenType(ASSIGN),
			value:     "=",
			want:      Token{TokenType: TokenType(ASSIGN), Value: "="},
		},
		{
			name:      "Create IDENTIFIER token",
			tokenType: TokenType(IDENTIFIER),
			value:     "myVar",
			want:      Token{TokenType: TokenType(IDENTIFIER), Value: "myVar"},
		},
		{
			name:      "Create INT token",
			tokenType: TokenType(INT),
			value:     "42",
			want:      Token{TokenType: TokenType(INT), Value: "42"},
		},
		{
			name:      "Create MULT token",
			tokenType: TokenType(MULT),
			value:     "*",
			want:      Token{TokenType: TokenType(MULT), Value: "*"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CreateToken(tt.tokenType, tt.value)
			if got != tt.want {
				t.Errorf("createToken() = %v, want %v", got, tt.want)
			}
		})
	}
}
