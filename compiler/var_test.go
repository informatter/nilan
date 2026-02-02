package compiler

import (
	"nilan/ast"
	"nilan/token"
	"testing"
)

func TestCompilerVariableBehavior(t *testing.T) {
	tests := []struct {
		name       string
		statements []ast.Stmt
		hasError   bool
	}{
		{
			name: "var declared without initializer then accessed -> error",
			statements: []ast.Stmt{
				ast.VarStmt{Name: token.CreateLiteralToken(token.IDENTIFIER, nil, "a", 0, 0)},
				ast.PrintStmt{Expression: ast.Variable{Name: token.CreateLiteralToken(token.IDENTIFIER, nil, "a", 0, 0)}},
			},
			hasError: true,
		},
		{
			name: "var declared with initializer then accessed -> success",
			statements: []ast.Stmt{
				ast.VarStmt{Name: token.CreateLiteralToken(token.IDENTIFIER, nil, "a", 0, 0), Initializer: ast.Literal{Value: int64(0)}},
				ast.PrintStmt{Expression: ast.Variable{Name: token.CreateLiteralToken(token.IDENTIFIER, nil, "a", 0, 0)}},
			},
			hasError: false,
		},
		{
			name: "access undeclared variable -> error",
			statements: []ast.Stmt{
				ast.PrintStmt{Expression: ast.Variable{Name: token.CreateLiteralToken(token.IDENTIFIER, nil, "c", 0, 0)}},
			},
			hasError: true,
		},
		{
			name: "redeclaration of variable -> error",
			statements: []ast.Stmt{
				ast.VarStmt{Name: token.CreateLiteralToken(token.IDENTIFIER, nil, "a", 0, 0)},
				ast.VarStmt{Name: token.CreateLiteralToken(token.IDENTIFIER, nil, "a", 0, 0), Initializer: ast.Literal{Value: int64(9)}},
			},
			hasError: true,
		},
		{
			name: "assignment to existing variable -> success",
			statements: []ast.Stmt{
				ast.VarStmt{Name: token.CreateLiteralToken(token.IDENTIFIER, nil, "a", 0, 0)},
				ast.ExpressionStmt{Expression: ast.Assign{Name: token.CreateLiteralToken(token.IDENTIFIER, nil, "a", 0, 0), Value: ast.Literal{Value: int64(1)}}},
			},
			hasError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			compiler := NewASTCompiler()
			_, err := compiler.CompileAST(tt.statements)
			if tt.hasError && err == nil {
				t.Errorf("expected error but got nil")
			}
			if !tt.hasError && err != nil {
				t.Errorf("unexpected compilation error: %s", err.Error())
			}
		})
	}
}
