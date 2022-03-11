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
		{
			"Break-Continue",
			NewWhileStmt(
				NewLiteralExpr(true, 1),
				NewBlockStmt([]Stmt{
					NewContinueStmt(NewToken(TOKEN_CONTINUE, "continue", nil, 1)),
					NewBreakStmt(NewToken(TOKEN_BREAK, "break", nil, 1)),
				}),
			),
			"(while true)\n{\n  (continue);\n  (break);\n}\n",
		},
		{
			"Function Call",
			NewCallExpr(
				NewVariableExpr(
					NewToken(TOKEN_IDENTIFIER, "testFunction", "testFunction", 1),
				),
				NewToken(TOKEN_RIGHT_PAREN, ")", nil, 1),
				[]Expr{
					NewVariableExpr(
						NewToken(TOKEN_IDENTIFIER, "arg1", "arg1", 1),
					),
					NewVariableExpr(
						NewToken(TOKEN_IDENTIFIER, "arg2", "arg2", 1),
					),
					NewLiteralExpr("string arg3", 1),
					NewLiteralExpr(true, 1),
				},
			),
			"(call testFunction arg1 arg2 string arg3 true)",
		},
		{
			"Function Declaration",
			NewFunctionStmt(
				NewToken(TOKEN_IDENTIFIER, "testFunction", "testFunction", 1),
				[]Token{
					NewToken(TOKEN_IDENTIFIER, "arg1", "arg1", 1),
					NewToken(TOKEN_IDENTIFIER, "arg2", "arg2", 1),
					NewToken(TOKEN_IDENTIFIER, "arg3", "arg3", 1),
				},
				[]Stmt{
					NewPrintStmt(
						NewLiteralExpr(1.0, 1),
					),
				},
			),
			"(func testFunction(arg1, arg2, arg3)) {\n  (print 1);\n}\n",
		},
		{
			"Return Stmt",
			NewReturnStmt(
				NewToken(TOKEN_RETURN, "return", "return", 1),
				NewLiteralExpr(true, 1),
			),
			"(return true);",
		},
	}
	for _, testCase := range statements {
		astPrinter := AstPrinter{}
		stmt := testCase.stmt
		assert.Equal(t, testCase.expected_ast, astPrinter.Print([]Stmt{stmt}), testCase.name)
	}
}
