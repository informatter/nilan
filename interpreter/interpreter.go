package interpreter

import (
	"fmt"
	"nilan/parser"
	"nilan/token"
	"strconv"
)

type Interpreter struct{}

func (i Interpreter) Interpret(expression parser.Expression) any {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println(r)
		}
	}()
	return i.evaluate(expression)
}

// Evaluates a Binary expression
//
// Parameters:
//   - binary: parser.Binary
//     The Binary expression to evaluate
//
// Returns:
//   - any: The result from the binary expression.
func (i Interpreter) VisitBinary(binary parser.Binary) any {
	leftResult := i.evaluate(binary.Left)
	rightResult := i.evaluate(binary.Right)
	operator := binary.Operator.TokenType
	switch operator {
	case token.MULT:
		// TODO if string is being multipled, duplicate string by N
		// This will allow string multiplication like Python which can be
		// interesting?

		leftValue, rightValue, err := isOperandsNumeric(operator, leftResult, rightResult, binary.Operator)

		if err != nil {
			panic(err.Error())
		}
		return leftValue * rightValue
	case token.DIV:
		leftValue, rightValue, err := isOperandsNumeric(operator, leftResult, rightResult, binary.Operator)

		if err != nil {
			panic(err.Error())
		}
		return leftValue / rightValue

	case token.SUB:
		// TODO: There seems to be a weird bug in the lexers handleNumber method
		// 2-2 yields an error while 2 - 2 yields 0
		leftValue, rightValue, err := isOperandsNumeric(operator, leftResult, rightResult, binary.Operator)

		if err != nil {
			panic(err.Error())
		}
		return leftValue - rightValue

	case token.ADD:

		leftValue, rightValue, err := isOperandsNumeric(operator, leftResult, rightResult, binary.Operator)
		if err != nil {
			// Both operands are not numeric check if both are strings.
			leftValString, ok := leftResult.(string)
			rightValString, okk := rightResult.(string)
			if ok && okk {
				// Make sure none of them are numeric. If one of them is panic on the error
				// raised by `isOperandsNumeric`
				_, errA := strconv.ParseFloat(leftValString, 64)
				_, errB := strconv.ParseFloat(rightValString, 64)
				if errA == nil || errB == nil {
					panic(err.Error())
				}
				return leftValString + rightValString

			}
		}
		return leftValue + rightValue
	
	// NOTE: For now, we can use the same equality comparison as Go's
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
		message := fmt.Sprintf("operator '%s' not supported",operator)
		error := CreateRuntimeError(binary.Operator.Line, binary.Operator.Column, message)
		panic(error)
	}
}

func (i Interpreter) VisitUnary(unary parser.Unary) any {
	return nil
}

func (i Interpreter) VisitLiteral(literal parser.Literal) any {
	return literal.Value
}

func (i Interpreter) VisitGrouping(grouping parser.Grouping) any {
	return i.evaluate(grouping.Expression)
}

func (i Interpreter) evaluate(expression parser.Expression) any {
	return expression.Accept(Interpreter{})
}

// Converts the value of a Literal expression to a float64.
// Currently the value of a Literal expression is a string.
func literalToFloat64(value interface{}) (float64, error) {

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


func isOperandsNumeric(operator token.TokenType, left any, right any, token token.Token) (float64, float64, error) {

	l, lerr := literalToFloat64(left)
	r, rerr := literalToFloat64(right)
	if lerr == nil && rerr == nil {
		return l, r, nil
	}
	message := fmt.Sprintf("operands must be numeric values. '%s %s %s' is not allowed", left, operator, right)
	error := CreateRuntimeError(token.Line, token.Column, message)
	return 0, 0, error
}
