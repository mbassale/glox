package glox

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAstPrinterPrint(t *testing.T) {
	expression := NewBinaryExpr(
		NewUnaryExpr(NewToken(TOKEN_MINUS, "-", nil, 1), NewLiteralExpr(123, 1)),
		NewToken(TOKEN_STAR, "*", nil, 1),
		NewGroupingExpr(NewLiteralExpr(45.67, 1)),
	)
	astPrinter := AstPrinter{}
	assert.Equal(t, "(* (- 123) (group 45.67))", astPrinter.Print(expression))
}
