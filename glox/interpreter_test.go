package glox

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInterpreterExpressions(t *testing.T) {
	testCases := []struct {
		name          string
		source        string
		expectedValue interface{}
	}{
		{"addition", "2+2", 4.0},
		{"operator precedence", "2+3*4-4/2", 12.0},
		{"conditional expression, true", "2<=3?3-1:false", 2.0},
		{"conditional expression, false", "2==3?3-1:false", false},
	}
	for _, testCase := range testCases {
		errorReporter := NewConsoleErrorReporter()
		scanner := NewScanner(testCase.source, errorReporter)
		tokens := scanner.ScanTokens()
		assert.NotNil(t, tokens)
		parser := NewParser(tokens)
		expr, err := parser.Parse()
		assert.Nil(t, err)
		assert.NotNil(t, expr)
		interpreter := NewInterpreter()
		actualValue := interpreter.Interpret(expr)
		assert.Equal(t, testCase.expectedValue, actualValue, testCase.name)
	}
}
