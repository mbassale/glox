package glox

import (
	"fmt"
	"reflect"
	"strconv"
)

type Interpreter struct {
}

func NewInterpreter() Interpreter {
	return Interpreter{}
}

func (inter *Interpreter) Interpret(expr Expr) interface{} {
	return inter.evaluate(expr)
}

func isString(val interface{}) bool {
	switch val.(type) {
	case string:
		return true
	default:
		return false
	}
}

func isNumber(val interface{}) bool {
	switch val.(type) {
	case float64:
		return true
	default:
		return false
	}
}

func anyToString(val interface{}) (string, error) {
	switch val := val.(type) {
	case string:
		return val, nil
	case nil:
		return "nil", nil
	case float64:
		return strconv.FormatFloat(val, 'f', -1, 64), nil
	default:
		return "", fmt.Errorf("cannot convert to string: %v", val)
	}
}

func anyToFloat64(val interface{}) (float64, error) {
	switch val := val.(type) {
	case float64:
		return val, nil
	case string:
		return strconv.ParseFloat(val, 64)
	case float32:
		return float64(val), nil
	default:
		return 0, fmt.Errorf("cannot convert to float: %v", val)
	}
}

func isTruthy(val interface{}) (bool, error) {
	switch val := val.(type) {
	case bool:
		return val, nil
	case nil:
		return false, nil
	case string:
		return len(val) > 0, nil
	case float64:
		return val > 0, nil
	default:
		return false, fmt.Errorf("cannot convert to boolean: %v", val)
	}
}

func isEqual(left interface{}, right interface{}) bool {
	return reflect.DeepEqual(left, right)
}

func (inter *Interpreter) visitLiteralExpr(expr LiteralExpr) interface{} {
	return expr.Value
}

func (inter *Interpreter) visitGroupingExpr(expr GroupingExpr) interface{} {
	return inter.evaluate(expr.Expression)
}

func (inter *Interpreter) visitUnaryExpr(expr UnaryExpr) interface{} {
	right := inter.evaluate(expr.Right)

	switch expr.Operator.Type {
	case TOKEN_BANG:
		val, _ := isTruthy(right)
		return val
	case TOKEN_MINUS:
		val, _ := anyToFloat64(right)
		return val
	}

	// unreachable
	return nil
}

func (inter *Interpreter) visitBinaryExpr(expr BinaryExpr) interface{} {
	left := inter.evaluate(expr.Left)
	right := inter.evaluate(expr.Right)

	switch expr.Operator.Type {
	case TOKEN_GREATER:
		leftVal, _ := anyToFloat64(left)
		rightVal, _ := anyToFloat64(right)
		return leftVal > rightVal
	case TOKEN_GREATER_EQUAL:
		leftVal, _ := anyToFloat64(left)
		rightVal, _ := anyToFloat64(right)
		return leftVal >= rightVal
	case TOKEN_LESS:
		leftVal, _ := anyToFloat64(left)
		rightVal, _ := anyToFloat64(right)
		return leftVal < rightVal
	case TOKEN_LESS_EQUAL:
		leftVal, _ := anyToFloat64(left)
		rightVal, _ := anyToFloat64(right)
		return leftVal <= rightVal
	case TOKEN_MINUS:
		leftVal, _ := anyToFloat64(left)
		rightVal, _ := anyToFloat64(right)
		return leftVal - rightVal
	case TOKEN_SLASH:
		leftVal, _ := anyToFloat64(left)
		rightVal, _ := anyToFloat64(right)
		return leftVal / rightVal
	case TOKEN_STAR:
		leftVal, _ := anyToFloat64(left)
		rightVal, _ := anyToFloat64(right)
		return leftVal * rightVal
	case TOKEN_PLUS:
		if isNumber(left) && isNumber(right) {
			leftVal, _ := anyToFloat64(left)
			rightVal, _ := anyToFloat64(right)
			return leftVal + rightVal
		}
		if isString(left) && isString(right) {
			leftVal, _ := anyToString(left)
			rightVal, _ := anyToString(right)
			return leftVal + rightVal
		}
	case TOKEN_BANG_EQUAL:
		return !isEqual(left, right)
	case TOKEN_EQUAL_EQUAL:
		return isEqual(left, right)
	}

	// unreachable
	return nil
}

func (inter *Interpreter) visitConditionalExpr(expr ConditionalExpr) interface{} {
	condition := inter.evaluate(expr.Condition)
	if val, _ := isTruthy(condition); val {
		return inter.evaluate(expr.Left)
	} else {
		return inter.evaluate(expr.Right)
	}
}

func (inter *Interpreter) evaluate(expr Expr) interface{} {
	return expr.accept(inter)
}
