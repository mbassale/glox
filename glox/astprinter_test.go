package glox

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAstPrinterPrint(t *testing.T) {
	expression := NewBinaryExpr(
		NewUnaryExpr(NewToken(TOKEN_MINUS, "-", nil, 1), NewLiteralExpr(1, 123)),
		NewToken(TOKEN_STAR, "*", nil, 1),
		NewGroupingExpr(NewLiteralExpr(1, 45.67)),
	)
	astPrinter := AstPrinter{}
	assert.Equal(t, "(* (- 123) (group 45.67))", astPrinter.Print(expression))
}
