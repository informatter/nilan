package parser

import (
	"encoding/json"
	"fmt"
	"nilan/ast"
	"os"
)

const (
	colorYellow = "\033[33m"
	colorReset  = "\033[0m"
)

// astPrinter implements the Visitor interfaces and builds a
// JSON-friendly representation of the AST using maps and slices.
// Each Visit method returns an object that can be marshaled to JSON.
type astPrinter struct{}

func (p astPrinter) VisitExpressionStmt(exprStmt ast.ExpressionStmt) any {
	return map[string]any{
		"type":       "ExpressionStmt",
		"expression": exprStmt.Expression.Accept(p),
	}
}

func (p astPrinter) VisitPrintStmt(printStmt ast.PrintStmt) any {
	return map[string]any{
		"type":       "PrintStmt",
		"expression": printStmt.Expression.Accept(p),
	}
}

func (p astPrinter) VisitVarStmt(varStmt ast.VarStmt) any {
	return map[string]any{
		"type":        "VarStmt",
		"name":        varStmt.Name.Lexeme,
		"initializer": nilOrAccept(varStmt.Initializer, p),
	}
}

func (p astPrinter) VisitBlockStmt(blockStmt ast.BlockStmt) any {
	stmts := make([]any, 0, len(blockStmt.Statements))
	for _, stmt := range blockStmt.Statements {
		stmts = append(stmts, stmt.Accept(p))
	}
	return map[string]any{
		"type":       "BlockStmt",
		"statements": stmts,
	}
}

func (p astPrinter) VisitWhileStmt(stmt ast.WhileStmt) any {
	return map[string]any{
		"type":      "WhileStmt",
		"condition": stmt.Condition.Accept(p),
		"body":      stmt.Body.Accept(p),
	}
}

func (p astPrinter) VisitIfStmt(stmt ast.IfStmt) any {
	var elseVal any
	if stmt.Else != nil {
		elseVal = stmt.Else.Accept(p)
	} else {
		elseVal = nil
	}
	return map[string]any{
		"type":      "IfStmt",
		"condition": stmt.Condition.Accept(p),
		"then":      stmt.Then.Accept(p),
		"else":      elseVal,
	}
}

func (p astPrinter) VisitLogicalExpression(expr ast.Logical) any {
	return map[string]any{
		"type":     "Logical",
		"operator": expr.Operator.Lexeme,
		"left":     expr.Left.Accept(p),
		"right":    expr.Right.Accept(p),
	}
}

func (p astPrinter) VisitAssignExpression(assign ast.Assign) any {
	return map[string]any{
		"type":  "Assign",
		"name":  assign.Name.Lexeme,
		"value": assign.Value.Accept(p),
	}
}

func (p astPrinter) VisitVariableExpression(variable ast.Variable) any {
	return map[string]any{
		"type": "Variable",
		"name": variable.Name.Lexeme,
	}
}

func (p astPrinter) VisitBinary(b ast.Binary) any {
	return map[string]any{
		"type":     "Binary",
		"operator": b.Operator.Lexeme,
		"left":     b.Left.Accept(p),
		"right":    b.Right.Accept(p),
	}
}

func (p astPrinter) VisitUnary(u ast.Unary) any {
	return map[string]any{
		"type":     "Unary",
		"operator": u.Operator.Lexeme,
		"right":    u.Right.Accept(p),
	}
}

func (p astPrinter) VisitLiteral(l ast.Literal) any {
	// literals are terminal values and can be used directly in JSON
	return l.Value
}

func (p astPrinter) VisitGrouping(g ast.Grouping) any {
	return map[string]any{
		"type":       "Grouping",
		"expression": g.Expression.Accept(p),
	}
}

// nilOrAccept returns nil if expr is nil, otherwise it continues
// processintg the expression and returns the result.
func nilOrAccept(expr ast.Expression, p ast.ExpressionVisitor) any {
	if expr == nil {
		return nil
	}
	return expr.Accept(p)
}

// PrintASTJSON converts a slice of statements into a prettified JSON string.
func PrintASTJSON(statements []ast.Stmt) (string, error) {
	printer := astPrinter{}
	out := make([]any, 0, len(statements))
	for _, s := range statements {
		out = append(out, s.Accept(printer))
	}
	bytes, err := json.MarshalIndent(out, "", "  ")
	if err != nil {
		return "", err
	}

	jsonStr := string(bytes)
	fmt.Println(colorYellow + "----- AST JSON -----")
	fmt.Println(colorYellow + jsonStr)
	fmt.Println(colorYellow + "-----" + colorReset)
	fmt.Println("")
	return jsonStr, nil
}

// WriteASTJSONToFile writes the prettified AST JSON to the given file path.
func WriteASTJSONToFile(statements []ast.Stmt, path string) error {
	s, err := PrintASTJSON(statements)
	if err != nil {
		return err
	}
	fDescriptor, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("error creating AST file: %s", err.Error())
	}

	_, error := fDescriptor.Write([]byte(s))
	if error != nil {
		return fmt.Errorf("error writing AST to file: %s", error.Error())
	}
	defer fDescriptor.Close()
	return nil
}
