package compiler

import (
	"nilan/ast"
	"nilan/lexer"
	"nilan/parser"
	"nilan/token"
	"testing"
)

// TestFullPipeline demonstrates the complete pipeline: tokens -> AST -> bytecode
// This test shows that the AST-to-bytecode compiler can successfully compile
// arithmetic expressions
func TestFullPipeline(t *testing.T) {
	tests := []struct {
		name             string
		source           string
		expectedBytecode Bytecode
	}{
		{
			name:   "Simple addition",
			source: "5 + 1",
			expectedBytecode: Bytecode{
				Instructions:  []byte{byte(OP_CONSTANT), 0, 0, byte(OP_CONSTANT), 0, 1, byte(OP_ADD), byte(OP_END)},
				ConstantsPool: []any{int64(5), int64(1)},
			},
		},
		{
			name:   "Multiplication",
			source: "5 * 3",
			expectedBytecode: Bytecode{
				Instructions:  []byte{byte(OP_CONSTANT), 0, 0, byte(OP_CONSTANT), 0, 1, byte(OP_MULTIPLY), byte(OP_END)},
				ConstantsPool: []any{int64(5), int64(3)},
			},
		},
		{
			name:   "Negation",
			source: "-5",
			expectedBytecode: Bytecode{
				Instructions:  []byte{byte(OP_CONSTANT), 0, 0, byte(OP_NEGATE), byte(OP_END)},
				ConstantsPool: []any{int64(5)},
			},
		},
		{
			name:   "Complex expression",
			source: "5 * 3 + 2",
			expectedBytecode: Bytecode{
				Instructions:  []byte{byte(OP_CONSTANT), 0, 0, byte(OP_CONSTANT), 0, 1, byte(OP_MULTIPLY), byte(OP_CONSTANT), 0, 2, byte(OP_ADD), byte(OP_END)},
				ConstantsPool: []any{int64(5), int64(3), int64(2)},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			lex := lexer.New(tt.source)
			tokens, err := lex.Scan()
			if err != nil {
				t.Fatalf("lexing failed: %v", err)
			}

			parser := parser.Make(tokens)
			statements, parseErrors := parser.Parse()
			if len(parseErrors) > 0 {
				t.Fatalf("parsing failed: %v", parseErrors[0])
			}

			compiler := NewASTCompiler()
			bytecode, err := compiler.CompileAST(statements)
			if err != nil {
				t.Fatalf("compilation failed: %v", err)
			}

			// Verify the bytecode matches expected
			if len(bytecode.Instructions) != len(tt.expectedBytecode.Instructions) {
				t.Errorf("bytecode length mismatch - got: %d, want: %d", len(bytecode.Instructions), len(tt.expectedBytecode.Instructions))
			}

			for i, instr := range bytecode.Instructions {
				if instr != tt.expectedBytecode.Instructions[i] {
					t.Errorf("instruction mismatch at index %d - got: %d, want: %d", i, instr, tt.expectedBytecode.Instructions[i])
				}
			}

			if len(bytecode.ConstantsPool) != len(tt.expectedBytecode.ConstantsPool) {
				t.Errorf("constants pool length mismatch - got: %d, want: %d", len(bytecode.ConstantsPool), len(tt.expectedBytecode.ConstantsPool))
			}

			for i, constant := range bytecode.ConstantsPool {
				if constant != tt.expectedBytecode.ConstantsPool[i] {
					t.Errorf("constant mismatch at index %d - got: %v, want: %v", i, constant, tt.expectedBytecode.ConstantsPool[i])
				}
			}
		})
	}
}

// TestPipelineWithParser demonstrates integration with the parser package
// This ensures the AST produced by the parser is compatible with the ASTCompiler
func TestPipelineWithParser(t *testing.T) {
	// Create a simple arithmetic expression AST manually
	five := ast.Literal{Value: int64(5)}
	three := ast.Literal{Value: int64(3)}

	binaryExpr := ast.Binary{
		Left:     five,
		Operator: token.CreateToken(token.MULT, 0, 0),
		Right:    three,
	}

	exprStmt := ast.ExpressionStmt{
		Expression: binaryExpr,
	}

	statements := []ast.Stmt{exprStmt}

	// Compile the AST to bytecode
	compiler := NewASTCompiler()
	bytecode, err := compiler.CompileAST(statements)
	if err != nil {
		t.Fatalf("compilation failed: %v", err)
	}

	// Verify the bytecode is correct for 5 * 3
	if len(bytecode.Instructions) != 8 {
		t.Errorf("bytecode length mismatch - got: %d, want: 8", len(bytecode.Instructions))
	}

	if len(bytecode.ConstantsPool) != 2 {
		t.Errorf("constants pool length mismatch - got: %d, want: 2", len(bytecode.ConstantsPool))
	}

	if bytecode.ConstantsPool[0] != int64(5) {
		t.Errorf("first constant mismatch - got: %v, want: 5", bytecode.ConstantsPool[0])
	}

	if bytecode.ConstantsPool[1] != int64(3) {
		t.Errorf("second constant mismatch - got: %v, want: 3", bytecode.ConstantsPool[1])
	}
}
