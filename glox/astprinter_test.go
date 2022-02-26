package glox

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAstPrinterPrint(t *testing.T) {
	statements := []struct {
		name         string
		stmt         Stmt
		expected_ast string
	}{
		{
			"Simple ExpressionStmt",
			NewExpressionStmt(NewBinaryExpr(
				NewUnaryExpr(NewToken(TOKEN_MINUS, "-", nil, 1), NewLiteralExpr(123, 1)),
				NewToken(TOKEN_STAR, "*", nil, 1),
				NewGroupingExpr(NewLiteralExpr(45.67, 1)),
			)),
			"( (* (- 123) (group 45.67)));",
		},
		{
			"Var Declaration",
			NewVarStmt(
				NewToken(TOKEN_IDENTIFIER, "test", "test", 1),
				NewBinaryExpr(
					NewUnaryExpr(NewToken(TOKEN_MINUS, "-", nil, 1), NewLiteralExpr(42.42, 1)),
					NewToken(TOKEN_STAR, "+", nil, 1),
					NewGroupingExpr(NewLiteralExpr(2, 1)),
				),
			),
			"(var test (+ (- 42.42) (group 2)));",
		},
	}
	for _, testCase := range statements {
		astPrinter := AstPrinter{}
		stmt := testCase.stmt
		assert.Equal(t, testCase.expected_ast, astPrinter.Print([]Stmt{stmt}), testCase.name)
	}
}
