package glox

import "log"

type ErrorReporter interface {
	Error(line int, message string)
	HasError() bool
}

type ConsoleErrorReporter struct {
	hasError bool
}

func (er *ConsoleErrorReporter) HasError() bool {
	return er.hasError
}

func (er *ConsoleErrorReporter) Error(line int, message string) {
	er.report(line, "", message)
}

func (er *ConsoleErrorReporter) report(line int, where string, message string) {
	if len(where) == 0 {
		log.Printf("[line %d] Error: %s", line, message)
	} else {
		log.Printf("[line %d] Error %s: %s", line, where, message)
	}
	er.hasError = true
}

func NewConsoleErrorReporter() *ConsoleErrorReporter {
	return &ConsoleErrorReporter{
		hasError: false,
	}
}
