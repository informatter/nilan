package parser

import (
	"nilan/ast"
	"nilan/lexer"
	"strings"
	"testing"
)

func parseSource(t *testing.T, source string, fromREPL bool) ([]ast.Stmt, []error) {
	t.Helper()

	lex := lexer.New(source)
	tokens, err := lex.Scan()
	if err != nil {
		t.Fatalf("lexing failed: %v", err)
	}

	p := Make(tokens)
	return p.Parse(fromREPL)
}

func assertParseErrorContains(t *testing.T, source string, expectedSubstring string) {
	t.Helper()

	_, errs := parseSource(t, source, false)
	if len(errs) == 0 {
		t.Fatalf("expected parse error for source %q, got none", source)
	}

	if !strings.Contains(errs[0].Error(), expectedSubstring) {
		t.Fatalf("expected first parse error to contain %q, got %q", expectedSubstring, errs[0].Error())
	}
}

func TestParseFuncDecl_NoParameters(t *testing.T) {
	source := `fn main() {}`

	stmts, errs := parseSource(t, source, false)
	if len(errs) != 0 {
		t.Fatalf("unexpected parse errors: %v", errs)
	}
	if len(stmts) != 1 {
		t.Fatalf("expected 1 top-level statement, got %d", len(stmts))
	}

	fn, ok := stmts[0].(ast.FunctionStmt)
	if !ok {
		t.Fatalf("expected ast.FunctionStmt, got %T", stmts[0])
	}
	if fn.Name.Lexeme != "main" {
		t.Fatalf("expected function name 'main', got %q", fn.Name.Lexeme)
	}
	if len(fn.Parameters) != 0 {
		t.Fatalf("expected 0 parameters, got %d", len(fn.Parameters))
	}
	if len(fn.Body.Statements) != 0 {
		t.Fatalf("expected empty body, got %d statements", len(fn.Body.Statements))
	}
}

func TestParseFuncDecl_WithParameters(t *testing.T) {
	source := `fn add(a, b) { print(a) }`

	stmts, errs := parseSource(t, source, false)
	if len(errs) != 0 {
		t.Fatalf("unexpected parse errors: %v", errs)
	}
	if len(stmts) != 1 {
		t.Fatalf("expected 1 top-level statement, got %d", len(stmts))
	}

	fn, ok := stmts[0].(ast.FunctionStmt)
	if !ok {
		t.Fatalf("expected ast.FunctionStmt, got %T", stmts[0])
	}

	if got, want := len(fn.Parameters), 2; got != want {
		t.Fatalf("expected %d parameters, got %d", want, got)
	}
	if fn.Parameters[0].Lexeme != "a" || fn.Parameters[1].Lexeme != "b" {
		t.Fatalf("expected parameters [a b], got [%s %s]", fn.Parameters[0].Lexeme, fn.Parameters[1].Lexeme)
	}
	if len(fn.Body.Statements) != 1 {
		t.Fatalf("expected 1 statement in function body, got %d", len(fn.Body.Statements))
	}
	if _, ok := fn.Body.Statements[0].(ast.PrintStmt); !ok {
		t.Fatalf("expected body statement to be ast.PrintStmt, got %T", fn.Body.Statements[0])
	}
}

func TestParseCallExpr_NoArguments(t *testing.T) {
	source := `foo()`

	stmts, errs := parseSource(t, source, false)
	if len(errs) != 0 {
		t.Fatalf("unexpected parse errors: %v", errs)
	}
	if len(stmts) != 1 {
		t.Fatalf("expected 1 statement, got %d", len(stmts))
	}

	exprStmt, ok := stmts[0].(ast.ExpressionStmt)
	if !ok {
		t.Fatalf("expected ast.ExpressionStmt, got %T", stmts[0])
	}

	call, ok := exprStmt.Expression.(ast.CallExpr)
	if !ok {
		t.Fatalf("expected ast.CallExpr, got %T", exprStmt.Expression)
	}

	callee, ok := call.Callee.(ast.VariableExpr)
	if !ok {
		t.Fatalf("expected call callee ast.VariableExpr, got %T", call.Callee)
	}
	if callee.Name.Lexeme != "foo" {
		t.Fatalf("expected callee 'foo', got %q", callee.Name.Lexeme)
	}
	if len(call.Arguments) != 0 {
		t.Fatalf("expected 0 arguments, got %d", len(call.Arguments))
	}
}

func TestParseCallExpr_WithArguments(t *testing.T) {
	source := `foo(1, 2)`

	stmts, errs := parseSource(t, source, false)
	if len(errs) != 0 {
		t.Fatalf("unexpected parse errors: %v", errs)
	}
	if len(stmts) != 1 {
		t.Fatalf("expected 1 statement, got %d", len(stmts))
	}

	exprStmt, ok := stmts[0].(ast.ExpressionStmt)
	if !ok {
		t.Fatalf("expected ast.ExpressionStmt, got %T", stmts[0])
	}
	call, ok := exprStmt.Expression.(ast.CallExpr)
	if !ok {
		t.Fatalf("expected ast.CallExpr, got %T", exprStmt.Expression)
	}

	if got, want := len(call.Arguments), 2; got != want {
		t.Fatalf("expected %d call arguments, got %d", want, got)
	}

	arg0, ok := call.Arguments[0].(ast.LiteralExpr)
	if !ok {
		t.Fatalf("expected first argument to be ast.LiteralExpr, got %T", call.Arguments[0])
	}
	arg1, ok := call.Arguments[1].(ast.LiteralExpr)
	if !ok {
		t.Fatalf("expected second argument to be ast.LiteralExpr, got %T", call.Arguments[1])
	}

	if arg0.Value != int64(1) {
		t.Fatalf("expected first argument value 1, got %v", arg0.Value)
	}
	if arg1.Value != int64(2) {
		t.Fatalf("expected second argument value 2, got %v", arg1.Value)
	}
}

func TestParseFunctionBody_CallExpression(t *testing.T) {
	source := `fn main() { foo(1, 2) }`

	stmts, errs := parseSource(t, source, false)
	if len(errs) != 0 {
		t.Fatalf("unexpected parse errors: %v", errs)
	}
	if len(stmts) != 1 {
		t.Fatalf("expected 1 top-level statement, got %d", len(stmts))
	}

	fn, ok := stmts[0].(ast.FunctionStmt)
	if !ok {
		t.Fatalf("expected ast.FunctionStmt, got %T", stmts[0])
	}
	if len(fn.Body.Statements) != 1 {
		t.Fatalf("expected 1 statement in function body, got %d", len(fn.Body.Statements))
	}

	exprStmt, ok := fn.Body.Statements[0].(ast.ExpressionStmt)
	if !ok {
		t.Fatalf("expected ast.ExpressionStmt, got %T", fn.Body.Statements[0])
	}
	if _, ok := exprStmt.Expression.(ast.CallExpr); !ok {
		t.Fatalf("expected function body expression to be ast.CallExpr, got %T", exprStmt.Expression)
	}
}

func TestParseCallExpr_MalformedCalls(t *testing.T) {
	tests := []struct {
		name              string
		source            string
		expectedErrorPart string
	}{
		{
			name:              "missing closing parenthesis",
			source:            `foo(1, 2`,
			expectedErrorPart: "expected ')' after arguments",
		},
		{
			name:              "leading comma in args",
			source:            `foo(, 1)`,
			expectedErrorPart: "Unrecognised expression",
		},
		{
			name:              "double comma in args",
			source:            `foo(1,, 2)`,
			expectedErrorPart: "Unrecognised expression",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assertParseErrorContains(t, tt.source, tt.expectedErrorPart)
		})
	}
}

func TestParseFuncDecl_MalformedParameterLists(t *testing.T) {
	tests := []struct {
		name              string
		source            string
		expectedErrorPart string
	}{
		{
			name:              "leading comma in parameters",
			source:            `fn add(, a) {}`,
			expectedErrorPart: "expected parameter name",
		},
		{
			name:              "trailing comma in parameters",
			source:            `fn add(a, ) {}`,
			expectedErrorPart: "expected parameter name",
		},
		{
			name:              "missing comma between parameters",
			source:            `fn add(a b) {}`,
			expectedErrorPart: "expected ')' after parameters",
		},
		{
			name:              "missing closing parenthesis",
			source:            `fn add(a, b {}`,
			expectedErrorPart: "expected ')' after parameters",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assertParseErrorContains(t, tt.source, tt.expectedErrorPart)
		})
	}
}
