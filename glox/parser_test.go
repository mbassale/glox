package glox

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParserExpressions(t *testing.T) {
	errorReporter := NewConsoleErrorReporter()
	scanner := NewScanner("1*2+3/4", errorReporter)
	tokens := scanner.ScanTokens()
	parser := NewParser(tokens)
	currentExpr, err := parser.Parse()
	assert.Nil(t, err)
	assert.NotNil(t, currentExpr)
	expectedExpr := NewBinaryExpr(
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
	)
	assert.Equal(t, expectedExpr, currentExpr)
}
