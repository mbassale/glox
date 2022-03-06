package glox

type Visitor interface {
	visitBlockStmt(stmt BlockStmt) interface{}
	visitExpressionStmt(stmt ExpressionStmt) interface{}
	visitPrintStmt(stmt PrintStmt) interface{}
	visitVarStmt(stmt VarStmt) interface{}
	visitIfStmt(stmt IfStmt) interface{}
	visitWhileStmt(stmt WhileStmt) interface{}
	visitBreakStmt(stmt BreakStmt) interface{}
	visitContinueStmt(stmt ContinueStmt) interface{}
	visitBinaryExpr(expr BinaryExpr) interface{}
	visitConditionalExpr(expr ConditionalExpr) interface{}
	visitGroupingExpr(expr GroupingExpr) interface{}
	visitLiteralExpr(expr LiteralExpr) interface{}
	visitLogicalExpr(expr LogicalExpr) interface{}
	visitUnaryExpr(expr UnaryExpr) interface{}
	visitVariableExpr(expr VariableExpr) interface{}
	visitAssignExpr(expr AssignExpr) interface{}
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

type LogicalExpr struct {
	Left     Expr
	Operator Token
	Right    Expr
}

func (e LogicalExpr) accept(visitor Visitor) interface{} {
	return visitor.visitLogicalExpr(e)
}

func (e LogicalExpr) getLine() int {
	return e.Operator.Line
}

func NewLogicalExpr(left Expr, operator Token, right Expr) LogicalExpr {
	return LogicalExpr{
		Left:     left,
		Operator: operator,
		Right:    right,
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

type AssignExpr struct {
	Name  Token
	Value Expr
}

func (e AssignExpr) accept(visitor Visitor) interface{} {
	return visitor.visitAssignExpr(e)
}

func (e AssignExpr) getLine() int {
	return e.Name.Line
}

func NewAssignExpr(name Token, value Expr) AssignExpr {
	return AssignExpr{
		Name:  name,
		Value: value,
	}
}

type Stmt interface {
	accept(visitor Visitor) interface{}
	getLine() int
}

type BlockStmt struct {
	Statements []Stmt
}

func (stmt BlockStmt) accept(visitor Visitor) interface{} {
	return visitor.visitBlockStmt(stmt)
}

func (stmt BlockStmt) getLine() int {
	if len(stmt.Statements) > 0 {
		return stmt.Statements[0].getLine()
	}
	return 0
}

func NewBlockStmt(statements []Stmt) BlockStmt {
	return BlockStmt{
		Statements: statements,
	}
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

type IfStmt struct {
	Condition  Expr
	ThenBranch Stmt
	ElseBranch Stmt
}

func (stmt IfStmt) accept(visitor Visitor) interface{} {
	return visitor.visitIfStmt(stmt)
}

func (stmt IfStmt) getLine() int {
	return stmt.Condition.getLine()
}

func NewIfStmt(condition Expr, thenBranch Stmt, elseBranch Stmt) IfStmt {
	return IfStmt{
		Condition:  condition,
		ThenBranch: thenBranch,
		ElseBranch: elseBranch,
	}
}

type WhileStmt struct {
	Condition Expr
	Body      Stmt
}

func (stmt WhileStmt) accept(visitor Visitor) interface{} {
	return visitor.visitWhileStmt(stmt)
}

func (stmt WhileStmt) getLine() int {
	return stmt.Condition.getLine()
}

func NewWhileStmt(condition Expr, body Stmt) WhileStmt {
	return WhileStmt{
		Condition: condition,
		Body:      body,
	}
}

type BreakStmt struct {
	Token Token
}

func (stmt BreakStmt) accept(visitor Visitor) interface{} {
	return visitor.visitBreakStmt(stmt)
}

func (stmt BreakStmt) getLine() int {
	return stmt.Token.Line
}

func NewBreakStmt(token Token) BreakStmt {
	return BreakStmt{
		Token: token,
	}
}

type ContinueStmt struct {
	Token Token
}

func (stmt ContinueStmt) accept(visitor Visitor) interface{} {
	return visitor.visitContinueStmt(stmt)
}

func (stmt ContinueStmt) getLine() int {
	return stmt.Token.Line
}

func NewContinueStmt(token Token) ContinueStmt {
	return ContinueStmt{
		Token: token,
	}
}
