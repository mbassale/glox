package glox

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParserExpressions(t *testing.T) {
	testCases := []struct {
		name         string
		source       string
		expectedExpr Expr
	}{
		{"basic operators", "1*2+3/4;", NewBinaryExpr(
			NewBinaryExpr(
				NewLiteralExpr(1.0, 1),
				NewToken(TOKEN_STAR, "*", nil, 1),
				NewLiteralExpr(2.0, 1),
			),
			NewToken(TOKEN_PLUS, "+", nil, 1),
			NewBinaryExpr(
				NewLiteralExpr(3.0, 1),
				NewToken(TOKEN_SLASH, "/", nil, 1),
				NewLiteralExpr(4.0, 1),
			),
		)},
		{"ternary expression", "2>=1?\"2\"==2?true:false:false;", NewConditionalExpr(
			NewBinaryExpr(
				NewLiteralExpr(2.0, 1),
				NewToken(TOKEN_GREATER_EQUAL, ">=", nil, 1),
				NewLiteralExpr(1.0, 1),
			),
			NewConditionalExpr(
				NewBinaryExpr(
					NewLiteralExpr("2", 1),
					NewToken(TOKEN_EQUAL_EQUAL, "==", nil, 1),
					NewLiteralExpr(2.0, 1),
				),
				NewLiteralExpr(true, 1),
				NewLiteralExpr(false, 1),
			),
			NewLiteralExpr(false, 1),
		)},
	}
	for _, testCase := range testCases {
		errorReporter := NewConsoleErrorReporter()
		scanner := NewScanner(testCase.source, errorReporter)
		tokens := scanner.ScanTokens()
		assert.False(t, errorReporter.HasError())
		errorReporter.ClearError()
		parser := NewParser(tokens, errorReporter)
		statements := parser.Parse()
		assert.False(t, errorReporter.HasError())
		if assert.NotEmpty(t, statements) {
			if assert.IsType(t, NewExpressionStmt(testCase.expectedExpr), statements[0]) {
				expressionStmt := statements[0].(ExpressionStmt)
				currentExpr := expressionStmt.Expression
				assert.Equal(t, testCase.expectedExpr, currentExpr, testCase.name)
			}
		}
	}

}
