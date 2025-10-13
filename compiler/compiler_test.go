package compiler

import (
	"nilan/token"
	"testing"
)

func assertBytecodeEquals(t *testing.T, got Bytecode, want Bytecode) {
	for i, instruction := range got.Instructions {
		if instruction != want.Instructions[i] {
			t.Errorf("expected instruction does not equal computed instruction at index %d", i)
		}
	}
	for i, constant := range got.ConstantsPool {
		if constant != want.ConstantsPool[i] {
			t.Errorf("expected constant does not equal computed constant at index %d - want: %v, got: %v", i, want.ConstantsPool[i], constant)
		}
	}
}

func TestCompileNumericTokens_BinaryExpressions(t *testing.T) {
	tests := []struct {
		tokens           []token.Token
		expectedBytecode Bytecode
	}{
		{
			tokens: []token.Token{
				token.CreateLiteralToken(token.INT, int64(5), "5", 0, 0),
				token.CreateToken(token.ADD, 0, 0),
				token.CreateLiteralToken(token.INT, int64(1), "1", 0, 0),
				token.CreateToken(token.EOF, 0, 0),
			},
			expectedBytecode: Bytecode{
				Instructions:  []byte{byte(OP_CONSTANT), 0, 0, byte(OP_CONSTANT), 0, 1, byte(OP_ADD), byte(OP_END)},
				ConstantsPool: []any{int64(5), int64(1)},
			},
		},
		{
			tokens: []token.Token{
				token.CreateLiteralToken(token.INT, int64(5), "5", 0, 0),
				token.CreateToken(token.MULT, 0, 0),
				token.CreateLiteralToken(token.INT, int64(1), "1", 0, 0),
				token.CreateToken(token.EOF, 0, 0),
			},
			expectedBytecode: Bytecode{
				Instructions:  []byte{byte(OP_CONSTANT), 0, 0, byte(OP_CONSTANT), 0, 1, byte(OP_MULTIPLY), byte(OP_END)},
				ConstantsPool: []any{int64(5), int64(1)},
			},
		},
		{
			tokens: []token.Token{
				token.CreateLiteralToken(token.INT, int64(5), "5", 0, 0),
				token.CreateToken(token.DIV, 0, 0),
				token.CreateLiteralToken(token.INT, int64(1), "1", 0, 0),
				token.CreateToken(token.EOF, 0, 0),
			},
			expectedBytecode: Bytecode{
				Instructions:  []byte{byte(OP_CONSTANT), 0, 0, byte(OP_CONSTANT), 0, 1, byte(OP_DIVIDE), byte(OP_END)},
				ConstantsPool: []any{int64(5), int64(1)},
			},
		},
		{
			tokens: []token.Token{
				token.CreateLiteralToken(token.INT, int64(5), "5", 0, 0),
				token.CreateToken(token.SUB, 0, 0),
				token.CreateLiteralToken(token.INT, int64(1), "1", 0, 0),
				token.CreateToken(token.EOF, 0, 0),
			},
			expectedBytecode: Bytecode{
				Instructions:  []byte{byte(OP_CONSTANT), 0, 0, byte(OP_CONSTANT), 0, 1, byte(OP_SUBTRACT), byte(OP_END)},
				ConstantsPool: []any{int64(5), int64(1)},
			},
		},
	}

	for _, tt := range tests {
		compiler := NewCompiler(tt.tokens)
		bytecode, err := compiler.Compile()
		if err != nil {
			t.Errorf("compilation error occurred: %s", err.Error())
		}
		assertBytecodeEquals(t, bytecode, tt.expectedBytecode)
	}
}

func TestCompileNumericTokens_UnaryExpressions(t *testing.T) {
	tests := []struct {
		tokens           []token.Token
		expectedBytecode Bytecode
	}{
		{
			tokens: []token.Token{
				token.CreateToken(token.SUB, 0, 0),
				token.CreateLiteralToken(token.INT, int64(5), "5", 0, 0),
				token.CreateToken(token.EOF, 0, 0),
			},
			expectedBytecode: Bytecode{
				Instructions:  []byte{byte(OP_CONSTANT), 0, 0, byte(OP_NEGATE), byte(OP_END)},
				ConstantsPool: []any{int64(5)},
			},
		},
	}

	for _, tt := range tests {
		compiler := NewCompiler(tt.tokens)
		bytecode, err := compiler.Compile()
		if err != nil {
			t.Errorf("compilation error occurred: %s", err.Error())
		}
		assertBytecodeEquals(t, bytecode, tt.expectedBytecode)
	}
}
