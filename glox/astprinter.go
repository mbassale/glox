package glox

import (
	"fmt"
	"strings"
)

type AstPrinter struct {
}

func (p AstPrinter) Print(statements []Stmt) string {
	astStrs := []string{}
	for _, stmt := range statements {
		astStr := stmt.accept(p).(string)
		astStrs = append(astStrs, astStr)
	}
	return strings.Join(astStrs, "\n")
}

func (p AstPrinter) visitBlockStmt(stmt BlockStmt) interface{} {
	astStrs := []string{}
	for _, stmt := range stmt.Statements {
		astStr := "  " + stmt.accept(p).(string)
		astStrs = append(astStrs, astStr)
	}
	return fmt.Sprintf("{\n%s\n}\n", strings.Join(astStrs, "\n"))
}

func (p AstPrinter) visitExpressionStmt(stmt ExpressionStmt) interface{} {
	return p.parenthesize("", stmt.Expression) + ";"
}

func (p AstPrinter) visitPrintStmt(stmt PrintStmt) interface{} {
	return p.parenthesize("print", stmt.Print) + ";"
}

func (p AstPrinter) visitVarStmt(stmt VarStmt) interface{} {
	return p.parenthesize("var "+stmt.Name.Lexeme, stmt.Initializer) + ";"
}

func (p AstPrinter) visitBinaryExpr(expr BinaryExpr) interface{} {
	return p.parenthesize(expr.Operator.Lexeme, expr.Left, expr.Right)
}

func (p AstPrinter) visitConditionalExpr(expr ConditionalExpr) interface{} {
	return p.parenthesize("?", expr.Condition, expr.Left, expr.Right)
}

func (p AstPrinter) visitGroupingExpr(expr GroupingExpr) interface{} {
	return p.parenthesize("group", expr.Expression)
}

func (p AstPrinter) visitLiteralExpr(expr LiteralExpr) interface{} {
	if expr.Value == nil {
		return nil
	}
	return fmt.Sprintf("%v", expr.Value)
}

func (p AstPrinter) visitUnaryExpr(expr UnaryExpr) interface{} {
	return p.parenthesize(expr.Operator.Lexeme, expr.Right)
}

func (p AstPrinter) visitVariableExpr(expr VariableExpr) interface{} {
	return expr.Name.Lexeme
}

func (p AstPrinter) visitAssignExpr(expr AssignExpr) interface{} {
	return p.parenthesize("= "+expr.Name.Lexeme, expr.Value)
}

func (p AstPrinter) parenthesize(name string, exprs ...interface{}) string {
	var builder string
	builder += "(" + name
	for _, expr := range exprs {
		builder += " "
		builder += fmt.Sprintf("%v", expr.(Expr).accept(p))
	}
	builder += ")"
	return builder
}
