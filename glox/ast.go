package glox

type Visitor interface {
	visitBinaryExpr(expr BinaryExpr) interface{}
	visitConditionalExpr(expr ConditionalExpr) interface{}
	visitGroupingExpr(expr GroupingExpr) interface{}
	visitLiteralExpr(expr LiteralExpr) interface{}
	visitUnaryExpr(expr UnaryExpr) interface{}
}

type Expr interface {
	accept(visitor Visitor) interface{}
	getLine() int
}

type BinaryExpr struct {
	Left     Expr
	Operator Token
	Right    Expr
}

func (e BinaryExpr) accept(visitor Visitor) interface{} {
	return visitor.visitBinaryExpr(e)
}

func (e BinaryExpr) getLine() int {
	return e.Operator.Line
}

func NewBinaryExpr(left Expr, operator Token, right Expr) BinaryExpr {
	return BinaryExpr{
		Left:     left,
		Operator: operator,
		Right:    right,
	}
}

type ConditionalExpr struct {
	Condition Expr
	Left      Expr
	Right     Expr
}

func (e ConditionalExpr) accept(visitor Visitor) interface{} {
	return visitor.visitConditionalExpr(e)
}

func (e ConditionalExpr) getLine() int {
	return e.Condition.getLine()
}

func NewConditionalExpr(condition Expr, left Expr, right Expr) ConditionalExpr {
	return ConditionalExpr{
		Condition: condition,
		Left:      left,
		Right:     right,
	}
}

type GroupingExpr struct {
	Expression Expr
}

func (e GroupingExpr) accept(visitor Visitor) interface{} {
	return visitor.visitGroupingExpr(e)
}

func (e GroupingExpr) getLine() int {
	return e.Expression.getLine()
}

func NewGroupingExpr(expression Expr) GroupingExpr {
	return GroupingExpr{
		Expression: expression,
	}
}

type LiteralExpr struct {
	Line  int
	Value interface{}
}

func (e LiteralExpr) accept(visitor Visitor) interface{} {
	return visitor.visitLiteralExpr(e)
}

func (e LiteralExpr) getLine() int {
	return e.Line
}

func NewLiteralExpr(line int, value interface{}) LiteralExpr {
	return LiteralExpr{
		Value: value,
		Line:  line,
	}
}

type UnaryExpr struct {
	Operator Token
	Right    Expr
}

func (e UnaryExpr) accept(visitor Visitor) interface{} {
	return visitor.visitUnaryExpr(e)
}

func (e UnaryExpr) getLine() int {
	return e.Operator.Line
}

func NewUnaryExpr(operator Token, right Expr) UnaryExpr {
	return UnaryExpr{
		Operator: operator,
		Right:    right,
	}
}
