package glox

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInterpreterStatements(t *testing.T) {
	testCases := []struct {
		name          string
		source        string
		expectedValue interface{}
	}{
		{"addition", "2+2;", 4.0},
		{"operator precedence", "2+3*4-4/2;", 12.0},
		{"conditional expression, true", "2<=3?3-1:false;", 2.0},
		{"conditional expression, false", "2==3?3-1:false;", false},
		{"default variable definition", "var test;", nil},
		{"simple variable definition", "var test = 2*3;", 6.0},
	}
	for _, testCase := range testCases {
		errorReporter := NewConsoleErrorReporter()
		scanner := NewScanner(testCase.source, errorReporter)
		tokens := scanner.ScanTokens()
		assert.NotNil(t, tokens, testCase.name)
		parser := NewParser(tokens, errorReporter)
		statements := parser.Parse()
		assert.False(t, errorReporter.HasError())
		if assert.NotEmpty(t, statements, testCase.name) {
			interpreter := NewInterpreter(errorReporter)
			lastValue, err := interpreter.Interpret(statements)
			if assert.Nil(t, err, testCase.name) {
				assert.Equal(t, testCase.expectedValue, lastValue, testCase.name)
			}
		}
	}
}
