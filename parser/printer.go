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

type blockStmtJSON struct {
	Type       string `json:"type"`
	Statements []any  `json:"statements"`
}

type expressionStmtJSON struct {
	Type       string `json:"type"`
	Expression any    `json:"expression"`
}

type printStmtJSON struct {
	Type       string `json:"type"`
	Expression any    `json:"expression"`
}

type varStmtJSON struct {
	Type        string `json:"type"`
	Name        string `json:"name"`
	Initializer any    `json:"initializer"`
}

type whileStmtJSON struct {
	Type      string `json:"type"`
	Condition any    `json:"condition"`
	Body      any    `json:"body"`
}

type ifStmtJSON struct {
	Type      string `json:"type"`
	Condition any    `json:"condition"`
	Then      any    `json:"then"`
	Else      any    `json:"else"`
}

type assignExprJSON struct {
	Type  string `json:"type"`
	Name  string `json:"name"`
	Value any    `json:"value"`
}

type logicalExprJSON struct {
	Type     string `json:"type"`
	Operator string `json:"operator"`
	Left     any    `json:"left"`
	Right    any    `json:"right"`
}

type groupingExprJSON struct {
	Type       string `json:"type"`
	Expression any    `json:"expression"`
}

type variableExprJSON struct {
	Type string `json:"type"`
	Name string `json:"name"`
}

type binaryExprJSON struct {
	Type     string `json:"type"`
	Operator string `json:"operator"`
	Left     any    `json:"left"`
	Right    any    `json:"right"`
}

type unaryExprJSON struct {
	Type     string `json:"type"`
	Operator string `json:"operator"`
	Right    any    `json:"right"`
}

// astPrinter implements the Visitor interfaces and builds a
// JSON-friendly representation of the AST using maps and slices.
// Each Visit method returns an object that can be marshaled to JSON.
type astPrinter struct{}

func (p astPrinter) VisitExpressionStmt(exprStmt ast.ExpressionStmt) any {
	return expressionStmtJSON{
		Type:       "ExpressionStmt",
		Expression: exprStmt.Expression.Accept(p),
	}
}

func (p astPrinter) VisitPrintStmt(printStmt ast.PrintStmt) any {
	return printStmtJSON{
		Type:       "PrintStmt",
		Expression: printStmt.Expression.Accept(p),
	}
}

func (p astPrinter) VisitVarStmt(varStmt ast.VarStmt) any {
	return varStmtJSON{
		Type:        "VarStmt",
		Name:        varStmt.Name.Lexeme,
		Initializer: nilOrAccept(varStmt.Initializer, p),
	}
}

func (p astPrinter) VisitBlockStmt(blockStmt ast.BlockStmt) any {
	stmts := make([]any, 0, len(blockStmt.Statements))
	for _, stmt := range blockStmt.Statements {
		stmts = append(stmts, stmt.Accept(p))
	}
	return blockStmtJSON{
		Type:       "BlockStmt",
		Statements: stmts,
	}
}

func (p astPrinter) VisitWhileStmt(stmt ast.WhileStmt) any {
	return whileStmtJSON{
		Type:      "WhileStmt",
		Condition: stmt.Condition.Accept(p),
		Body:      stmt.Body.Accept(p),
	}
}

func (p astPrinter) VisitIfStmt(stmt ast.IfStmt) any {
	var elseVal any
	if stmt.Else != nil {
		elseVal = stmt.Else.Accept(p)
	} else {
		elseVal = nil
	}
	return ifStmtJSON{
		Type:      "IfStmt",
		Condition: stmt.Condition.Accept(p),
		Then:      stmt.Then.Accept(p),
		Else:      elseVal,
	}
}

func (p astPrinter) VisitLogicalExpression(expr ast.Logical) any {
	return logicalExprJSON{
		Type:     "Logical",
		Operator: expr.Operator.Lexeme,
		Left:     expr.Left.Accept(p),
		Right:    expr.Right.Accept(p),
	}
}

func (p astPrinter) VisitAssignExpression(assign ast.Assign) any {
	return assignExprJSON{
		Type:  "Assign",
		Name:  assign.Name.Lexeme,
		Value: assign.Value.Accept(p),
	}
}

func (p astPrinter) VisitVariableExpression(variable ast.Variable) any {
	return variableExprJSON{
		Type: "Variable",
		Name: variable.Name.Lexeme,
	}
}

func (p astPrinter) VisitBinary(b ast.Binary) any {
	return binaryExprJSON{
		Type:     "Binary",
		Operator: b.Operator.Lexeme,
		Left:     b.Left.Accept(p),
		Right:    b.Right.Accept(p),
	}
}

func (p astPrinter) VisitUnary(u ast.Unary) any {
	return unaryExprJSON{
		Type:     "Unary",
		Operator: u.Operator.Lexeme,
		Right:    u.Right.Accept(p),
	}
}

func (p astPrinter) VisitLiteral(l ast.Literal) any {
	// literals are terminal values and can be used directly in JSON
	return l.Value
}

func (p astPrinter) VisitGrouping(g ast.Grouping) any {
	return groupingExprJSON{
		Type:       "Grouping",
		Expression: g.Expression.Accept(p),
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
