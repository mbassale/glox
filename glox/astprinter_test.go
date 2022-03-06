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
		{
			"If-Then-Else",
			NewIfStmt(
				NewLiteralExpr(true, 1),
				NewPrintStmt(NewLiteralExpr("then", 1)),
				NewPrintStmt(NewLiteralExpr("else", 2)),
			),
			"(if true)\n(print then);\nelse\n(print else);",
		},
		{
			"While",
			NewWhileStmt(
				NewLiteralExpr(true, 1),
				NewBlockStmt([]Stmt{
					NewPrintStmt(NewLiteralExpr("block", 1)),
				}),
			),
			"(while true)\n{\n  (print block);\n}\n",
		},
	}
	for _, testCase := range statements {
		astPrinter := AstPrinter{}
		stmt := testCase.stmt
		assert.Equal(t, testCase.expected_ast, astPrinter.Print([]Stmt{stmt}), testCase.name)
	}
}
