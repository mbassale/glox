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
}

func (c FunctionCallable) getArity() int {
	return len(c.declaration.Params)
}

func (c FunctionCallable) call(inter *Interpreter, arguments []interface{}) (interface{}, error) {
	env := NewEnvironmentWithEnclosing(inter.globals)
	for i, arg := range c.declaration.Params {
		env.Define(arg.Lexeme, arguments[i])
	}
	return inter.executeBlock(c.declaration.Body, &env)
}

func NewFunctionCallable(declaration FunctionStmt) FunctionCallable {
	return FunctionCallable{
		declaration: declaration,
	}
}
