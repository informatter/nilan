package compiler

import (
	"nilan/ast"
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

func TestCompilerPrintStatement(t *testing.T) {
	tests := []struct {
		name             string
		statements       []ast.Stmt
		expectedBytecode Bytecode
	}{
		{
			name: "Print literal integer",
			statements: []ast.Stmt{
				ast.PrintStmt{
					Expression: ast.Grouping{
						Expression: ast.Literal{
							Value: int64(5),
						},
					},
				},
			},
			expectedBytecode: Bytecode{
				Instructions:  []byte{byte(OP_CONSTANT), 0, 0, byte(OP_PRINT), byte(OP_END)},
				ConstantsPool: []any{int64(5)},
			},
		},
		{
			name: "Print literal float",
			statements: []ast.Stmt{
				ast.PrintStmt{
					Expression: ast.Grouping{
						Expression: ast.Literal{
							Value: float64(5.545),
						},
					},
				},
			},
			expectedBytecode: Bytecode{
				Instructions:  []byte{byte(OP_CONSTANT), 0, 0, byte(OP_PRINT), byte(OP_END)},
				ConstantsPool: []any{float64(5.545)},
			},
		},
		{
			name: "Print expression result",
			statements: []ast.Stmt{
				ast.PrintStmt{
					Expression: ast.Grouping{
						Expression: ast.Binary{
							Left:     ast.Literal{Value: int64(2)},
							Operator: token.CreateToken(token.ADD, 0, 0),
							Right:    ast.Literal{Value: int64(10)},
						},
					},
				},
			},
			expectedBytecode: Bytecode{
				Instructions:  []byte{byte(OP_CONSTANT), 0, 0, byte(OP_CONSTANT), 0, 1, byte(OP_ADD), byte(OP_PRINT), byte(OP_END)},
				ConstantsPool: []any{int64(2), int64(10)},
			},
		},
	}
	for _, tt := range tests {
		compiler := NewASTCompiler()
		bytecode, err := compiler.CompileAST(tt.statements)
		if err != nil {
			t.Errorf("compilation error occurred: %s", err.Error())
		}
		assertBytecodeEquals(t, bytecode, tt.expectedBytecode)
	}
}

func TestCompileNumericTokens_BinaryExpressions(t *testing.T) {
	tests := []struct {
		statements       []ast.Stmt
		expectedBytecode Bytecode
	}{
		{
			statements: []ast.Stmt{
				ast.ExpressionStmt{
					Expression: ast.Binary{
						Left:     ast.Literal{Value: int64(5)},
						Operator: token.CreateToken(token.ADD, 0, 0),
						Right:    ast.Literal{Value: int64(1)},
					},
				},
			},
			expectedBytecode: Bytecode{
				Instructions:  []byte{byte(OP_CONSTANT), 0, 0, byte(OP_CONSTANT), 0, 1, byte(OP_ADD), byte(OP_END)},
				ConstantsPool: []any{int64(5), int64(1)},
			},
		},
		{
			statements: []ast.Stmt{
				ast.ExpressionStmt{
					Expression: ast.Binary{
						Left:     ast.Literal{Value: int64(5)},
						Operator: token.CreateToken(token.MULT, 0, 0),
						Right:    ast.Literal{Value: int64(1)},
					},
				},
			},
			expectedBytecode: Bytecode{
				Instructions:  []byte{byte(OP_CONSTANT), 0, 0, byte(OP_CONSTANT), 0, 1, byte(OP_MULTIPLY), byte(OP_END)},
				ConstantsPool: []any{int64(5), int64(1)},
			},
		},
		{
			statements: []ast.Stmt{
				ast.ExpressionStmt{
					Expression: ast.Binary{
						Left:     ast.Literal{Value: int64(5)},
						Operator: token.CreateToken(token.DIV, 0, 0),
						Right:    ast.Literal{Value: int64(1)},
					},
				},
			},
			expectedBytecode: Bytecode{
				Instructions:  []byte{byte(OP_CONSTANT), 0, 0, byte(OP_CONSTANT), 0, 1, byte(OP_DIVIDE), byte(OP_END)},
				ConstantsPool: []any{int64(5), int64(1)},
			},
		},
		{
			statements: []ast.Stmt{
				ast.ExpressionStmt{
					Expression: ast.Binary{
						Left:     ast.Literal{Value: int64(5)},
						Operator: token.CreateToken(token.SUB, 0, 0),
						Right:    ast.Literal{Value: int64(1)},
					},
				},
			},
			expectedBytecode: Bytecode{
				Instructions:  []byte{byte(OP_CONSTANT), 0, 0, byte(OP_CONSTANT), 0, 1, byte(OP_SUBTRACT), byte(OP_END)},
				ConstantsPool: []any{int64(5), int64(1)},
			},
		},
	}

	for _, tt := range tests {
		compiler := NewASTCompiler()
		bytecode, err := compiler.CompileAST(tt.statements)
		if err != nil {
			t.Errorf("compilation error occurred: %s", err.Error())
		}
		assertBytecodeEquals(t, bytecode, tt.expectedBytecode)
	}
}

func TestCompileNumericTokens_UnaryExpressions(t *testing.T) {
	tests := []struct {
		statements       []ast.Stmt
		expectedBytecode Bytecode
	}{
		{
			statements: []ast.Stmt{
				ast.ExpressionStmt{
					Expression: ast.Unary{
						Operator: token.CreateToken(token.SUB, 0, 0),
						Right:    ast.Literal{Value: int64(5)},
					},
				},
			},
			expectedBytecode: Bytecode{
				Instructions:  []byte{byte(OP_CONSTANT), 0, 0, byte(OP_NEGATE), byte(OP_END)},
				ConstantsPool: []any{int64(5)},
			},
		},
	}

	for _, tt := range tests {
		compiler := NewASTCompiler()
		bytecode, err := compiler.CompileAST(tt.statements)
		if err != nil {
			t.Errorf("compilation error occurred: %s", err.Error())
		}
		assertBytecodeEquals(t, bytecode, tt.expectedBytecode)
	}
}

func TestDiassembleBytecode(t *testing.T) {

	tests := []struct {
		statements []ast.Stmt
		expected   string
	}{
		{
			statements: []ast.Stmt{
				ast.ExpressionStmt{
					Expression: ast.Binary{
						Left: ast.Binary{
							Left:     ast.Literal{Value: int64(1)},
							Operator: token.CreateToken(token.ADD, 0, 0),
							Right: ast.Binary{
								Left:     ast.Literal{Value: int64(2)},
								Operator: token.CreateToken(token.MULT, 0, 0),
								Right:    ast.Literal{Value: int64(4)},
							},
						},
						Operator: token.CreateToken(token.ADD, 0, 0),
						Right:    ast.Literal{Value: int64(3)},
					},
				},
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

		compiler := NewASTCompiler()
		_, err := compiler.CompileAST(tt.statements)
		if err != nil {
			t.Errorf("compilation error occurred: %s", err.Error())
		}

		result, err := compiler.DiassembleBytecode(false, "")
		if err != nil {
			t.Errorf("bytecode diassembly error: %s", err.Error())
		}

		if strings.TrimSpace(result) != strings.TrimSpace(tt.expected) {
			t.Errorf("\n\nwant:\n%s\n\ngot:\n%s", tt.expected, result)
		}
	}

}

// TestASTCompileArithmetic tests the AST-to-bytecode compiler with arithmetic expressions
func TestASTCompileArithmetic(t *testing.T) {
	tests := []struct {
		name             string
		statements       []ast.Stmt
		expectedBytecode Bytecode
	}{
		{
			name: "Binary Addition",
			statements: []ast.Stmt{
				ast.ExpressionStmt{
					Expression: ast.Binary{
						Left:     ast.Literal{Value: int64(5)},
						Operator: token.CreateToken(token.ADD, 0, 0),
						Right:    ast.Literal{Value: int64(1)},
					},
				},
			},
			expectedBytecode: Bytecode{
				Instructions:  []byte{byte(OP_CONSTANT), 0, 0, byte(OP_CONSTANT), 0, 1, byte(OP_ADD), byte(OP_END)},
				ConstantsPool: []any{int64(5), int64(1)},
			},
		},
		{
			name: "Binary Multiplication",
			statements: []ast.Stmt{
				ast.ExpressionStmt{
					Expression: ast.Binary{
						Left:     ast.Literal{Value: int64(5)},
						Operator: token.CreateToken(token.MULT, 0, 0),
						Right:    ast.Literal{Value: int64(3)},
					},
				},
			},
			expectedBytecode: Bytecode{
				Instructions:  []byte{byte(OP_CONSTANT), 0, 0, byte(OP_CONSTANT), 0, 1, byte(OP_MULTIPLY), byte(OP_END)},
				ConstantsPool: []any{int64(5), int64(3)},
			},
		},
		{
			name: "Unary Negation",
			statements: []ast.Stmt{
				ast.ExpressionStmt{
					Expression: ast.Unary{
						Operator: token.CreateToken(token.SUB, 0, 0),
						Right:    ast.Literal{Value: int64(5)},
					},
				},
			},
			expectedBytecode: Bytecode{
				Instructions:  []byte{byte(OP_CONSTANT), 0, 0, byte(OP_NEGATE), byte(OP_END)},
				ConstantsPool: []any{int64(5)},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			compiler := NewASTCompiler()
			bytecode, err := compiler.CompileAST(tt.statements)
			if err != nil {
				t.Errorf("compilation error occurred: %s", err.Error())
			}
			assertBytecodeEquals(t, bytecode, tt.expectedBytecode)
		})
	}
}

// TestASTCompilerDisassembleBytecode tests the ASTCompiler's DisassembleBytecode method
func TestASTCompilerDisassembleBytecode(t *testing.T) {
	tests := []struct {
		name       string
		statements []ast.Stmt
		expected   string
	}{
		{
			name: "Simple Addition",
			statements: []ast.Stmt{
				ast.ExpressionStmt{
					Expression: ast.Binary{
						Left:     ast.Literal{Value: int64(5)},
						Operator: token.CreateToken(token.ADD, 0, 0),
						Right:    ast.Literal{Value: int64(3)},
					},
				},
			},
			expected: `opcode: OP_CONSTANT, operand: 0, operand widths: 2 bytes, value: 5
opcode: OP_CONSTANT, operand: 1, operand widths: 2 bytes, value: 3
opcode: OP_ADD, operand: None, operand widths: 0 bytes
opcode: OP_END, operand: None, operand widths: 0 bytes`,
		},
		{
			name: "Complex Expression",
			statements: []ast.Stmt{
				ast.ExpressionStmt{
					Expression: ast.Binary{
						Left: ast.Binary{
							Left:     ast.Literal{Value: int64(1)},
							Operator: token.CreateToken(token.ADD, 0, 0),
							Right:    ast.Literal{Value: int64(2)},
						},
						Operator: token.CreateToken(token.MULT, 0, 0),
						Right: ast.Binary{
							Left:     ast.Literal{Value: int64(4)},
							Operator: token.CreateToken(token.ADD, 0, 0),
							Right:    ast.Literal{Value: int64(3)},
						},
					},
				},
			},
			expected: `opcode: OP_CONSTANT, operand: 0, operand widths: 2 bytes, value: 1
opcode: OP_CONSTANT, operand: 1, operand widths: 2 bytes, value: 2
opcode: OP_ADD, operand: None, operand widths: 0 bytes
opcode: OP_CONSTANT, operand: 2, operand widths: 2 bytes, value: 4
opcode: OP_CONSTANT, operand: 3, operand widths: 2 bytes, value: 3
opcode: OP_ADD, operand: None, operand widths: 0 bytes
opcode: OP_MULTIPLY, operand: None, operand widths: 0 bytes
opcode: OP_END, operand: None, operand widths: 0 bytes`,
		},
		{
			name: "Division and Subtraction",
			statements: []ast.Stmt{
				ast.ExpressionStmt{
					Expression: ast.Binary{
						Left: ast.Binary{
							Left:     ast.Literal{Value: int64(10)},
							Operator: token.CreateToken(token.DIV, 0, 0),
							Right:    ast.Literal{Value: int64(2)},
						},
						Operator: token.CreateToken(token.SUB, 0, 0),
						Right:    ast.Literal{Value: int64(1)},
					},
				},
			},
			expected: `opcode: OP_CONSTANT, operand: 0, operand widths: 2 bytes, value: 10
opcode: OP_CONSTANT, operand: 1, operand widths: 2 bytes, value: 2
opcode: OP_DIVIDE, operand: None, operand widths: 0 bytes
opcode: OP_CONSTANT, operand: 2, operand widths: 2 bytes, value: 1
opcode: OP_SUBTRACT, operand: None, operand widths: 0 bytes
opcode: OP_END, operand: None, operand widths: 0 bytes`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			compiler := NewASTCompiler()
			_, err := compiler.CompileAST(tt.statements)
			if err != nil {
				t.Errorf("compilation error occurred: %s", err.Error())
			}

			result, err := compiler.DiassembleBytecode(false, "")
			if err != nil {
				t.Errorf("bytecode disassembly error: %s", err.Error())
			}

			if strings.TrimSpace(result) != strings.TrimSpace(tt.expected) {
				t.Errorf("\n\nwant:\n%s\n\ngot:\n%s", tt.expected, result)
			}
		})
	}
}
