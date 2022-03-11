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
		astStr, _ := stmt.accept(p)
		astStrs = append(astStrs, astStr.(string))
	}
	return strings.Join(astStrs, "\n")
}

func (p AstPrinter) visitBlockStmt(stmt BlockStmt) (interface{}, error) {
	astStrs := []string{}
	for _, stmt := range stmt.Statements {
		ast, _ := stmt.accept(p)
		astStr := "  " + ast.(string)
		astStrs = append(astStrs, astStr)
	}
	return fmt.Sprintf("{\n%s\n}\n", strings.Join(astStrs, "\n")), nil
}

func (p AstPrinter) visitExpressionStmt(stmt ExpressionStmt) (interface{}, error) {
	return p.parenthesize("", stmt.Expression) + ";", nil
}

func (p AstPrinter) visitPrintStmt(stmt PrintStmt) (interface{}, error) {
	return p.parenthesize("print", stmt.Print) + ";", nil
}

func (p AstPrinter) visitVarStmt(stmt VarStmt) (interface{}, error) {
	return p.parenthesize("var "+stmt.Name.Lexeme, stmt.Initializer) + ";", nil
}

func (p AstPrinter) visitIfStmt(stmt IfStmt) (interface{}, error) {
	astStr := p.parenthesize("if", stmt.Condition) + "\n"
	if stmt.ThenBranch != nil {
		ast, _ := stmt.ThenBranch.accept(p)
		astStr += ast.(string)
	} else {
		astStr += ";\n"
	}
	if stmt.ElseBranch != nil {
		ast, _ := stmt.ElseBranch.accept(p)
		astStr += "\nelse\n"
		astStr += ast.(string)
	}
	return astStr, nil
}

func (p AstPrinter) visitWhileStmt(stmt WhileStmt) (interface{}, error) {
	ast, _ := stmt.Body.accept(p)
	return p.parenthesize("while", stmt.Condition) + "\n" + ast.(string), nil
}

func (p AstPrinter) visitBreakStmt(stmt BreakStmt) (interface{}, error) {
	return p.parenthesize("break") + ";", nil
}

func (p AstPrinter) visitContinueStmt(stmt ContinueStmt) (interface{}, error) {
	return p.parenthesize("continue") + ";", nil
}

func (p AstPrinter) visitFunctionStmt(stmt FunctionStmt) (interface{}, error) {
	args := []string{}
	for _, arg := range stmt.Params {
		args = append(args, arg.Lexeme)
	}
	funcSignature := p.parenthesize("func " + stmt.Name.Lexeme + "(" + strings.Join(args, ", ") + ")")
	astStrs := []string{}
	for _, stmt := range stmt.Body {
		ast, _ := stmt.accept(p)
		astStr := "  " + ast.(string)
		astStrs = append(astStrs, astStr)
	}
	return fmt.Sprintf("%s {\n%s\n}\n", funcSignature, strings.Join(astStrs, "\n")), nil
}

func (p AstPrinter) visitReturnStmt(stmt ReturnStmt) (interface{}, error) {
	return p.parenthesize("return", stmt.Value) + ";", nil
}

func (p AstPrinter) visitBinaryExpr(expr BinaryExpr) (interface{}, error) {
	return p.parenthesize(expr.Operator.Lexeme, expr.Left, expr.Right), nil
}

func (p AstPrinter) visitConditionalExpr(expr ConditionalExpr) (interface{}, error) {
	return p.parenthesize("?", expr.Condition, expr.Left, expr.Right), nil
}

func (p AstPrinter) visitGroupingExpr(expr GroupingExpr) (interface{}, error) {
	return p.parenthesize("group", expr.Expression), nil
}

func (p AstPrinter) visitLiteralExpr(expr LiteralExpr) (interface{}, error) {
	if expr.Value == nil {
		return nil, nil
	}
	return fmt.Sprintf("%v", expr.Value), nil
}

func (p AstPrinter) visitLogicalExpr(expr LogicalExpr) (interface{}, error) {
	return p.parenthesize(expr.Operator.Lexeme, expr.Left, expr.Right), nil
}

func (p AstPrinter) visitUnaryExpr(expr UnaryExpr) (interface{}, error) {
	return p.parenthesize(expr.Operator.Lexeme, expr.Right), nil
}

func (p AstPrinter) visitVariableExpr(expr VariableExpr) (interface{}, error) {
	return expr.Name.Lexeme, nil
}

func (p AstPrinter) visitAssignExpr(expr AssignExpr) (interface{}, error) {
	return p.parenthesize("= "+expr.Name.Lexeme, expr.Value), nil
}

func (p AstPrinter) visitCallExpr(expr CallExpr) (interface{}, error) {
	args := []interface{}{expr.Callee}
	for _, arg := range expr.Arguments {
		args = append(args, arg)
	}
	return p.parenthesize("call", args...), nil
}

func (p AstPrinter) parenthesize(name string, exprs ...interface{}) string {
	var builder string
	builder += "(" + name
	for _, expr := range exprs {
		ast, _ := expr.(Expr).accept(p)
		builder += " "
		builder += fmt.Sprintf("%v", ast.(string))
	}
	builder += ")"
	return builder
}
