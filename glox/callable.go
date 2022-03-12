package glox

import "time"

type Callable interface {
	getArity() int
	call(inter *Interpreter, arguments []interface{}) (interface{}, error)
}

type ClockCallable struct {
}

func (c ClockCallable) getArity() int {
	return 0
}

func (c ClockCallable) call(inter *Interpreter, arguments []interface{}) (interface{}, error) {
	return time.Now().Unix(), nil
}

func NewClockCallable() ClockCallable {
	return ClockCallable{}
}

type FunctionCallable struct {
	declaration FunctionStmt
	closure     *Environment
}

func (c FunctionCallable) getArity() int {
	return len(c.declaration.Params)
}

func (c FunctionCallable) call(inter *Interpreter, arguments []interface{}) (interface{}, error) {
	env := NewEnvironmentWithEnclosing(c.closure)
	for i, arg := range c.declaration.Params {
		env.Define(arg.Lexeme, arguments[i])
	}
	value, result := inter.executeBlock(c.declaration.Body, &env)
	switch result := result.(type) {
	case ReturnResult:
		return result.value, nil
	}
	return value, result
}

func NewFunctionCallable(declaration FunctionStmt, closure *Environment) FunctionCallable {
	return FunctionCallable{
		declaration: declaration,
		closure:     closure,
	}
}
