package compiler

import (
	"nilan/token"
	"testing"
)

func TestCompileNumericTokens(t *testing.T) {

	tests := []struct {
		tokens           []token.Token
		expectedBytecode Bytecode
	}{
		{
			tokens: []token.Token{
				token.CreateLiteralToken(token.INT, int64(5), "5", 0, 0),
				token.CreateLiteralToken(token.INT, int64(1), "1", 0, 0),
				token.CreateToken(token.EOF, 0, 0),
			},
			expectedBytecode: Bytecode{
				// NOTE: expected instruction operands are the indices in the constants pool
				// where 5 and 1 where added, encoded in Big Endian. In this case 5 is added
				// to index 0 and 1 to index 1. 0 in Big Endian uint16 is represented as 0,0 (decimal)
				// and 1 as 0,1 (decimal)
				Instructions:  []byte{byte(OP_CONSTANT), 0, 0, byte(OP_CONSTANT), 0, 1, byte(OP_END)},
				ConstantsPool: []any{int64(5), int64(1)},
			},
		},
	}

	for _, tt := range tests {

		compiler := NewCompiler()
		bytecode, err := compiler.Compile(tt.tokens)
		if err != nil {
			t.Errorf("compilation error occurred: %s", err.Error())
		}

		for i, instruction := range bytecode.Instructions {
			expectedInstruction := tt.expectedBytecode.Instructions[i]
			if expectedInstruction != instruction {
				t.Errorf("expected instruction does not equal computed instruction")
			}
		}

		for i, constant := range bytecode.ConstantsPool {
			expectedConstant := tt.expectedBytecode.ConstantsPool[i]
			if constant != expectedConstant {
				t.Errorf("expected constant does not equal computed constant - want: %v, got: %v", expectedConstant, constant)
			}
		}

	}
}
