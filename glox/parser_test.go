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
		{
			"While-Stmt",
			"while(true==1){print 1;break;continue;}",
			[]Stmt{
				NewWhileStmt(
					NewBinaryExpr(
						NewLiteralExpr(true, 1),
						NewToken(TOKEN_EQUAL_EQUAL, "==", nil, 1),
						NewLiteralExpr(1.0, 1),
					),
					NewBlockStmt([]Stmt{
						NewPrintStmt(
							NewLiteralExpr(1.0, 1),
						),
						NewBreakStmt(
							NewToken(TOKEN_BREAK, "break", nil, 1),
						),
						NewContinueStmt(
							NewToken(TOKEN_CONTINUE, "continue", nil, 1),
						),
					}),
				),
			},
		},
		{
			"For-Stmt",
			"for(var i=0;i<10;i=i+1){print 1;}",
			[]Stmt{
				NewBlockStmt([]Stmt{
					NewVarStmt(
						NewToken(TOKEN_IDENTIFIER, "i", "i", 1),
						NewLiteralExpr(0.0, 1),
					),
					NewWhileStmt(
						NewBinaryExpr(
							NewVariableExpr(
								NewToken(TOKEN_IDENTIFIER, "i", "i", 1),
							),
							NewToken(TOKEN_LESS, "<", nil, 1),
							NewLiteralExpr(10.0, 1),
						),
						NewBlockStmt([]Stmt{
							NewBlockStmt([]Stmt{
								NewPrintStmt(
									NewLiteralExpr(1.0, 1),
								),
							}),
							NewExpressionStmt(
								NewAssignExpr(
									NewToken(TOKEN_IDENTIFIER, "i", "i", 1),
									NewBinaryExpr(
										NewVariableExpr(
											NewToken(TOKEN_IDENTIFIER, "i", "i", 1),
										),
										NewToken(TOKEN_PLUS, "+", nil, 1),
										NewLiteralExpr(1.0, 1),
									),
								),
							),
						}),
					),
				}),
			},
		},
		{
			"Function Call Expression",
			"testFunction(arg1, arg2, \"string arg3\", true);",
			[]Stmt{
				NewExpressionStmt(
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
				),
			},
		},
		{
			"Function Declaration Statement",
			"fun testFunction(arg1,arg2,arg3){print 1;}",
			[]Stmt{
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
