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

	if len(got.ConstantsPool) != len(want.ConstantsPool) {
		t.Errorf("computed constants pool has a different length than the expected constants pool - got: %d, want: %d", len(got.ConstantsPool), len(want.ConstantsPool))
	}
	for i, constant := range got.ConstantsPool {
		if constant != want.ConstantsPool[i] {
			t.Errorf("computed constant does not equal expected constant at index %d - want: %v, got: %v", i, want.ConstantsPool[i], constant)
		}
	}
}

func TestASTCompilerVisitIfStmt(t *testing.T) {
	tests := []struct {
		name  string
		stmts []ast.Stmt
		want  Bytecode
	}{
		{
			name: "if(true){print(1)}else{print(2)}",
			stmts: []ast.Stmt{
				ast.IfStmt{
					Condition: ast.Literal{Value: true},
					Then:      ast.PrintStmt{Expression: ast.Literal{Value: int64(1)}},
					Else:      ast.PrintStmt{Expression: ast.Literal{Value: int64(2)}},
				},
			},
			want: Bytecode{
				Instructions: []byte{
					byte(OP_CONSTANT), 0, 0, // true
					byte(OP_JUMP_IF_FALSE), 0, 13, // jump to else (offset 13)
					byte(OP_CONSTANT), 0, 1, // 1
					byte(OP_PRINT),
					byte(OP_JUMP), 0, 17, // jump to end (offset 17)
					byte(OP_CONSTANT), 0, 2, // 2
					byte(OP_PRINT),
					byte(OP_POP),
					byte(OP_END),
				},
				ConstantsPool: []any{true, int64(1), int64(2)},
			},
		},
		{
			name: "if(false){print(1)}else{print(2)}",
			stmts: []ast.Stmt{
				ast.IfStmt{
					Condition: ast.Literal{Value: false},
					Then:      ast.PrintStmt{Expression: ast.Literal{Value: int64(1)}},
					Else:      ast.PrintStmt{Expression: ast.Literal{Value: int64(2)}},
				},
			},
			want: Bytecode{
				Instructions: []byte{
					byte(OP_CONSTANT), 0, 0, // false
					byte(OP_JUMP_IF_FALSE), 0, 13, // jump to else (offset 13)
					byte(OP_CONSTANT), 0, 1, // 1
					byte(OP_PRINT),
					byte(OP_JUMP), 0, 17, // jump to end (offset 17)
					byte(OP_CONSTANT), 0, 2, // 2
					byte(OP_PRINT),
					byte(OP_POP),
					byte(OP_END),
				},
				ConstantsPool: []any{false, int64(1), int64(2)},
			},
		},
		{

			name: "if (true){print(42)}",
			stmts: []ast.Stmt{
				ast.IfStmt{
					Condition: ast.Literal{Value: true},
					Then:      ast.PrintStmt{Expression: ast.Literal{Value: int64(42)}},
					Else:      nil,
				},
			},
			want: Bytecode{
				Instructions: []byte{
					byte(OP_CONSTANT), 0, 0, // true
					byte(OP_JUMP_IF_FALSE), 0, 10, // jump past then (offset 10)
					byte(OP_CONSTANT), 0, 1, // 42
					byte(OP_PRINT),
					byte(OP_POP),
					byte(OP_END),
				},
				ConstantsPool: []any{true, int64(42)},
			},
		},
		{

			name: "if (false){print(42)}",
			stmts: []ast.Stmt{
				ast.IfStmt{
					Condition: ast.Literal{Value: false},
					Then:      ast.PrintStmt{Expression: ast.Literal{Value: int64(42)}},
					Else:      nil,
				},
			},
			want: Bytecode{
				Instructions: []byte{
					byte(OP_CONSTANT), 0, 0, // false
					byte(OP_JUMP_IF_FALSE), 0, 10, // jump past then (offset 10)
					byte(OP_CONSTANT), 0, 1, // 42
					byte(OP_PRINT),
					byte(OP_POP),
					byte(OP_END),
				},
				ConstantsPool: []any{false, int64(42)},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			compiler := NewASTCompiler()
			bytecode, err := compiler.CompileAST(tt.stmts)
			if err != nil {
				t.Fatalf("compilation error: %v", err)
			}
			assertBytecodeEquals(t, bytecode, tt.want)
		})
	}
}

// TestASTCompilerLocalVariableDeclaration tests local variable declaration and initialization
func TestASTCompilerLocalVariableDeclaration(t *testing.T) {
	tests := []struct {
		name  string
		stmts []ast.Stmt
		want  Bytecode
	}{
		{
			// var x = 5; in a block scope
			name: "simple local variable declaration",
			stmts: []ast.Stmt{
				ast.BlockStmt{
					Statements: []ast.Stmt{
						ast.VarStmt{
							Name:        token.Token{Lexeme: "x", TokenType: token.IDENTIFIER},
							Initializer: ast.Literal{Value: int64(5)},
						},
					},
				},
			},
			want: Bytecode{
				Instructions: []byte{
					byte(OP_CONSTANT), 0, 0, // 5
					byte(OP_SET_LOCAL), 0, 0, // set x in slot 0
					byte(OP_SCOPE_EXIT), 0, 1, // exit scope, pop 1 local variable
					byte(OP_END),
				},
				ConstantsPool: []any{int64(5)},
			},
		},
		{
			// var x = 5; var y = 10; in a block scope
			name: "multiple local variable declarations",
			stmts: []ast.Stmt{
				ast.BlockStmt{
					Statements: []ast.Stmt{
						ast.VarStmt{
							Name:        token.Token{Lexeme: "x", TokenType: token.IDENTIFIER},
							Initializer: ast.Literal{Value: int64(5)},
						},
						ast.VarStmt{
							Name:        token.Token{Lexeme: "y", TokenType: token.IDENTIFIER},
							Initializer: ast.Literal{Value: int64(10)},
						},
					},
				},
			},
			want: Bytecode{
				Instructions: []byte{
					byte(OP_CONSTANT), 0, 0, // 5
					byte(OP_SET_LOCAL), 0, 0, // set x in slot 0
					byte(OP_CONSTANT), 0, 1, // 10
					byte(OP_SET_LOCAL), 0, 1, // set y in slot 1
					byte(OP_SCOPE_EXIT), 0, 2, // exit scope, pop 2 local variables
					byte(OP_END),
				},
				ConstantsPool: []any{int64(5), int64(10)},
			},
		},
		{
			// var x; (declaration without initializer)
			name: "local variable declaration without initializer",
			stmts: []ast.Stmt{
				ast.BlockStmt{
					Statements: []ast.Stmt{
						ast.VarStmt{
							Name:        token.Token{Lexeme: "x", TokenType: token.IDENTIFIER},
							Initializer: nil,
						},
					},
				},
			},
			want: Bytecode{
				Instructions: []byte{
					byte(OP_CONSTANT), 0, 0, // nil
					byte(OP_SET_LOCAL), 0, 0, // set x in slot 0
					byte(OP_SCOPE_EXIT), 0, 1, // exit scope, pop 1 local variable
					byte(OP_END),
				},
				ConstantsPool: []any{nil},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			compiler := NewASTCompiler()
			bytecode, err := compiler.CompileAST(tt.stmts)
			if err != nil {
				t.Fatalf("compilation error: %v", err)
			}
			assertBytecodeEquals(t, bytecode, tt.want)
		})
	}
}

// TestASTCompilerLocalVariableAccess tests local variable access and usage
func TestASTCompilerLocalVariableAccess(t *testing.T) {
	tests := []struct {
		name  string
		stmts []ast.Stmt
		want  Bytecode
	}{
		{
			// var x = 5; x + 3
			name: "local variable access in expression",
			stmts: []ast.Stmt{
				ast.BlockStmt{
					Statements: []ast.Stmt{
						ast.VarStmt{
							Name:        token.Token{Lexeme: "x", TokenType: token.IDENTIFIER},
							Initializer: ast.Literal{Value: int64(5)},
						},
						ast.ExpressionStmt{
							Expression: ast.Binary{
								Left:     ast.Variable{Name: token.Token{Lexeme: "x", TokenType: token.IDENTIFIER}},
								Operator: token.Token{TokenType: token.ADD},
								Right:    ast.Literal{Value: int64(3)},
							},
						},
					},
				},
			},
			want: Bytecode{
				Instructions: []byte{
					byte(OP_CONSTANT), 0, 0, // 5
					byte(OP_SET_LOCAL), 0, 0, // set x in slot 0
					byte(OP_GET_LOCAL), 0, 0, // load x
					byte(OP_CONSTANT), 0, 1, // 3
					byte(OP_ADD),
					byte(OP_SCOPE_EXIT), 0, 1, // exit scope, pop 1 local variable
					byte(OP_END),
				},
				ConstantsPool: []any{int64(5), int64(3)},
			},
		},
		{
			// var x = 5; var y = 10; x + y
			name: "multiple local variables in expression",
			stmts: []ast.Stmt{
				ast.BlockStmt{
					Statements: []ast.Stmt{
						ast.VarStmt{
							Name:        token.Token{Lexeme: "x", TokenType: token.IDENTIFIER},
							Initializer: ast.Literal{Value: int64(5)},
						},
						ast.VarStmt{
							Name:        token.Token{Lexeme: "y", TokenType: token.IDENTIFIER},
							Initializer: ast.Literal{Value: int64(10)},
						},
						ast.ExpressionStmt{
							Expression: ast.Binary{
								Left:     ast.Variable{Name: token.Token{Lexeme: "x", TokenType: token.IDENTIFIER}},
								Operator: token.Token{TokenType: token.ADD},
								Right:    ast.Variable{Name: token.Token{Lexeme: "y", TokenType: token.IDENTIFIER}},
							},
						},
					},
				},
			},
			want: Bytecode{
				Instructions: []byte{
					byte(OP_CONSTANT), 0, 0, // 5
					byte(OP_SET_LOCAL), 0, 0, // set x in slot 0
					byte(OP_CONSTANT), 0, 1, // 10
					byte(OP_SET_LOCAL), 0, 1, // set y in slot 1
					byte(OP_GET_LOCAL), 0, 0, // load x
					byte(OP_GET_LOCAL), 0, 1, // load y
					byte(OP_ADD),
					byte(OP_SCOPE_EXIT), 0, 2, // exit scope, pop 2 local variables
					byte(OP_END),
				},
				ConstantsPool: []any{int64(5), int64(10)},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			compiler := NewASTCompiler()
			bytecode, err := compiler.CompileAST(tt.stmts)
			if err != nil {
				t.Fatalf("compilation error: %v", err)
			}
			assertBytecodeEquals(t, bytecode, tt.want)
		})
	}
}

// TestASTCompilerLocalVariableAssignment tests local variable assignment
func TestASTCompilerLocalVariableAssignment(t *testing.T) {
	tests := []struct {
		name  string
		stmts []ast.Stmt
		want  Bytecode
	}{
		{
			// var x = 5; x = 10;
			name: "local variable reassignment",
			stmts: []ast.Stmt{
				ast.BlockStmt{
					Statements: []ast.Stmt{
						ast.VarStmt{
							Name:        token.Token{Lexeme: "x", TokenType: token.IDENTIFIER},
							Initializer: ast.Literal{Value: int64(5)},
						},
						ast.ExpressionStmt{
							Expression: ast.Assign{
								Name:  token.Token{Lexeme: "x", TokenType: token.IDENTIFIER},
								Value: ast.Literal{Value: int64(10)},
							},
						},
					},
				},
			},
			want: Bytecode{
				Instructions: []byte{
					byte(OP_CONSTANT), 0, 0, // 5
					byte(OP_SET_LOCAL), 0, 0, // set x in slot 0
					byte(OP_CONSTANT), 0, 1, // 10
					byte(OP_SET_LOCAL), 0, 0, // reassign x (same slot 0)
					byte(OP_SCOPE_EXIT), 0, 1, // exit scope, pop 1 local variable
					byte(OP_END),
				},
				ConstantsPool: []any{int64(5), int64(10)},
			},
		},
		{
			// var x = 5; x = x + 3;
			name: "local variable reassignment with self-reference",
			stmts: []ast.Stmt{
				ast.BlockStmt{
					Statements: []ast.Stmt{
						ast.VarStmt{
							Name:        token.Token{Lexeme: "x", TokenType: token.IDENTIFIER},
							Initializer: ast.Literal{Value: int64(5)},
						},
						ast.ExpressionStmt{
							Expression: ast.Assign{
								Name: token.Token{Lexeme: "x", TokenType: token.IDENTIFIER},
								Value: ast.Binary{
									Left:     ast.Variable{Name: token.Token{Lexeme: "x", TokenType: token.IDENTIFIER}},
									Operator: token.Token{TokenType: token.ADD},
									Right:    ast.Literal{Value: int64(3)},
								},
							},
						},
					},
				},
			},
			want: Bytecode{
				Instructions: []byte{
					byte(OP_CONSTANT), 0, 0, // 5
					byte(OP_SET_LOCAL), 0, 0, // set x in slot 0
					byte(OP_GET_LOCAL), 0, 0, // load x
					byte(OP_CONSTANT), 0, 1, // 3
					byte(OP_ADD),
					byte(OP_SET_LOCAL), 0, 0, // reassign x (same slot 0)
					byte(OP_SCOPE_EXIT), 0, 1, // exit scope, pop 1 local variable
					byte(OP_END),
				},
				ConstantsPool: []any{int64(5), int64(3)},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			compiler := NewASTCompiler()
			bytecode, err := compiler.CompileAST(tt.stmts)
			if err != nil {
				t.Fatalf("compilation error: %v", err)
			}
			assertBytecodeEquals(t, bytecode, tt.want)
		})
	}
}

// TestASTCompilerNestedScopes tests nested local scopes
func TestASTCompilerNestedScopes(t *testing.T) {
	tests := []struct {
		name  string
		stmts []ast.Stmt
		want  Bytecode
	}{
		{
			// Outer scope: var x = 5
			// Inner scope: var y = 10
			// Test proper scope nesting and cleanup
			name: "nested block scopes",
			stmts: []ast.Stmt{
				ast.BlockStmt{
					Statements: []ast.Stmt{
						ast.VarStmt{
							Name:        token.Token{Lexeme: "x", TokenType: token.IDENTIFIER},
							Initializer: ast.Literal{Value: int64(5)},
						},
						ast.BlockStmt{
							Statements: []ast.Stmt{
								ast.VarStmt{
									Name:        token.Token{Lexeme: "y", TokenType: token.IDENTIFIER},
									Initializer: ast.Literal{Value: int64(10)},
								},
							},
						},
					},
				},
			},
			want: Bytecode{
				Instructions: []byte{
					byte(OP_CONSTANT), 0, 0, // 5
					byte(OP_SET_LOCAL), 0, 0, // set x in slot 0
					byte(OP_CONSTANT), 0, 1, // 10
					byte(OP_SET_LOCAL), 0, 1, // set y in slot 1 (inner scope)
					byte(OP_SCOPE_EXIT), 0, 1, // exit inner scope, pop 1 local (y)
					byte(OP_SCOPE_EXIT), 0, 1, // exit outer scope, pop 1 local (x)
					byte(OP_END),
				},
				ConstantsPool: []any{int64(5), int64(10)},
			},
		},
		{
			// Test deeply nested scopes (3 levels)
			// {var a = 1 { var b = 2 { var c = 3 } }}
			name: "deeply nested scopes",
			stmts: []ast.Stmt{
				ast.BlockStmt{
					Statements: []ast.Stmt{
						ast.VarStmt{
							Name:        token.Token{Lexeme: "a", TokenType: token.IDENTIFIER},
							Initializer: ast.Literal{Value: int64(1)},
						},
						ast.BlockStmt{
							Statements: []ast.Stmt{
								ast.VarStmt{
									Name:        token.Token{Lexeme: "b", TokenType: token.IDENTIFIER},
									Initializer: ast.Literal{Value: int64(2)},
								},
								ast.BlockStmt{
									Statements: []ast.Stmt{
										ast.VarStmt{
											Name:        token.Token{Lexeme: "c", TokenType: token.IDENTIFIER},
											Initializer: ast.Literal{Value: int64(3)},
										},
									},
								},
							},
						},
					},
				},
			},
			want: Bytecode{
				Instructions: []byte{
					byte(OP_CONSTANT), 0, 0, // 1
					byte(OP_SET_LOCAL), 0, 0, // set a in slot 0
					byte(OP_CONSTANT), 0, 1, // 2
					byte(OP_SET_LOCAL), 0, 1, // set b in slot 1
					byte(OP_CONSTANT), 0, 2, // 3
					byte(OP_SET_LOCAL), 0, 2, // set c in slot 2
					byte(OP_SCOPE_EXIT), 0, 1, // exit innermost scope, pop 1 (c)
					byte(OP_SCOPE_EXIT), 0, 1, // exit middle scope, pop 1 (b)
					byte(OP_SCOPE_EXIT), 0, 1, // exit outer scope, pop 1 (a)
					byte(OP_END),
				},
				ConstantsPool: []any{int64(1), int64(2), int64(3)},
			},
		},
		{
			// Test accessing outer scope variable from inner scope
			// var x = 5 {x + 3}
			name: "access outer scope variable from inner scope",
			stmts: []ast.Stmt{
				ast.VarStmt{
					Name:        token.Token{Lexeme: "x", TokenType: token.IDENTIFIER},
					Initializer: ast.Literal{Value: int64(5)},
				},
				ast.BlockStmt{
					Statements: []ast.Stmt{
						ast.ExpressionStmt{
							Expression: ast.Binary{
								Left:     ast.Variable{Name: token.Token{Lexeme: "x", TokenType: token.IDENTIFIER}},
								Operator: token.Token{TokenType: token.ADD},
								Right:    ast.Literal{Value: int64(3)},
							},
						},
					},
				},
			},
			want: Bytecode{
				Instructions: []byte{
					byte(OP_CONSTANT), 0, 0, // declare 5
					byte(OP_SET_GLOBAL), 0, 0, // assign 5 to x in global scope
					byte(OP_GET_GLOBAL), 0, 0, // load x from global scope
					byte(OP_CONSTANT), 0, 1, // declare 3
					byte(OP_ADD),
					byte(OP_END),
				},
				ConstantsPool: []any{int64(5), int64(3)},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			compiler := NewASTCompiler()
			bytecode, err := compiler.CompileAST(tt.stmts)
			if err != nil {
				t.Fatalf("compilation error: %v", err)
			}
			assertBytecodeEquals(t, bytecode, tt.want)
		})
	}
}

// TestASTCompilerScopeWithIfStatement tests local scopes used with if statements
func TestASTCompilerScopeWithIfStatement(t *testing.T) {
	tests := []struct {
		name  string
		stmts []ast.Stmt
		want  Bytecode
	}{
		{
			// var x = 5; if (x > 3) { var y = 10 }
			name: "local variable in if block",
			stmts: []ast.Stmt{
				ast.VarStmt{
					Name:        token.Token{Lexeme: "x", TokenType: token.IDENTIFIER},
					Initializer: ast.Literal{Value: int64(5)},
				},
				ast.IfStmt{
					Condition: ast.Grouping{
						Expression: ast.Binary{
							Left:     ast.Variable{Name: token.Token{Lexeme: "x", TokenType: token.IDENTIFIER}},
							Operator: token.Token{TokenType: token.LARGER},
							Right:    ast.Literal{Value: int64(3)},
						},
					},
					Then: ast.BlockStmt{
						Statements: []ast.Stmt{
							ast.VarStmt{
								Name:        token.Token{Lexeme: "y", TokenType: token.IDENTIFIER},
								Initializer: ast.Literal{Value: int64(10)},
							},
						},
					},
					Else: nil,
				},
			},
			want: Bytecode{
				Instructions: []byte{
					byte(OP_CONSTANT), 0, 0, // 5
					byte(OP_SET_GLOBAL), 0, 0, // set x
					byte(OP_GET_GLOBAL), 0, 0, // load x
					byte(OP_CONSTANT), 0, 1, // declare 3
					byte(OP_LARGER),
					byte(OP_JUMP_IF_FALSE), 0, 25, // jump if false to end
					byte(OP_CONSTANT), 0, 2, // 10
					byte(OP_SET_LOCAL), 0, 0, // set y
					byte(OP_SCOPE_EXIT), 0, 1, // exit if block scope, pop 1 (y)
					byte(OP_POP), // pop condition
					byte(OP_END),
				},
				ConstantsPool: []any{int64(5), int64(3), int64(10)},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			compiler := NewASTCompiler()
			bytecode, err := compiler.CompileAST(tt.stmts)
			if err != nil {
				t.Fatalf("compilation error: %v", err)
			}
			assertBytecodeEquals(t, bytecode, tt.want)
		})
	}
}

// TestASTCompilerLocalVariableShadowing tests that local variable shadowing is correctly compiled
func TestASTCompilerLocalVariableShadowing(t *testing.T) {
	tests := []struct {
		name  string
		stmts []ast.Stmt
		want  Bytecode
	}{
		{
			// Outer scope: var x = 5
			// Inner scope: var x = 10 (shadows outer x)
			// Test that inner x gets a different slot
			name: "variable shadowing in nested scopes",
			stmts: []ast.Stmt{
				ast.BlockStmt{
					Statements: []ast.Stmt{
						ast.VarStmt{
							Name:        token.Token{Lexeme: "x", TokenType: token.IDENTIFIER},
							Initializer: ast.Literal{Value: int64(5)},
						},
						ast.BlockStmt{
							Statements: []ast.Stmt{
								ast.VarStmt{
									Name:        token.Token{Lexeme: "x", TokenType: token.IDENTIFIER},
									Initializer: ast.Literal{Value: int64(10)},
								},
								ast.ExpressionStmt{
									Expression: ast.Variable{Name: token.Token{Lexeme: "x", TokenType: token.IDENTIFIER}},
								},
							},
						},
					},
				},
			},
			want: Bytecode{
				Instructions: []byte{
					byte(OP_CONSTANT), 0, 0, // 5
					byte(OP_SET_LOCAL), 0, 0, // set x in slot 0 (outer)
					byte(OP_CONSTANT), 0, 1, // 10
					byte(OP_SET_LOCAL), 0, 1, // set x in slot 1 (inner, shadows outer)
					byte(OP_GET_LOCAL), 0, 1, // load x from slot 1 (shadowed, inner)
					byte(OP_SCOPE_EXIT), 0, 1, // exit inner scope, pop 1 (inner x)
					byte(OP_SCOPE_EXIT), 0, 1, // exit outer scope, pop 1 (outer x)
					byte(OP_END),
				},
				ConstantsPool: []any{int64(5), int64(10)},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			compiler := NewASTCompiler()
			bytecode, err := compiler.CompileAST(tt.stmts)
			if err != nil {
				t.Fatalf("compilation error: %v", err)
			}
			assertBytecodeEquals(t, bytecode, tt.want)
		})
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
		bytecode   *Bytecode
		expected   string
	}{
		{
			name: "Nested arithmetic",
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
			if tt.statements != nil {
				_, err := compiler.CompileAST(tt.statements)
				if err != nil {
					t.Errorf("compilation error occurred: %s", err.Error())
				}
			}

			if tt.bytecode != nil {
				compiler.bytecode = *tt.bytecode
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

func TestASTCompilerVisitWhileStmt(t *testing.T) {
	tests := []struct {
		name  string
		stmts []ast.Stmt
		want  Bytecode
	}{
		{

			name: "while(true){print(1)}",
			stmts: []ast.Stmt{
				ast.WhileStmt{
					Condition: ast.Literal{Value: true},
					Body: ast.BlockStmt{
						Statements: []ast.Stmt{
							ast.PrintStmt{Expression: ast.Literal{Value: int64(1)}},
						},
					},
				},
			},
			want: Bytecode{
				Instructions: []byte{
					byte(OP_CONSTANT), 0, 0, // true (loop start)
					byte(OP_JUMP_IF_FALSE), 0, 14, // jump to end if false (offset 14)
					byte(OP_CONSTANT), 0, 1, // 1
					byte(OP_PRINT),
					byte(OP_POP),        // pop condition
					byte(OP_JUMP), 0, 0, // jump back to loop start (offset 0)
					byte(OP_POP), // pop condition at end
					byte(OP_END),
				},
				ConstantsPool: []any{true, int64(1)},
			},
		},
		{
			// While loop with false condition: while(false){print(1)}
			// This tests that the loop body is skipped when condition is false
			name: "while(false){print(1)}",
			stmts: []ast.Stmt{
				ast.WhileStmt{
					Condition: ast.Literal{Value: false},
					Body: ast.BlockStmt{
						Statements: []ast.Stmt{
							ast.PrintStmt{Expression: ast.Literal{Value: int64(1)}},
						},
					},
				},
			},
			want: Bytecode{
				Instructions: []byte{
					byte(OP_CONSTANT), 0, 0, // false (loop start)
					byte(OP_JUMP_IF_FALSE), 0, 14, // jump to end if false (offset 14)
					byte(OP_CONSTANT), 0, 1, // 1
					byte(OP_PRINT),
					byte(OP_POP),        // pop condition
					byte(OP_JUMP), 0, 0, // jump back to loop start (offset 0)
					byte(OP_POP), // pop condition at end
					byte(OP_END),
				},
				ConstantsPool: []any{false, int64(1)},
			},
		},
		{
			name: "while(1 < 5){print('true')}",
			stmts: []ast.Stmt{
				ast.WhileStmt{
					Condition: ast.Binary{
						Left:     ast.Literal{Value: int64(1)},
						Operator: token.CreateToken(token.LESS, 0, 0),
						Right:    ast.Literal{Value: int64(5)},
					},
					Body: ast.BlockStmt{
						Statements: []ast.Stmt{
							ast.PrintStmt{Expression: ast.Literal{Value: "true"}},
						},
					},
				},
			},
			want: Bytecode{
				Instructions: []byte{
					byte(OP_CONSTANT), 0, 0, // 1
					byte(OP_CONSTANT), 0, 1, // 5
					byte(OP_LESS),                 // 1 < 5
					byte(OP_JUMP_IF_FALSE), 0, 18, // jump to end if false (offset 25)
					byte(OP_CONSTANT), 0, 2, // "true"
					byte(OP_PRINT),
					byte(OP_POP),        // pop condition
					byte(OP_JUMP), 0, 0, // jump back to loop start (offset 6)
					byte(OP_POP), // pop condition at end
					byte(OP_END),
				},
				ConstantsPool: []any{int64(1), int64(5), "true"},
				NameConstants: []string{},
			},
		},
		{
			name: "var x = 1 while(x < 5){print(x)}",
			stmts: []ast.Stmt{
				ast.VarStmt{
					Name:        token.CreateLiteralToken(token.IDENTIFIER, "x", "x", 0, 0),
					Initializer: ast.Literal{Value: int64(1)},
				},
				ast.WhileStmt{
					Condition: ast.Grouping{
						Expression: ast.Binary{
							Left:     ast.Variable{Name: token.CreateLiteralToken(token.IDENTIFIER, "x", "x", 0, 0)},
							Operator: token.CreateToken(token.LESS, 0, 0),
							Right:    ast.Literal{Value: int64(5)},
						},
					},
					Body: ast.BlockStmt{
						Statements: []ast.Stmt{
							ast.PrintStmt{
								Expression: ast.Grouping{
									Expression: ast.Variable{Name: token.CreateLiteralToken(token.IDENTIFIER, "x", "x", 0, 0)},
								},
							},
						},
					},
				},
			},
			want: Bytecode{
				Instructions: []byte{
					byte(OP_CONSTANT), 0, 0, // 1
					byte(OP_SET_GLOBAL), 0, 0, // 1
					byte(OP_GET_GLOBAL), 0, 0, // 1
					byte(OP_CONSTANT), 0, 1, // 5
					byte(OP_LESS),                 // 1 < 5
					byte(OP_JUMP_IF_FALSE), 0, 24, // jump to end if false (offset 25)
					byte(OP_GET_GLOBAL), 0, 0, // 1
					byte(OP_PRINT),
					byte(OP_POP),        // pop condition
					byte(OP_JUMP), 0, 6, // jump back to loop start (offset 6)
					byte(OP_POP), // pop condition at end
					byte(OP_END),
				},
				ConstantsPool: []any{int64(1), int64(5)},
				NameConstants: []string{"x"},
			},
		},
		{
			// Nested while loop: while(true){while(false){print(1)}}
			// This tests nested loop structures
			name: "nested while loops",
			stmts: []ast.Stmt{
				ast.WhileStmt{
					Condition: ast.Literal{Value: true},
					Body: ast.BlockStmt{
						Statements: []ast.Stmt{
							ast.WhileStmt{
								Condition: ast.Literal{Value: false},
								Body: ast.BlockStmt{
									Statements: []ast.Stmt{
										ast.PrintStmt{Expression: ast.Literal{Value: int64(1)}},
									},
								},
							},
						},
					},
				},
			},
			want: Bytecode{
				Instructions: []byte{
					byte(OP_CONSTANT), 0, 0, // true (outer loop start)
					byte(OP_JUMP_IF_FALSE), 0, 25, // jump to end of outer loop if false
					// Inner loop
					byte(OP_CONSTANT), 0, 1, // false (inner loop start)
					byte(OP_JUMP_IF_FALSE), 0, 20, // jump to end of inner loop if false
					byte(OP_CONSTANT), 0, 2, // 1
					byte(OP_PRINT),
					byte(OP_POP),        // pop condition
					byte(OP_JUMP), 0, 6, // jump back to inner loop start
					byte(OP_POP),        // pop condition at end of inner loop
					byte(OP_POP),        // pop condition (outer)
					byte(OP_JUMP), 0, 0, // jump back to outer loop start
					byte(OP_POP), // pop condition at end of outer loop
					byte(OP_END),
				},
				ConstantsPool: []any{true, false, int64(1)},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			compiler := NewASTCompiler()
			bytecode, err := compiler.CompileAST(tt.stmts)
			if err != nil {
				t.Fatalf("compilation error: %v", err)
			}
			assertBytecodeEquals(t, bytecode, tt.want)
		})
	}
}
