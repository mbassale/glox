package glox

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParserStatements(t *testing.T) {
	testCases := []struct {
		name               string
		source             string
		expectedStatements []Stmt
	}{
		{
			"basic operators",
			"1*2+3/4;",
			[]Stmt{
				NewExpressionStmt(
					NewBinaryExpr(
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
					),
				),
			},
		},
		{
			"ternary expression",
			"2>=1?\"2\"==2?true:false:false;",
			[]Stmt{
				NewExpressionStmt(
					NewConditionalExpr(
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
					),
				),
			},
		},
		{
			"assign expression",
			"var test=1; test=test*2+1;",
			[]Stmt{
				NewVarStmt(
					NewToken(TOKEN_IDENTIFIER, "test", "test", 1),
					NewLiteralExpr(1.0, 1),
				),
				NewExpressionStmt(
					NewAssignExpr(
						NewToken(TOKEN_IDENTIFIER, "test", "test", 1),
						NewBinaryExpr(
							NewBinaryExpr(
								NewVariableExpr(
									NewToken(TOKEN_IDENTIFIER, "test", "test", 1),
								),
								NewToken(TOKEN_STAR, "*", nil, 1),
								NewLiteralExpr(2.0, 1),
							),
							NewToken(TOKEN_PLUS, "+", nil, 1),
							NewLiteralExpr(1.0, 1),
						),
					),
				),
			},
		},
		{
			"If-Then-Else Stmt",
			"if (true) { print 1; } else { print 2; }",
			[]Stmt{
				NewIfStmt(
					NewLiteralExpr(true, 1),
					NewBlockStmt([]Stmt{
						NewPrintStmt(
							NewLiteralExpr(1.0, 1),
						),
					}),
					NewBlockStmt([]Stmt{
						NewPrintStmt(
							NewLiteralExpr(2.0, 1),
						),
					}),
				),
			},
		},
		{
			"If(OrLogicalExpr)-Then-Else",
			"if (3<1 or 1>=3 or 1==1) { print 1; } else { print 2; }",
			[]Stmt{
				NewIfStmt(
					NewLogicalExpr(
						NewLogicalExpr(
							NewBinaryExpr(
								NewLiteralExpr(3.0, 1),
								NewToken(TOKEN_LESS, "<", nil, 1),
								NewLiteralExpr(1.0, 1),
							),
							NewToken(TOKEN_OR, "or", nil, 1),
							NewBinaryExpr(
								NewLiteralExpr(1.0, 1),
								NewToken(TOKEN_GREATER_EQUAL, ">=", nil, 1),
								NewLiteralExpr(3.0, 1),
							),
						),
						NewToken(TOKEN_OR, "or", nil, 1),
						NewBinaryExpr(
							NewLiteralExpr(1.0, 1),
							NewToken(TOKEN_EQUAL_EQUAL, "==", nil, 1),
							NewLiteralExpr(1.0, 1),
						),
					),
					NewBlockStmt([]Stmt{
						NewPrintStmt(
							NewLiteralExpr(1.0, 1),
						),
					}),
					NewBlockStmt([]Stmt{
						NewPrintStmt(
							NewLiteralExpr(2.0, 1),
						),
					}),
				),
			},
		},
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
			assert.Equal(t, testCase.expectedStatements, statements, testCase.name)
		}
	}

}
