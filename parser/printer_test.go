package parser

import (
	"encoding/json"
	"nilan/ast"
	"nilan/token"
	"os"
	"path/filepath"
	"testing"
)

func TestPrintASTJSON_PrintLiteral(t *testing.T) {
	stmts := []ast.Stmt{
		ast.PrintStmt{Expression: ast.Literal{Value: 42}},
	}

	jsonString, err := PrintASTJSON(stmts)
	if err != nil {
		t.Fatalf("PrintASTJSON error: %v", err)
	}

	var out []map[string]any
	if err := json.Unmarshal([]byte(jsonString), &out); err != nil {
		t.Fatalf("unmarshal json: %v", err)
	}

	if len(out) != 1 {
		t.Fatalf("expected 1 statement, got %d", len(out))
	}

	node := out[0]
	if typ, ok := node["type"].(string); !ok || typ != "PrintStmt" {
		t.Fatalf("expected type PrintStmt, got %v", node["type"])
	}

	expr := node["expression"]
	if num, ok := expr.(float64); !ok || num != 42 {
		t.Fatalf("expected expression 42, got %v", expr)
	}
}

func TestPrintASTJSON_VarStmt_NilInitializer(t *testing.T) {
	name := token.CreateLiteralToken(token.IDENTIFIER, nil, "x", 0, 0)
	stmts := []ast.Stmt{
		ast.VarStmt{Name: name, Initializer: nil},
	}

	jsonStr, err := PrintASTJSON(stmts)
	if err != nil {
		t.Fatalf("PrintASTJSON error: %v", err)
	}

	var out []map[string]any
	if err := json.Unmarshal([]byte(jsonStr), &out); err != nil {
		t.Fatalf("unmarshal json: %v", err)
	}

	if len(out) != 1 {
		t.Fatalf("expected 1 statement, got %d", len(out))
	}

	node := out[0]
	if typ, ok := node["type"].(string); !ok || typ != "VarStmt" {
		t.Fatalf("expected type VarStmt, got %v", node["type"])
	}

	if nameVal, ok := node["name"].(string); !ok || nameVal != "x" {
		t.Fatalf("expected name 'x', got %v", node["name"])
	}

	if initVal, exists := node["initializer"]; !exists || initVal != nil {
		t.Fatalf("expected initializer to be nil, got %v", initVal)
	}
}

func TestPrintASTJSON_BinaryExpression(t *testing.T) {
	stmts := []ast.Stmt{
		ast.ExpressionStmt{Expression: ast.Binary{
			Left:     ast.Literal{Value: 1},
			Operator: token.CreateToken(token.ADD, 0, 0),
			Right:    ast.Literal{Value: 2},
		}},
	}

	jsonStr, err := PrintASTJSON(stmts)
	if err != nil {
		t.Fatalf("PrintASTJSON error: %v", err)
	}

	var out []map[string]any
	if err := json.Unmarshal([]byte(jsonStr), &out); err != nil {
		t.Fatalf("unmarshal json: %v", err)
	}

	if len(out) != 1 {
		t.Fatalf("expected 1 statement, got %d", len(out))
	}

	node := out[0]
	if typ, ok := node["type"].(string); !ok || typ != "ExpressionStmt" {
		t.Fatalf("expected type ExpressionStmt, got %v", node["type"])
	}

	expr, ok := node["expression"].(map[string]any)
	if !ok {
		t.Fatalf("expected expression object, got %v", node["expression"])
	}

	if typ, ok := expr["type"].(string); !ok || typ != "Binary" {
		t.Fatalf("expected Binary expression, got %v", expr["type"])
	}

	if op, ok := expr["operator"].(string); !ok || op != "+" {
		t.Fatalf("expected operator '+', got %v", expr["operator"])
	}

	if left, ok := expr["left"].(float64); !ok || left != 1 {
		t.Fatalf("expected left 1, got %v", expr["left"])
	}
	if right, ok := expr["right"].(float64); !ok || right != 2 {
		t.Fatalf("expected right 2, got %v", expr["right"])
	}
}

func TestWriteASTJSONToFile(t *testing.T) {
	stmts := []ast.Stmt{
		ast.PrintStmt{Expression: ast.Literal{Value: "hellow nilan!"}},
	}

	filePath := filepath.Join(os.TempDir(), "nilan_ast_printer_test.json")
	defer os.Remove(filePath)

	if err := WriteASTJSONToFile(stmts, filePath); err != nil {
		t.Fatalf("WriteASTJSONToFile error: %v", err)
	}

	bytes, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("read file: %v", err)
	}

	var out []map[string]any
	if err := json.Unmarshal(bytes, &out); err != nil {
		t.Fatalf("unmarshal json: %v", err)
	}

	if len(out) != 1 {
		t.Fatalf("expected 1 statement, got %d", len(out))
	}

	node := out[0]
	if typ, ok := node["type"].(string); !ok || typ != "PrintStmt" {
		t.Fatalf("expected type PrintStmt, got %v", node["type"])
	}

	if expr, ok := node["expression"].(string); !ok || expr != "hellow nilan!" {
		t.Fatalf("expected expression 'hellow nilan!', got %v", node["expression"])
	}
}
