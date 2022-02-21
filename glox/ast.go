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
}

type BinaryExpr struct {
	Left     Expr
	Operator Token
	Right    Expr
}

func (e BinaryExpr) accept(visitor Visitor) interface{} {
	return visitor.visitBinaryExpr(e)
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

func NewGroupingExpr(expression Expr) GroupingExpr {
	return GroupingExpr{
		Expression: expression,
	}
}

type LiteralExpr struct {
	Value interface{}
}

func (e LiteralExpr) accept(visitor Visitor) interface{} {
	return visitor.visitLiteralExpr(e)
}

func NewLiteralExpr(value interface{}) LiteralExpr {
	return LiteralExpr{
		Value: value,
	}
}

type UnaryExpr struct {
	Operator Token
	Right    Expr
}

func (e UnaryExpr) accept(visitor Visitor) interface{} {
	return visitor.visitUnaryExpr(e)
}

func NewUnaryExpr(operator Token, right Expr) UnaryExpr {
	return UnaryExpr{
		Operator: operator,
		Right:    right,
	}
}
