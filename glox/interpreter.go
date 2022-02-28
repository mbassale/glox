package glox

import (
	"fmt"
	"reflect"
	"strconv"
)

const INTERPRETER_WHERE = "interpreter"

type Interpreter struct {
	errorReporter ErrorReporter
	environment   Environment
	lastValue     interface{}
}

func NewInterpreter(errorReporter ErrorReporter) Interpreter {
	return Interpreter{
		errorReporter: errorReporter,
		environment:   NewEnvironment(),
		lastValue:     nil,
	}
}

func (inter *Interpreter) Interpret(statements []Stmt) (interface{}, error) {
	for _, stmt := range statements {
		inter.execute(stmt)
	}
	return inter.lastValue, nil
}

func (inter *Interpreter) GetLastValue() interface{} {
	return inter.lastValue
}

func (inter *Interpreter) execute(stmt Stmt) {
	stmt.accept(inter)
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

func (inter *Interpreter) visitBlockStmt(stmt BlockStmt) interface{} {
	inter.executeBlock(stmt.Statements, NewEnvironmentWithEnclosing(&inter.environment))
	return nil
}

func (inter *Interpreter) visitExpressionStmt(stmt ExpressionStmt) interface{} {
	inter.lastValue = inter.evaluate(stmt.Expression)
	return inter.lastValue
}

func (inter *Interpreter) visitPrintStmt(stmt PrintStmt) interface{} {
	value := inter.evaluate(stmt.Print)
	fmt.Println(value)
	inter.lastValue = value
	return inter.lastValue
}

func (inter *Interpreter) visitVarStmt(stmt VarStmt) interface{} {
	var value interface{} = nil
	if stmt.Initializer != nil {
		value = inter.evaluate(stmt.Initializer)
	}
	inter.environment.Define(stmt.Name.Lexeme, value)
	inter.lastValue = value
	return value
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
		val, err := anyToFloat64(right)
		if err != nil {
			inter.errorReporter.Error(expr.Operator.Line, err.Error())
			return nil
		}
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
		leftVal, rightVal, err := inter.checkNumberOperands(expr, expr.Operator, left, right)
		if err != nil {
			return nil
		}
		return leftVal > rightVal
	case TOKEN_GREATER_EQUAL:
		leftVal, rightVal, err := inter.checkNumberOperands(expr, expr.Operator, left, right)
		if err != nil {
			return nil
		}
		return leftVal >= rightVal
	case TOKEN_LESS:
		leftVal, rightVal, err := inter.checkNumberOperands(expr, expr.Operator, left, right)
		if err != nil {
			return nil
		}
		return leftVal < rightVal
	case TOKEN_LESS_EQUAL:
		leftVal, rightVal, err := inter.checkNumberOperands(expr, expr.Operator, left, right)
		if err != nil {
			return nil
		}
		return leftVal <= rightVal
	case TOKEN_MINUS:
		leftVal, rightVal, err := inter.checkNumberOperands(expr, expr.Operator, left, right)
		if err != nil {
			return nil
		}
		return leftVal - rightVal
	case TOKEN_SLASH:
		leftVal, rightVal, err := inter.checkNumberOperands(expr, expr.Operator, left, right)
		if err != nil {
			return nil
		}
		return leftVal / rightVal
	case TOKEN_STAR:
		leftVal, rightVal, err := inter.checkNumberOperands(expr, expr.Operator, left, right)
		if err != nil {
			return nil
		}
		return leftVal * rightVal
	case TOKEN_PLUS:
		if isNumber(left) && isNumber(right) {
			leftVal, rightVal, err := inter.checkNumberOperands(expr, expr.Operator, left, right)
			if err != nil {
				return nil
			}
			return leftVal + rightVal
		} else if isString(left) && isString(right) {
			leftVal, rightVal, err := inter.checkStringOperands(expr, expr.Operator, left, right)
			if err != nil {
				return nil
			}
			return leftVal + rightVal
		}
		inter.errorReporter.Push(expr.getLine(), INTERPRETER_WHERE, fmt.Errorf("operands must be two numbers or two strings"))
		return nil
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

func (inter *Interpreter) visitVariableExpr(expr VariableExpr) interface{} {
	value, err := inter.environment.Get(expr.Name.Lexeme)
	if err != nil {
		inter.errorReporter.Push(expr.getLine(), INTERPRETER_WHERE, err)
		return nil
	}
	return value
}

func (inter *Interpreter) visitAssignExpr(expr AssignExpr) interface{} {
	value := inter.evaluate(expr.Value)
	inter.environment.Assign(expr.Name.Lexeme, value)
	return value
}

func (inter *Interpreter) evaluate(expr Expr) interface{} {
	return expr.accept(inter)
}

func (inter *Interpreter) executeBlock(statements []Stmt, localEnv Environment) {
	previousEnv := inter.environment
	inter.environment = localEnv
	for _, stmt := range statements {
		inter.execute(stmt)
	}
	inter.environment = previousEnv
}

func (inter *Interpreter) checkNumberOperands(expr Expr, operator Token, left interface{}, right interface{}) (float64, float64, error) {
	returnError := func(err error) (float64, float64, error) {
		inter.errorReporter.Push(expr.getLine(), INTERPRETER_WHERE, fmt.Errorf("operator %s: operands must be numbers: %w", operator.Lexeme, err))
		return 0, 0, nil
	}
	leftVal, err := anyToFloat64(left)
	if err != nil {
		return returnError(err)
	}
	rightVal, err := anyToFloat64(right)
	if err != nil {
		return returnError(err)
	}
	return leftVal, rightVal, nil
}

func (inter *Interpreter) checkStringOperands(expr Expr, operator Token, left interface{}, right interface{}) (string, string, error) {
	returnError := func(err error) (string, string, error) {
		inter.errorReporter.Push(expr.getLine(), INTERPRETER_WHERE, fmt.Errorf("operator %s: operands must be strings: %w", operator.Lexeme, err))
		return "", "", nil
	}
	leftVal, err := anyToString(left)
	if err != nil {
		return returnError(err)
	}
	rightVal, err := anyToString(right)
	if err != nil {
		return returnError(err)
	}
	return leftVal, rightVal, nil
}
