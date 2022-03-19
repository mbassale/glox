package glox

import "fmt"

type Resolver struct {
	inter  *Interpreter
	scopes []map[string]bool
}

func NewResolver(inter *Interpreter) Resolver {
	return Resolver{
		inter:  inter,
		scopes: []map[string]bool{},
	}
}

func (r *Resolver) beginScope() {
	r.scopes = append(r.scopes, map[string]bool{})
}

func (r *Resolver) endScope() {
	r.scopes = r.scopes[:len(r.scopes)-1]
}

func (r *Resolver) ResolveStatements(statements []Stmt) {
	for _, stmt := range statements {
		r.resolveStatement(stmt)
	}
}

func (r *Resolver) resolveStatement(stmt Stmt) {
	stmt.accept(r)
}

func (r *Resolver) resolveExpression(expr Expr) {
	expr.accept(r)
}

func (r *Resolver) resolveLocal(expr Expr, name Token) {
	for i := len(r.scopes) - 1; i >= 0; i-- {
		if _, ok := r.scopes[i][name.Lexeme]; ok {
			r.inter.resolve(expr, len(r.scopes)-1-i)
			return
		}
	}
}

func (r *Resolver) resolveFunction(funcStmt FunctionStmt) {
	r.beginScope()
	for _, param := range funcStmt.Params {
		r.declare(param)
		r.define(param)
	}
	r.ResolveStatements(funcStmt.Body)
	r.endScope()
}

func (r *Resolver) declare(name Token) {
	if len(r.scopes) == 0 {
		return
	}
	r.scopes[len(r.scopes)-1][name.Lexeme] = false
}

func (r *Resolver) define(name Token) {
	if len(r.scopes) == 0 {
		return
	}
	r.scopes[len(r.scopes)-1][name.Lexeme] = true
}

func (r *Resolver) visitBlockStmt(stmt BlockStmt) (interface{}, error) {
	r.beginScope()
	r.ResolveStatements(stmt.Statements)
	r.endScope()
	return nil, nil
}

func (r *Resolver) visitExpressionStmt(stmt ExpressionStmt) (interface{}, error) {
	r.resolveExpression(stmt.Expression)
	return nil, nil
}

func (r *Resolver) visitPrintStmt(stmt PrintStmt) (interface{}, error) {
	r.resolveExpression(stmt.Print)
	return nil, nil
}

func (r *Resolver) visitVarStmt(stmt VarStmt) (interface{}, error) {
	r.declare(stmt.Name)
	if stmt.Initializer != nil {
		r.resolveExpression(stmt.Initializer)
	}
	r.define(stmt.Name)
	return nil, nil
}

func (r *Resolver) visitIfStmt(stmt IfStmt) (interface{}, error) {
	r.resolveExpression(stmt.Condition)
	r.resolveStatement(stmt.ThenBranch)
	if stmt.ElseBranch != nil {
		r.resolveStatement(stmt.ElseBranch)
	}
	return nil, nil
}

func (r *Resolver) visitWhileStmt(stmt WhileStmt) (interface{}, error) {
	r.resolveExpression(stmt.Condition)
	r.resolveStatement(stmt.Body)
	return nil, nil
}

func (r *Resolver) visitBreakStmt(stmt BreakStmt) (interface{}, error) {
	return nil, nil
}

func (r *Resolver) visitContinueStmt(stmt ContinueStmt) (interface{}, error) {
	return nil, nil
}

func (r *Resolver) visitFunctionStmt(stmt FunctionStmt) (interface{}, error) {
	r.declare(stmt.Name)
	r.define(stmt.Name)
	r.resolveFunction(stmt)
	return nil, nil
}

func (r *Resolver) visitReturnStmt(stmt ReturnStmt) (interface{}, error) {
	if stmt.Value != nil {
		r.resolveExpression(stmt.Value)
	}
	return nil, nil
}

func (r *Resolver) visitBinaryExpr(expr BinaryExpr) (interface{}, error) {
	r.resolveExpression(expr.Left)
	r.resolveExpression(expr.Right)
	return nil, nil
}

func (r *Resolver) visitConditionalExpr(expr ConditionalExpr) (interface{}, error) {
	return nil, nil
}

func (r *Resolver) visitGroupingExpr(expr GroupingExpr) (interface{}, error) {
	r.resolveExpression(expr.Expression)
	return nil, nil
}

func (r *Resolver) visitLiteralExpr(expr LiteralExpr) (interface{}, error) {
	return nil, nil
}

func (r *Resolver) visitLogicalExpr(expr LogicalExpr) (interface{}, error) {
	r.resolveExpression(expr.Left)
	r.resolveExpression(expr.Right)
	return nil, nil
}

func (r *Resolver) visitUnaryExpr(expr UnaryExpr) (interface{}, error) {
	r.resolveExpression(expr.Right)
	return nil, nil
}

func (r *Resolver) visitVariableExpr(expr VariableExpr) (interface{}, error) {
	if len(r.scopes) > 0 && !r.scopes[len(r.scopes)-1][expr.Name.Lexeme] {
		return nil, fmt.Errorf("can't read local variable in its own initializer")
	}
	r.resolveLocal(expr, expr.Name)
	return nil, nil
}

func (r *Resolver) visitAssignExpr(expr AssignExpr) (interface{}, error) {
	r.resolveExpression(expr.Value)
	r.resolveLocal(expr, expr.Name)
	return nil, nil
}

func (r *Resolver) visitCallExpr(expr CallExpr) (interface{}, error) {
	r.resolveExpression(expr.Callee)
	for _, argument := range expr.Arguments {
		r.resolveExpression(argument)
	}
	return nil, nil
}
