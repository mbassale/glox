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
		{"basic operators", "1*2+3/4", NewBinaryExpr(
			NewBinaryExpr(
				NewLiteralExpr(NewToken(TOKEN_NUMBER, "1", 1.0, 1)),
				NewToken(TOKEN_STAR, "*", nil, 1),
				NewLiteralExpr(NewToken(TOKEN_NUMBER, "2", 2.0, 1)),
			),
			NewToken(TOKEN_PLUS, "+", nil, 1),
			NewBinaryExpr(
				NewLiteralExpr(NewToken(TOKEN_NUMBER, "3", 3.0, 1)),
				NewToken(TOKEN_SLASH, "/", nil, 1),
				NewLiteralExpr(NewToken(TOKEN_NUMBER, "4", 4.0, 1)),
			),
		)},
		{"ternary expression", "2>=1?\"2\"==2?true:false:false", NewConditionalExpr(
			NewBinaryExpr(
				NewLiteralExpr(NewToken(TOKEN_NUMBER, "2", 2.0, 1)),
				NewToken(TOKEN_GREATER_EQUAL, ">=", nil, 1),
				NewLiteralExpr(NewToken(TOKEN_NUMBER, "1", 1.0, 1)),
			),
			NewConditionalExpr(
				NewBinaryExpr(
					NewLiteralExpr(NewToken(TOKEN_STRING, "\"2\"", "2", 1)),
					NewToken(TOKEN_EQUAL_EQUAL, "==", nil, 1),
					NewLiteralExpr(NewToken(TOKEN_NUMBER, "2", 2.0, 1)),
				),
				NewLiteralExpr(true),
				NewLiteralExpr(false),
			),
			NewLiteralExpr(false),
		)},
	}
	for _, testCase := range testCases {
		errorReporter := NewConsoleErrorReporter()
		scanner := NewScanner(testCase.source, errorReporter)
		tokens := scanner.ScanTokens()
		parser := NewParser(tokens)
		currentExpr, err := parser.Parse()
		assert.Nil(t, err)
		assert.NotNil(t, currentExpr)
		assert.Equal(t, testCase.expectedExpr, currentExpr, testCase.name)
	}

}
