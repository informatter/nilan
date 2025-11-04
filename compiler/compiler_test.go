package compiler

import (
	"nilan/token"
	"strings"
	"testing"
)

func assertBytecodeEquals(t *testing.T, got Bytecode, want Bytecode) {

	if len(got.Instructions) != len(want.Instructions) {
		t.Errorf("computed instructions has a different length than the expected instructions - got: %d, want: %d", len(got.Instructions), len(want.Instructions))
	}

	for i, instruction := range got.Instructions {
		if instruction != want.Instructions[i] {
			t.Errorf("computed instruction does not equal expected instruction at index %d", i)
		}
	}
	for i, constant := range got.ConstantsPool {
		if constant != want.ConstantsPool[i] {
			t.Errorf("computed constant does not equal expected constant at index %d - want: %v, got: %v", i, want.ConstantsPool[i], constant)
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
		compiler := New(tt.tokens)
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
		compiler := New(tt.tokens)
		bytecode, err := compiler.Compile()
		if err != nil {
			t.Errorf("compilation error occurred: %s", err.Error())
		}
		assertBytecodeEquals(t, bytecode, tt.expectedBytecode)
	}
}

func TestDiassembleBytecode(t *testing.T) {

	tests := []struct {
		tokens   []token.Token
		expected string
	}{
		{
			tokens: []token.Token{
				token.CreateLiteralToken(token.INT, int64(1), "1", 0, 0),
				token.CreateToken(token.ADD, 0, 0),
				token.CreateLiteralToken(token.INT, int64(2), "2", 0, 0),
				token.CreateToken(token.MULT, 0, 0),
				token.CreateLiteralToken(token.INT, int64(4), "4", 0, 0),
				token.CreateToken(token.ADD, 0, 0),
				token.CreateLiteralToken(token.INT, int64(3), "3", 0, 0),
				token.CreateToken(token.EOF, 0, 0),
			},
			expected: `opcode: OP_CONSTANT, operand: 0, operand widths: 2 bytes, value: 1
opcode: OP_CONSTANT, operand: 1, operand widths: 2 bytes, value: 2
opcode: OP_CONSTANT, operand: 2, operand widths: 2 bytes, value: 4
opcode: OP_MULTIPLY, operand: None, operand widths: 0 bytes
opcode: OP_ADD, operand: None, operand widths: 0 bytes
opcode: OP_CONSTANT, operand: 3, operand widths: 2 bytes, value: 3
opcode: OP_ADD, operand: None, operand widths: 0 bytes
opcode: OP_END, operand: None, operand widths: 0 bytes`,
		},
	}

	for _, tt := range tests {

		compiler := New(tt.tokens)
		_, err := compiler.Compile()
		if err != nil {
			t.Errorf("compilation error occurred: %s", err.Error())
		}

		result, err := compiler.DiassembleBytecode(false,"")
		if err != nil {
			t.Errorf("bytecode diassembly error: %s", err.Error())
		}

		if strings.TrimSpace(result) != strings.TrimSpace(tt.expected) {
			t.Errorf("\n\nwant:\n%s\n\ngot:\n%s", tt.expected, result)
		}
	}

}
