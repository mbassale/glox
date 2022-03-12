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
		{"default variable declaration", "var test;", nil},
		{"simple variable declaration", "var test = 2*3;", 6.0},
		{"if-true-then", "if(2<3){var test=1;}", 1.0},
		{"if-false-then", "if(2>3){var test=1;}", nil},
		{"if-true-then-false-else", "if(2<3){var test=1;}else{var test=2;}", 1.0},
		{"if-false-then-true-else", "if(2>3){var test=1;}else{var test=2;}", 2.0},
		{"if(logicalExpr)-true-then-else", "if(3>=3 and 2==2 and 3<4 and true==true){var test=1;}else{var test=2;}", 1.0},
		{"if(logicalExpr)-true-then-else", "if(3>3 or 2==1 or 3>4 or true==true){var test=1;}else{var test=2;}", 1.0},
		{"if(logicalExpr)-true-then-else", "if(3>3 or 2==1 or 3>4 or (true==false)){var test=1;}else{var test=2;}", 2.0},
		{"while(trueLogicalExpr)-block", "var counter=0;while(counter<5){counter=counter+1;}", 5.0},
		{"while(falseLogicalExpr)-block", "while(false){print 1;}", nil},
		{"ForStmt", "for(var i=0;i<5;i=i+1){i;}", 5.0},
		{"ContinueStmt", "var counter=0;while(counter<5){counter=counter+1;continue;counter=0;}", 5.0},
		{"BreakStmt", "var counter=0;while(counter<5){counter=counter+1;break;counter=0;}", 1.0},
		{"CallExpr", "if(clock()>0){var counter=1;}", 1.0},
		{"FunctionStmt", "fun testFunction(arg){var counter=arg;}testFunction(1.0);", 1.0},
		{"ReturnStmt", "fun testFunction(num) { for (var i = 0; i < num; i=i+1) { if (i >= 2) { return i; } } } testFunction(10.0);", 2.0},
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
