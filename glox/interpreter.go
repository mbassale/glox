package glox

import (
	"fmt"
	"reflect"
	"strconv"
)

const INTERPRETER_WHERE = "interpreter"

type Interpreter struct {
	errorReporter ErrorReporter
	environment   *Environment
	globals       *Environment
	lastValue     interface{}
}

type BreakResult struct {
}

func (b BreakResult) Error() string {
	return "break"
}

type ContinueResult struct {
}

func (c ContinueResult) Error() string {
	return "continue"
}

func NewInterpreter(errorReporter ErrorReporter) Interpreter {
	globals := NewEnvironment()
	globals.Define("clock", NewClockCallable())
	return Interpreter{
		errorReporter: errorReporter,
		globals:       &globals,
		environment:   &globals,
		lastValue:     nil,
	}
}

func (inter *Interpreter) Interpret(statements []Stmt) (interface{}, error) {
	for _, stmt := range statements {
		inter.execute(stmt)
	}
	return inter.lastValue, nil
}

func (inter *Interpreter) GetLastValue() (interface{}, error) {
	return inter.lastValue, nil
}

func (inter *Interpreter) execute(stmt Stmt) (interface{}, error) {
	return stmt.accept(inter)
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
	case int32:
		return float64(val), nil
	case int64:
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

func (inter *Interpreter) visitBlockStmt(stmt BlockStmt) (interface{}, error) {
	localEnv := NewEnvironmentWithEnclosing(inter.environment)
	return inter.executeBlock(stmt.Statements, &localEnv)
}

func (inter *Interpreter) visitExpressionStmt(stmt ExpressionStmt) (interface{}, error) {
	var err error = nil
	inter.lastValue, err = inter.evaluate(stmt.Expression)
	return inter.lastValue, err
}

func (inter *Interpreter) visitPrintStmt(stmt PrintStmt) (interface{}, error) {
	value, err := inter.evaluate(stmt.Print)
	fmt.Println(value)
	inter.lastValue = value
	return inter.lastValue, err
}

func (inter *Interpreter) visitVarStmt(stmt VarStmt) (interface{}, error) {
	var value interface{} = nil
	var err error = nil
	if stmt.Initializer != nil {
		value, err = inter.evaluate(stmt.Initializer)
		if err != nil {
			return nil, err
		}
	}
	inter.environment.Define(stmt.Name.Lexeme, value)
	inter.lastValue = value
	return value, nil
}

func (inter *Interpreter) visitIfStmt(stmt IfStmt) (interface{}, error) {
	evalResult, err := inter.evaluate(stmt.Condition)
	if err != nil {
		return nil, err
	}
	conditionVal, err := isTruthy(evalResult)
	if err != nil {
		inter.errorReporter.Push(stmt.Condition.getLine(), INTERPRETER_WHERE, err)
		return nil, err
	}
	if conditionVal {
		return inter.execute(stmt.ThenBranch)
	} else if stmt.ElseBranch != nil {
		return inter.execute(stmt.ElseBranch)
	}
	return nil, nil
}

func (inter *Interpreter) visitWhileStmt(stmt WhileStmt) (interface{}, error) {
	for {
		evalResult, err := inter.evaluate(stmt.Condition)
		if err != nil {
			return nil, err
		}
		keepRunning, err := isTruthy(evalResult)
		if err != nil {
			return nil, err
		}
		if !keepRunning {
			break
		}
		_, result := inter.execute(stmt.Body)
		switch result.(type) {
		case BreakResult:
			keepRunning = false
		case ContinueResult:
			continue
		}
		if !keepRunning {
			break
		}
	}
	return nil, nil
}

func (inter *Interpreter) visitBreakStmt(stmt BreakStmt) (interface{}, error) {
	return nil, BreakResult{}
}

func (inter *Interpreter) visitContinueStmt(stmt ContinueStmt) (interface{}, error) {
	return nil, ContinueResult{}
}

func (inter *Interpreter) visitFunctionStmt(stmt FunctionStmt) (interface{}, error) {
	return nil, nil
}

func (inter *Interpreter) visitLiteralExpr(expr LiteralExpr) (interface{}, error) {
	return expr.Value, nil
}

func (inter *Interpreter) visitGroupingExpr(expr GroupingExpr) (interface{}, error) {
	return inter.evaluate(expr.Expression)
}

func (inter *Interpreter) visitLogicalExpr(expr LogicalExpr) (interface{}, error) {
	left, err := inter.evaluate(expr.Left)
	if err != nil {
		return nil, err
	}
	leftVal, err := isTruthy(left)
	if err != nil {
		inter.errorReporter.Push(expr.Left.getLine(), INTERPRETER_WHERE, err)
		return nil, err
	}
	if expr.Operator.Type == TOKEN_OR {
		if leftVal {
			return left, nil
		}
	} else {
		if !leftVal {
			return leftVal, nil
		}
	}
	return inter.evaluate(expr.Right)
}

func (inter *Interpreter) visitUnaryExpr(expr UnaryExpr) (interface{}, error) {
	right, err := inter.evaluate(expr.Right)
	if err != nil {
		return nil, err
	}

	switch expr.Operator.Type {
	case TOKEN_BANG:
		val, err := isTruthy(right)
		return val, err
	case TOKEN_MINUS:
		val, err := anyToFloat64(right)
		return val, err
	}

	// unreachable
	return nil, nil
}

func (inter *Interpreter) visitBinaryExpr(expr BinaryExpr) (interface{}, error) {
	left, err := inter.evaluate(expr.Left)
	if err != nil {
		return nil, err
	}
	right, err := inter.evaluate(expr.Right)
	if err != nil {
		return nil, err
	}

	switch expr.Operator.Type {
	case TOKEN_GREATER:
		leftVal, rightVal, err := inter.checkNumberOperands(expr, expr.Operator, left, right)
		if err != nil {
			return nil, err
		}
		return leftVal > rightVal, nil
	case TOKEN_GREATER_EQUAL:
		leftVal, rightVal, err := inter.checkNumberOperands(expr, expr.Operator, left, right)
		if err != nil {
			return nil, err
		}
		return leftVal >= rightVal, nil
	case TOKEN_LESS:
		leftVal, rightVal, err := inter.checkNumberOperands(expr, expr.Operator, left, right)
		if err != nil {
			return nil, err
		}
		return leftVal < rightVal, nil
	case TOKEN_LESS_EQUAL:
		leftVal, rightVal, err := inter.checkNumberOperands(expr, expr.Operator, left, right)
		if err != nil {
			return nil, err
		}
		return leftVal <= rightVal, nil
	case TOKEN_MINUS:
		leftVal, rightVal, err := inter.checkNumberOperands(expr, expr.Operator, left, right)
		if err != nil {
			return nil, err
		}
		return leftVal - rightVal, nil
	case TOKEN_SLASH:
		leftVal, rightVal, err := inter.checkNumberOperands(expr, expr.Operator, left, right)
		if err != nil {
			return nil, err
		}
		return leftVal / rightVal, nil
	case TOKEN_STAR:
		leftVal, rightVal, err := inter.checkNumberOperands(expr, expr.Operator, left, right)
		if err != nil {
			return nil, err
		}
		return leftVal * rightVal, nil
	case TOKEN_PLUS:
		if isNumber(left) && isNumber(right) {
			leftVal, rightVal, err := inter.checkNumberOperands(expr, expr.Operator, left, right)
			if err != nil {
				return nil, err
			}
			return leftVal + rightVal, nil
		} else if isString(left) && isString(right) {
			leftVal, rightVal, err := inter.checkStringOperands(expr, expr.Operator, left, right)
			if err != nil {
				return nil, err
			}
			return leftVal + rightVal, nil
		}
		return nil, fmt.Errorf("operands must be two numbers or two strings")
	case TOKEN_BANG_EQUAL:
		return !isEqual(left, right), nil
	case TOKEN_EQUAL_EQUAL:
		return isEqual(left, right), nil
	}

	// unreachable
	return nil, fmt.Errorf("unreachable code")
}

func (inter *Interpreter) visitConditionalExpr(expr ConditionalExpr) (interface{}, error) {
	condition, err := inter.evaluate(expr.Condition)
	if err != nil {
		return nil, err
	}
	if val, _ := isTruthy(condition); val {
		return inter.evaluate(expr.Left)
	} else {
		return inter.evaluate(expr.Right)
	}
}

func (inter *Interpreter) visitVariableExpr(expr VariableExpr) (interface{}, error) {
	value, err := inter.environment.Get(expr.Name.Lexeme)
	if err != nil {
		return nil, err
	}
	return value, nil
}

func (inter *Interpreter) visitAssignExpr(expr AssignExpr) (interface{}, error) {
	value, err := inter.evaluate(expr.Value)
	if err != nil {
		return nil, err
	}
	inter.environment.Assign(expr.Name.Lexeme, value)
	return value, nil
}

func (inter *Interpreter) visitCallExpr(expr CallExpr) (interface{}, error) {
	callee, err := inter.evaluate(expr.Callee)
	if err != nil {
		return nil, err
	}

	var argumentValues []interface{} = []interface{}{}
	for _, argumentExpr := range expr.Arguments {
		argumentValue, err := inter.evaluate(argumentExpr)
		if err != nil {
			return nil, err
		}
		argumentValues = append(argumentValues, argumentValue)
	}

	switch callee := callee.(type) {
	case Callable:
		argumentCount := len(argumentValues)
		if argumentCount != callee.getArity() {
			err = fmt.Errorf("expected %d arguments but got %d", callee.getArity(), argumentCount)
			inter.errorReporter.Push(expr.getLine(), INTERPRETER_WHERE, err)
			return nil, err
		}
		return callee.call(inter, argumentValues)
	default:
		err = fmt.Errorf("can only call function and classes")
		inter.errorReporter.Push(expr.getLine(), INTERPRETER_WHERE, err)
		return nil, err
	}
}

func (inter *Interpreter) evaluate(expr Expr) (interface{}, error) {
	return expr.accept(inter)
}

func (inter *Interpreter) executeBlock(statements []Stmt, localEnv *Environment) (interface{}, error) {
	previousEnv := inter.environment
	inter.environment = localEnv
	for _, stmt := range statements {
		_, result := inter.execute(stmt)
		switch result.(type) {
		case ContinueResult:
			return nil, result
		case BreakResult:
			return nil, result
		}
	}
	inter.environment = previousEnv
	return nil, nil
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
