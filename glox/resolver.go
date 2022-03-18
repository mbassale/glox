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
		r.resolve(stmt)
	}
}

func (r *Resolver) resolve(stmt Stmt) {
	stmt.accept(r)
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
	r.resolve(stmt.Expression)
	return nil, nil
}

func (r *Resolver) visitPrintStmt(stmt PrintStmt) (interface{}, error) {
	r.resolve(stmt.Print)
	return nil, nil
}

func (r *Resolver) visitVarStmt(stmt VarStmt) (interface{}, error) {
	r.declare(stmt.Name)
	if stmt.Initializer != nil {
		r.resolve(stmt.Initializer)
	}
	r.define(stmt.Name)
	return nil, nil
}

func (r *Resolver) visitIfStmt(stmt IfStmt) (interface{}, error) {
	r.resolve(stmt.Condition)
	r.resolve(stmt.ThenBranch)
	if stmt.ElseBranch != nil {
		r.resolve(stmt.ElseBranch)
	}
	return nil, nil
}

func (r *Resolver) visitWhileStmt(stmt WhileStmt) (interface{}, error) {
	r.resolve(stmt.Condition)
	r.resolve(stmt.Body)
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
		r.resolve(stmt.Value)
	}
	return nil, nil
}

func (r *Resolver) visitBinaryExpr(expr BinaryExpr) (interface{}, error) {
	r.resolve(expr.Left)
	r.resolve(expr.Right)
	return nil, nil
}

func (r *Resolver) visitConditionalExpr(expr ConditionalExpr) (interface{}, error) {
	return nil, nil
}

func (r *Resolver) visitGroupingExpr(expr GroupingExpr) (interface{}, error) {
	r.resolve(expr.Expression)
	return nil, nil
}

func (r *Resolver) visitLiteralExpr(expr LiteralExpr) (interface{}, error) {
	return nil, nil
}

func (r *Resolver) visitLogicalExpr(expr LogicalExpr) (interface{}, error) {
	r.resolve(expr.Left)
	r.resolve(expr.Right)
	return nil, nil
}

func (r *Resolver) visitUnaryExpr(expr UnaryExpr) (interface{}, error) {
	r.resolve(expr.Right)
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
	r.resolve(expr.Value)
	r.resolveLocal(expr, expr.Name)
	return nil, nil
}

func (r *Resolver) visitCallExpr(expr CallExpr) (interface{}, error) {
	r.resolve(expr.Callee)
	for _, argument := range expr.Arguments {
		r.resolve(argument)
	}
	return nil, nil
}
