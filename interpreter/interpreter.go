package interpreter

import (
	"fmt"
	"nilan/ast"
	"nilan/token"
	"strconv"
)

// TreeWalkInterpreter executes parsed statements and evaluates expressions.
type TreeWalkInterpreter struct {
	environment *Environment
}

// Creates an instance of a "Tree-Walk Interpreter"
func Make() *TreeWalkInterpreter {
	return &TreeWalkInterpreter{
		environment: MakeEnvironment(),
	}
}

// Interpret executes a list of statements.
// It recovers from panics to print runtime errors without crashing.
func (i *TreeWalkInterpreter) Interpret(statements []ast.Stmt) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println(r)
		}
	}()
	i.executeStatements(statements)
}

// executeStatements executes each statement by invoking its Accept method.
func (i *TreeWalkInterpreter) executeStatements(statements []ast.Stmt) {
	for _, s := range statements {
		s.Accept(i)
	}
}

// executeStmt executes the given AST node statement by invoking its Accept method,
// which calls the appropriate Visit method of the interpreter.
func (i *TreeWalkInterpreter) executeStmt(stmt ast.Stmt) {

	// Implements the visitor pattern to process different
	// kinds of statements polymorphically.
	stmt.Accept(i)
}

// VisitBlockStmt executes all statements in the given ast.BlockStmt
// within a new nested environment. It temporarily replaces the current
// interpreter environment with a new one scoped as a child of the previous environment.
// A deferred function ensures that if a panic occurs, the environment
// is restored and the panic is printed. After executing the statements,
// the previous environment is always restored,
// providing block-scoped execution and panic safety.
func (i *TreeWalkInterpreter) VisitBlockStmt(blockStmt ast.BlockStmt) any {

	previous := i.environment
	i.environment = MakeNestedEnvironment(i.environment)
	defer func() {
		if r := recover(); r != nil {
			i.environment = previous
			fmt.Println(r)
		}
	}()

	i.executeStatements(blockStmt.Statements)
	i.environment = previous
	return nil
}

// VisitExpressionStmt visits an ExpressionStmt node.
// Evaluates the expression but does not return a value.
//
// Returns:
//   - any: always nil because statements do not produce values.
func (i *TreeWalkInterpreter) VisitExpressionStmt(exprStatement ast.ExpressionStmt) any {
	i.evaluate(exprStatement.Expression)
	return nil
}

// VisitIfStmt evaluates the condition of the given ast.IfStmt.
// If the condition evaluates to true (according to interpreter semantics),
// it executes the 'Then' branch.
// If an 'Else' branch is present and if the condition is false, it
// is executed.

// Returns:
//   - any: always nil because statements do not produce values.
func (i *TreeWalkInterpreter) VisitIfStmt(stmt ast.IfStmt) any {
	if i.isTrue(i.evaluate(stmt.Condition)) {
		i.executeStmt(stmt.Then)
	} else if stmt.Else != nil {
		i.executeStmt(stmt.Else)
	}

	return nil
}

// VisitPrintStmt visits a PrintStmt node.
// Evaluates the expression and prints the result.
//
// Returns:
//   - any: always nil because print statements have no return value.
func (i *TreeWalkInterpreter) VisitPrintStmt(printStmt ast.PrintStmt) any {
	value := i.evaluate(printStmt.Expression)
	if value == nil {
		fmt.Println("null")
		return nil
	}
	fmt.Println(value)
	return nil
}

// VisitVarStmt visits a VarStmt node.
// It evaluates the initialiser expression of the statement if it contains one
// and it sets the name of the variable to its evaluated value.
// Returns:
//   - nil: This method returns nil, as it mutates its own state to store
//     a variable name to its value
func (i *TreeWalkInterpreter) VisitVarStmt(varStmt ast.VarStmt) any {
	var value any = nil
	if varStmt.Initializer != nil {
		value = i.evaluate(varStmt.Initializer)
	}
	i.environment.set(varStmt.Name.Lexeme, value)
	return nil
}

// VisitAssignExpression evaluates an assignment expression node and updates
// the value of the corresponding variable in the environment.
//
// Steps:
//  1. The right-hand side expression (`assign.Value`) is evaluated using
//     the interpreter's `evaluate` method.
//  2. The resulting value is attempted to be assigned to the variable
//     identified by `assign.Name` via the environment's `assign` method.
//  3. If the variable is undefined in the current environment, a runtime
//     error is returned by `assign`.
//  4. On success, the new value is returned.
//
// Parameters:
//   - assign: An assignment AST node containing the variable token (Name)
//     and the expression to evaluate (Value).
//
// Returns:
//   - any: The value resulting from evaluating `assign.Value`, which is
//     also the value bound to the variable after the assignment.
func (i *TreeWalkInterpreter) VisitAssignExpression(assign ast.Assign) any {
	value := i.evaluate(assign.Value)
	err := i.environment.assign(assign.Name, value)
	if err != nil {
		panic(err.Error())
	}
	return value
}

// VisitBinary evaluates a binary expression node.
//
// Parameters:
//   - binary: the parser.Binary expression node.
//
// Returns:
//   - any: evaluated result of the binary expression (number, string, bool).
//
// Panics on invalid operands or unsupported operators.
func (i *TreeWalkInterpreter) VisitBinary(binary ast.Binary) any {
	leftResult := i.evaluate(binary.Left)
	rightResult := i.evaluate(binary.Right)
	operator := binary.Operator.TokenType

	switch operator {
	case token.MULT:
		leftValue, rightValue, err := isOperandsNumeric(operator, leftResult, rightResult, binary.Operator)
		if err != nil {
			panic(err.Error())
		}
		// TODO: support string multiplication by integer count
		return leftValue * rightValue

	case token.DIV:
		leftValue, rightValue, err := isOperandsNumeric(operator, leftResult, rightResult, binary.Operator)
		if err != nil {
			panic(err.Error())
		}
		if rightValue == 0 {
			return CreateRuntimeError(binary.Operator.Line, binary.Operator.Column+1, "Division by zero")
		}
		return leftValue / rightValue

	case token.SUB:
		leftValue, rightValue, err := isOperandsNumeric(operator, leftResult, rightResult, binary.Operator)
		if err != nil {
			panic(err.Error())
		}
		return leftValue - rightValue

	case token.ADD:
		leftValue, rightValue, err := isOperandsNumeric(operator, leftResult, rightResult, binary.Operator)
		if err != nil {
			// If not numeric, check if both are strings for concatenation
			leftValString, ok := leftResult.(string)
			rightValString, okk := rightResult.(string)
			if ok && okk {
				// Verify neither string parses as number
				_, errA := strconv.ParseFloat(leftValString, 64)
				_, errB := strconv.ParseFloat(rightValString, 64)
				if errA == nil || errB == nil {
					panic(err.Error())
				}
				return leftValString + rightValString
			}
			// Otherwise propagate the error
			panic(err.Error())
		}
		return leftValue + rightValue

	case token.EQUAL_EQUAL:
		return leftResult == rightResult

	case token.NOT_EQUAL:
		return leftResult != rightResult

	case token.LARGER:
		leftValue, rightValue, err := isOperandsNumeric(operator, leftResult, rightResult, binary.Operator)
		if err != nil {
			panic(err)
		}
		return leftValue > rightValue

	case token.LARGER_EQUAL:
		leftValue, rightValue, err := isOperandsNumeric(operator, leftResult, rightResult, binary.Operator)
		if err != nil {
			panic(err)
		}
		return leftValue >= rightValue

	case token.LESS:
		leftValue, rightValue, err := isOperandsNumeric(operator, leftResult, rightResult, binary.Operator)
		if err != nil {
			panic(err)
		}
		return leftValue < rightValue

	case token.LESS_EQUAL:
		leftValue, rightValue, err := isOperandsNumeric(operator, leftResult, rightResult, binary.Operator)
		if err != nil {
			panic(err)
		}
		return leftValue <= rightValue

	default:
		message := fmt.Sprintf("operator '%s' not supported", operator)
		error := CreateRuntimeError(binary.Operator.Line, binary.Operator.Column, message)
		panic(error)
	}
}

// VisitUnary evaluates a unary expression node.
//
// Parameters:
//   - unary: the parser.Unary expression node.
//
// Returns:
//   - any: the evaluated result of the unary operation.
//
// Panics on invalid operand types or unsupported operators.
func (i *TreeWalkInterpreter) VisitUnary(unary ast.Unary) any {
	rightResult := i.evaluate(unary.Right)
	operator := unary.Operator.TokenType
	switch operator {
	case token.SUB:
		r, err := literalToFloat64(rightResult)
		if err != nil {
			message := fmt.Sprintf("operand must be a numeric value. '%s %s' is not allowed", operator, rightResult)
			error := CreateRuntimeError(unary.Operator.Line, unary.Operator.Column, message)
			panic(error)
		}
		return -r
	case token.BANG:
		if rightResult == nil {
			return true
		}
		value, isBool := rightResult.(bool)
		if isBool {
			return !value
		}
		return false
	default:
		message := fmt.Sprintf("operator '%s' not supported for unary operations", operator)
		error := CreateRuntimeError(unary.Operator.Line, unary.Operator.Column, message)
		panic(error)
	}
}

// isTrue determines the "truthiness" of the given object according to interpreter rules.
// It returns false if the object is nil. If the object is explicitly a bool,
// it returns the boolean value. All other non-nil values are treated as true.
func (i *TreeWalkInterpreter) isTrue(object any) bool {
	if object == nil {
		return false
	}
	value, isBool := object.(bool)
	if isBool {
		return value
	}
	return true
}

// Retrieves the value for variable.
// Returns:
//   - The value of the variable
//
// Raises:
//   - RuntimeError: panics with a RuntimeError if attempting to access an undefined
//     variable
func (i *TreeWalkInterpreter) VisitVariableExpression(expression ast.Variable) any {
	value, error := i.environment.get(expression.Name)
	if error != nil {
		panic(error)
	}
	if value == nil {
		msg := fmt.Sprintf("Cant access uninitialised variable: %s", expression.Name.Lexeme)
		err := CreateRuntimeError(expression.Name.Line, expression.Name.Column, msg)
		panic(err)
	}
	return value
}

// VisitLiteral returns the value of a Literal node.
//
// Parameters:
//   - literal: the parser.Literal node.
//
// Returns:
//   - any: the literal's underlying value.
func (i *TreeWalkInterpreter) VisitLiteral(literal ast.Literal) any {
	return literal.Value
}

// VisitGrouping evaluates a Grouping expression by evaluating its inner expression.
//
// Parameters:
//   - grouping: the parser.Grouping node.
//
// Returns:
//   - any: the value of the enclosed expression.
func (i *TreeWalkInterpreter) VisitGrouping(grouping ast.Grouping) any {
	return i.evaluate(grouping.Expression)
}

// evaluate evaluates any expression node by invoking its Accept method
// with the Interpreter visitor.
//
// Returns:
//   - any: the evaluated value of the expression.
func (i *TreeWalkInterpreter) evaluate(expression ast.Expression) any {
	return expression.Accept(i)
}

// literalToFloat64 attempts to convert a literal value into a float64.
//
// Parameters:
//   - value: the literal value (various possible types).
//
// Returns:
//   - float64: the converted numeric value.
//   - error: on failure to convert value to float64.
func literalToFloat64(value any) (float64, error) {
	switch v := value.(type) {
	case float64:
		return v, nil
	case float32:
		return float64(v), nil
	case int:
		return float64(v), nil
	case int16:
		return float64(v), nil
	case int32:
		return float64(v), nil
	case int64:
		return float64(v), nil
	case string:
		result, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return 0, err
		}
		return result, nil
	default:
		return 0, fmt.Errorf("unsupported type: %T", value)
	}
}

// isOperandsNumeric validates that both operands are numeric and converts them to float64.
//
// Parameters:
//   - operator: the token type of the operator.
//   - left, right: values of the operands.
//   - token: token for error positioning.
//
// Returns:
//   - float64: numeric value of left operand.
//   - float64: numeric value of right operand.
//   - error: if either operand cannot be converted to float64.
func isOperandsNumeric(operator token.TokenType, left any, right any, token token.Token) (float64, float64, error) {
	l, lerr := literalToFloat64(left)
	r, rerr := literalToFloat64(right)

	if lerr == nil && rerr == nil {
		return l, r, nil
	}

	message := fmt.Sprintf("operands must be numeric values. '%v %s %v' is not allowed", left, operator, right)
	error := CreateRuntimeError(token.Line, token.Column, message)
	return 0, 0, error
}
