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

type funcStmtJSON struct {
	Type   string   `json:"type"`
	Name   string   `json:"name"`
	Params []string `json:"params"`
	Body   []any    `json:"body"`
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

type callExprJSON struct {
	Type      string `json:"type"`
	Callee    any    `json:"callee"`
	Arguments []any  `json:"arguments"`
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

func (p astPrinter) VisitLogicalExpr(expr ast.LogicalExpr) any {
	return logicalExprJSON{
		Type:     "Logical",
		Operator: expr.Operator.Lexeme,
		Left:     expr.Left.Accept(p),
		Right:    expr.Right.Accept(p),
	}
}

func (p astPrinter) VisitAssignExpr(assign ast.AssignExpr) any {
	return assignExprJSON{
		Type:  "Assign",
		Name:  assign.Name.Lexeme,
		Value: assign.Value.Accept(p),
	}
}

func (p astPrinter) VisitVariableExpr(variable ast.VariableExpr) any {
	return variableExprJSON{
		Type: "Variable",
		Name: variable.Name.Lexeme,
	}
}

func (p astPrinter) VisitBinaryExpr(b ast.BinaryExpr) any {
	return binaryExprJSON{
		Type:     "Binary",
		Operator: b.Operator.Lexeme,
		Left:     b.Left.Accept(p),
		Right:    b.Right.Accept(p),
	}
}

func (p astPrinter) VisitUnaryExpr(u ast.UnaryExpr) any {
	return unaryExprJSON{
		Type:     "Unary",
		Operator: u.Operator.Lexeme,
		Right:    u.Right.Accept(p),
	}
}

func (p astPrinter) VisitLiteralExpr(l ast.LiteralExpr) any {
	// literals are terminal values and can be used directly in JSON
	return l.Value
}

func (p astPrinter) VisitGroupingExpr(g ast.GroupingExpr) any {
	return groupingExprJSON{
		Type:       "Grouping",
		Expression: g.Expression.Accept(p),
	}
}

func (p astPrinter) VisitCallExpr(call ast.CallExpr) any {

	args := make([]any, 0, len(call.Arguments))
	for _, arg := range call.Arguments {
		args = append(args, arg.Accept(p))
	}

	return callExprJSON{
		Type:      "Call",
		Callee:    call.Callee.Accept(p),
		Arguments: args,
	}

}

func (p astPrinter) VisitFunctionStmt(functionStmt ast.FunctionStmt) any {
	params := make([]string, 0, len(functionStmt.Parameters))
	for _, param := range functionStmt.Parameters {
		params = append(params, param.Lexeme)
	}
	body := make([]any, 0, len(functionStmt.Body.Statements))
	for _, stmt := range functionStmt.Body.Statements {
		body = append(body, stmt.Accept(p))
	}
	return funcStmtJSON{
		Type:   "FunctionStmt",
		Name:   functionStmt.Name.Lexeme,
		Params: params,
		Body:   body,
	}
}

func (p astPrinter) VisitReturnStmt(returnStmt ast.ReturnStmt) any {

	var value any
	if returnStmt.Value != nil {
		value = returnStmt.Value.Accept(p)
	} else {
		value = nil
	}
	return map[string]any{
		"type":  "ReturnStmt",
		"value": value,
	}
}

// nilOrAccept returns nil if expr is nil, otherwise it continues
// processintg the expression and returns the result.
func nilOrAccept(expr ast.Expr, p ast.ExprVisitor) any {
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
