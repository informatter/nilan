package interpreter

import (
	"fmt"
	"nilan/parser"
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
func (i *TreeWalkInterpreter) Interpret(statements []parser.Stmt) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println(r)
		}
	}()
	i.executeStatements(statements)
}

// executeStatements executes each statement by invoking its Accept method
// with a fresh Interpreter visitor. It does not return anything.
func (i *TreeWalkInterpreter) executeStatements(statements []parser.Stmt) {
	for _, s := range statements {
		s.Accept(i)
	}
}

// VisitExpressionStmt visits an ExpressionStmt node.
// Evaluates the expression but does not return a value.
//
// Returns:
//   - any: always nil because statements do not produce values.
func (i *TreeWalkInterpreter) VisitExpressionStmt(exprStatement parser.ExpressionStmt) any {
	i.evaluate(exprStatement.Expression)
	return nil
}

// VisitPrintStmt visits a PrintStmt node.
// Evaluates the expression and prints the result.
//
// Returns:
//   - any: always nil because print statements have no return value.
func (i *TreeWalkInterpreter) VisitPrintStmt(printStmt parser.PrintStmt) any {
	value := i.evaluate(printStmt.Expression)
	fmt.Println(value)
	return nil
}

// VisitVarStmt visits a VarStmt node.
// It evaluates the initialiser expression of the statement if it contains one
// and it sets the name of the variable to its evaluated value.
// Returns:
//  - nil: This method returns nil, as it mutates its own state to store
//    a variable name to its value
func (i *TreeWalkInterpreter) VisitVarStmt(varStmt parser.VarStmt) any {
	var value any = nil
	if varStmt.Initializer != nil {
		value = i.evaluate(varStmt.Initializer)
	}
	i.environment.set(varStmt.Name.Lexeme, value)
	return nil
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
func (i *TreeWalkInterpreter) VisitBinary(binary parser.Binary) any {
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
func (i *TreeWalkInterpreter) VisitUnary(unary parser.Unary) any {
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

// Retrieves the value for variable.
// Returns:
//  - The value of the variable
// Raises:
//  - RuntimeError: panics with a RuntimeError if attempting to access an undefined
//    variable
func (i *TreeWalkInterpreter) VisitVariableExpression(expression parser.Variable) any {
	value, error := i.environment.get(expression.Name)
	if error != nil {
		panic(error)
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
func (i *TreeWalkInterpreter) VisitLiteral(literal parser.Literal) any {
	return literal.Value
}

// VisitGrouping evaluates a Grouping expression by evaluating its inner expression.
//
// Parameters:
//   - grouping: the parser.Grouping node.
//
// Returns:
//   - any: the value of the enclosed expression.
func (i *TreeWalkInterpreter) VisitGrouping(grouping parser.Grouping) any {
	return i.evaluate(grouping.Expression)
}

// evaluate evaluates any expression node by invoking its Accept method
// with the Interpreter visitor.
//
// Returns:
//   - any: the evaluated value of the expression.
func (i *TreeWalkInterpreter) evaluate(expression parser.Expression) any {
	// return expression.Accept(TreeWalkInterpreter{})
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
