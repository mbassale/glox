package glox

import "fmt"

type AstPrinter struct {
}

func (p AstPrinter) Print(expr Expr) string {
	return expr.accept(p).(string)
}

func (p AstPrinter) visitBinaryExpr(expr BinaryExpr) interface{} {
	return p.parenthesize(expr.Operator.Lexeme, expr.Left, expr.Right)
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
