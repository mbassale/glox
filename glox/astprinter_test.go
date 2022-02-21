package glox

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAstPrinterPrint(t *testing.T) {
	expression := BinaryExpr{
		Left: UnaryExpr{
			Operator: NewToken(TOKEN_MINUS, "-", nil, 1),
			Right: LiteralExpr{
				Value: 123,
			},
		},
		Operator: NewToken(TOKEN_STAR, "*", nil, 1),
		Right: GroupingExpr{
			Expression: LiteralExpr{
				Value: 45.67,
			},
		},
	}
	astPrinter := AstPrinter{}
	assert.Equal(t, "(* (- 123) (group 45.67))", astPrinter.Print(expression))
}
