package glox

type Visitor interface {
	visitExpressionStmt(stmt ExpressionStmt) interface{}
	visitPrintStmt(stmt PrintStmt) interface{}
	visitVarStmt(stmt VarStmt) interface{}
	visitBinaryExpr(expr BinaryExpr) interface{}
	visitConditionalExpr(expr ConditionalExpr) interface{}
	visitGroupingExpr(expr GroupingExpr) interface{}
	visitLiteralExpr(expr LiteralExpr) interface{}
	visitUnaryExpr(expr UnaryExpr) interface{}
	visitVariableExpr(expr VariableExpr) interface{}
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
	Value interface{}
	Line  int
}

func (e LiteralExpr) accept(visitor Visitor) interface{} {
	return visitor.visitLiteralExpr(e)
}

func (e LiteralExpr) getLine() int {
	return e.Line
}

func NewLiteralExpr(value interface{}, line int) LiteralExpr {
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

type VariableExpr struct {
	Name Token
}

func (e VariableExpr) accept(visitor Visitor) interface{} {
	return visitor.visitVariableExpr(e)
}

func (e VariableExpr) getLine() int {
	return e.Name.Line
}

func NewVariableExpr(name Token) VariableExpr {
	return VariableExpr{
		Name: name,
	}
}

type Stmt interface {
	accept(visitor Visitor) interface{}
	getLine() int
}

type ExpressionStmt struct {
	Expression Expr
}

func (stmt ExpressionStmt) accept(visitor Visitor) interface{} {
	return visitor.visitExpressionStmt(stmt)
}

func (stmt ExpressionStmt) getLine() int {
	return stmt.Expression.getLine()
}

func NewExpressionStmt(expression Expr) ExpressionStmt {
	return ExpressionStmt{
		Expression: expression,
	}
}

type PrintStmt struct {
	Print Expr
}

func (stmt PrintStmt) accept(visitor Visitor) interface{} {
	return visitor.visitPrintStmt(stmt)
}

func (stmt PrintStmt) getLine() int {
	return stmt.Print.getLine()
}

func NewPrintStmt(print Expr) PrintStmt {
	return PrintStmt{
		Print: print,
	}
}

type VarStmt struct {
	Name        Token
	Initializer Expr
}

func (stmt VarStmt) accept(visitor Visitor) interface{} {
	return visitor.visitVarStmt(stmt)
}

func (stmt VarStmt) getLine() int {
	return stmt.Name.Line
}

func NewVarStmt(name Token, initializer Expr) VarStmt {
	return VarStmt{
		Name:        name,
		Initializer: initializer,
	}
}
