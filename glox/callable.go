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
